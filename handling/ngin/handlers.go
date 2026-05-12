package ngin

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"

	"machin8/core/conf"
	"machin8/handling/auth"
	"machin8/processing/xai"
	"machin8/processing/xai/tooling"
)

type notebookState struct {
	mu      sync.Mutex
	history []xai.Message
	pending []xai.Message
	notify  chan struct{}
}

var (
	registryMu sync.RWMutex
	registry   = map[int64]*notebookState{}
)

func registerNotebook(id int64) {
	registryMu.Lock()
	defer registryMu.Unlock()
	registry[id] = &notebookState{
		history: make([]xai.Message, 0, 50),
		pending: make([]xai.Message, 0, 10),
		notify:  make(chan struct{}, 1),
	}
}

func getNotebook(id int64) (*notebookState, bool) {
	registryMu.RLock()
	defer registryMu.RUnlock()
	nb, ok := registry[id]
	return nb, ok
}

func Routing(router *gin.RouterGroup) {
	api := router.Group("/ngin/api/v1", auth.TokenRequired())
	api.POST("/notebook", APIv1NotebookCreate)
	api.GET("/notebook/:id/display", OwnerRequired(), APIv1NotebookDisplay)
	api.POST("/notebook/:id/message", OwnerRequired(), APIv1NotebookMessage)
}

func APIv1NotebookCreate(c *gin.Context) {
	userID := c.MustGet("user_id").(int64)

	db, err := sql.Open("postgres", conf.Settings.DBUrl)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database unavailable"})
		return
	}
	defer db.Close()

	var analistID int64
	err = db.QueryRow(`SELECT id FROM org.analist WHERE user_id = $1`, userID).Scan(&analistID)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusForbidden, gin.H{"error": "analist not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	var notebookID int64
	err = db.QueryRow(
		`INSERT INTO ngin.notebook (analist_id) VALUES ($1) RETURNING id`,
		analistID,
	).Scan(&notebookID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create notebook"})
		return
	}

	registerNotebook(notebookID)
	c.JSON(http.StatusCreated, gin.H{"id": notebookID})
}

func APIv1NotebookMessage(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	nb, ok := getNotebook(id)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "notebook not found"})
		return
	}

	var body struct {
		Message string `json:"message" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid message"})
		return
	}
	msg := xai.Message{Role: "user", Content: body.Message}

	db, err := sql.Open("postgres", conf.Settings.DBUrl)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database unavailable"})
		return
	}
	defer db.Close()

	var modelName, apiKey string
	err = db.QueryRow(
		`SELECT name, api_key
		FROM ngin.model
		WHERE provider = 'xai'
		ORDER BY created DESC
		LIMIT 1`,
	).Scan(&modelName, &apiKey)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusFailedDependency, gin.H{"error": "model not configured"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	nb.mu.Lock()
	nb.pending = append(nb.pending, msg)
	history := append([]xai.Message(nil), nb.history...)
	history = append(history, msg)
	nb.mu.Unlock()

	select {
	case nb.notify <- struct{}{}:
	default:
	}

	go func(notebook *notebookState, messages []xai.Message, model, key string) {
		parkingTools, err := tooling.ParkingTools()
		if err != nil {
			// TODO: surface error to notebook
			return
		}
		if err := xai.Stream(context.Background(), model, key, messages, parkingTools, func(chunk string) error {
			notebook.mu.Lock()
			notebook.pending = append(notebook.pending, xai.Message{Role: "assistant", Content: chunk})
			notebook.mu.Unlock()

			select {
			case notebook.notify <- struct{}{}:
			default:
			}

			return nil
		}); err != nil {
			notebook.mu.Lock()
			notebook.pending = append(notebook.pending, xai.Message{
				Role:    "assistant",
				Content: "Error al generar la respuesta: " + err.Error(),
			})
			notebook.mu.Unlock()

			select {
			case notebook.notify <- struct{}{}:
			default:
			}
		}
	}(nb, history, modelName, apiKey)

	c.JSON(http.StatusAccepted, gin.H{"queued": true})
}

func APIv1NotebookDisplay(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	nb, ok := getNotebook(id)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "notebook not found"})
		return
	}

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	ctx := c.Request.Context()
	for {
		select {
		case <-ctx.Done():
			return
		case <-nb.notify:
			nb.mu.Lock()
			flushed := nb.pending
			nb.pending = []xai.Message{}
			nb.history = append(nb.history, flushed...)
			nb.mu.Unlock()

			for _, msg := range flushed {
				data, _ := json.Marshal(msg)
				c.SSEvent("message", string(data))
				c.Writer.Flush()
			}
		}
	}
}
