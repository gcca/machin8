package analist

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"

	"machin8/core/conf"
)

func Routing(router gin.IRouter) {
	ui := router.Group("/")
	ui.GET("", IndexGet)
	ui.GET("/home", HomeGet)
	ui.GET("/notebook", NotebookGet)
	ui.GET("/notebook/list", NotebookListGet)
	ui.GET("/notebook/list/table", NotebookListTableGet)
	ui.GET("/notebook/create", NotebookCreateGet)
	ui.POST("/notebook/create", NotebookCreatePost)
}

func IndexGet(c *gin.Context) {
	tab := c.Query("t")
	if tab != "home" && tab != "notebook" {
		if tab != "" {
			tab = "blank"
		} else {
			tab = "home"
		}
	}
	if tab == "" {
		tab = "home"
	}

	c.HTML(http.StatusOK, "home/analist/index.tmpl", gin.H{
		"Tab": tab,
	})
}

func HomeGet(c *gin.Context) {
	c.HTML(http.StatusOK, "home/analist/home.tmpl", gin.H{})
}

func NotebookGet(c *gin.Context) {
	c.HTML(http.StatusOK, "home/analist/notebook/index.tmpl", gin.H{})
}

func NotebookListGet(c *gin.Context) {
	c.HTML(http.StatusOK, "home/analist/notebook/list.tmpl", gin.H{})
}

type notebookRow struct {
	ID      int64
	Title   sql.NullString
	Created string
	Updated string
}

type notebookListPage struct {
	Notebooks []notebookRow
	Page      int
	HasPrev   bool
	HasNext   bool
	PrevPage  int
	NextPage  int
}

func NotebookListTableGet(c *gin.Context) {
	userID := c.MustGet("user_id").(int64)
	page := 1
	if raw := c.Query("page"); raw != "" {
		if parsed, err := strconv.Atoi(raw); err == nil && parsed > 0 {
			page = parsed
		}
	}

	pageData, err := queryNotebookPage(userID, page)
	if err != nil {
		status := http.StatusInternalServerError
		message := "database error"
		if err == sql.ErrNoRows {
			status = http.StatusForbidden
			message = "analist not found"
		}
		if err == sql.ErrConnDone {
			message = "database unavailable"
		}
		c.HTML(status, "home/analist/notebook/list/table-error.tmpl", gin.H{
			"Error": message,
		})
		return
	}

	c.HTML(http.StatusOK, "home/analist/notebook/list/table.tmpl", gin.H{
		"Notebooks": pageData.Notebooks,
		"Page":      pageData.Page,
		"HasPrev":   pageData.HasPrev,
		"HasNext":   pageData.HasNext,
		"PrevPage":  pageData.PrevPage,
		"NextPage":  pageData.NextPage,
		"CloseModal": true,
		"SwapOOB":    true,
	})
}

func queryNotebookPage(userID int64, page int) (*notebookListPage, error) {
	db, err := sql.Open("postgres", conf.Settings.DBUrl)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	const pageSize = 11
	offset := (page - 1) * pageSize

	rows, err := db.Query(
		`SELECT n.id, n.title, n.created, n.updated
		 FROM ngin.notebook n
		 JOIN org.analist a ON a.id = n.analist_id
		 WHERE a.user_id = $1
		 ORDER BY n.created DESC
		 LIMIT $2 OFFSET $3`,
		userID,
		pageSize+1,
		offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notebooks []notebookRow
	for rows.Next() {
		var n notebookRow
		var created, updated sql.NullTime
		if err := rows.Scan(&n.ID, &n.Title, &created, &updated); err != nil {
			return nil, err
		}
		if created.Valid {
			n.Created = created.Time.Format(timeLayout)
		}
		if updated.Valid {
			n.Updated = updated.Time.Format(timeLayout)
		}
		notebooks = append(notebooks, n)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	hasNext := len(notebooks) > pageSize
	if hasNext {
		notebooks = notebooks[:pageSize]
	}

	prevPage := page - 1
	if prevPage < 1 {
		prevPage = 1
	}

	return &notebookListPage{
		Notebooks: notebooks,
		Page:      page,
		HasPrev:   page > 1,
		HasNext:   hasNext,
		PrevPage:  prevPage,
		NextPage:  page + 1,
	}, nil
}

func NotebookCreateGet(c *gin.Context) {
	c.HTML(http.StatusOK, "home/analist/notebook/create.tmpl", gin.H{})
}

func NotebookCreatePost(c *gin.Context) {
	userID := c.MustGet("user_id").(int64)
	title := c.PostForm("title")

	db, err := sql.Open("postgres", conf.Settings.DBUrl)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "home/analist/notebook/create.tmpl", gin.H{
			"Error": "database unavailable",
		})
		return
	}
	defer db.Close()

	var analistID int64
	err = db.QueryRow(`SELECT id FROM org.analist WHERE user_id = $1`, userID).Scan(&analistID)
	if err == sql.ErrNoRows {
		c.HTML(http.StatusForbidden, "home/analist/notebook/create.tmpl", gin.H{
			"Error": "analist not found",
		})
		return
	}
	if err != nil {
		c.HTML(http.StatusInternalServerError, "home/analist/notebook/create.tmpl", gin.H{
			"Error": "database error",
		})
		return
	}

	var notebookID int64
	err = db.QueryRow(
		`INSERT INTO ngin.notebook (analist_id, title) VALUES ($1, NULLIF($2, '')) RETURNING id`,
		analistID,
		title,
	).Scan(&notebookID)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "home/analist/notebook/create.tmpl", gin.H{
			"Error": "failed to create notebook",
		})
		return
	}

	pageData, err := queryNotebookPage(userID, 1)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "home/analist/notebook/create.tmpl", gin.H{
			"Error": "could not refresh notebook list",
		})
		return
	}

	c.HTML(http.StatusOK, "home/analist/notebook/list/table.tmpl", gin.H{
		"Notebooks":  pageData.Notebooks,
		"CloseModal": true,
		"SwapOOB":    true,
		"Page":       pageData.Page,
		"HasPrev":    pageData.HasPrev,
		"HasNext":    pageData.HasNext,
		"PrevPage":   pageData.PrevPage,
		"NextPage":   pageData.NextPage,
	})
}

const timeLayout = "2006-01-02 15:04"
