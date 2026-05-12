package auth

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"

	"machin8/core/conf"
)

func Routing(router *gin.RouterGroup) {
	ui := router.Group("/auth")
	ui.GET("/signin", SignInGet)
	ui.POST("/signin", SignInPost)

	api := router.Group("/auth/api/v1")
	api.POST("/signin", APIv1SignIn)
}

type SignInRequest struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

type APIv1SignInPostRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func APIv1SignIn(c *gin.Context) {
	var req APIv1SignInPostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	db, err := sql.Open("postgres", conf.Settings.DBUrl)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database unavailable"})
		return
	}
	defer db.Close()

	userID, err := Authenticate(db, req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "authentication error"})
		return
	}
	if userID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	expiry := time.Now().Add(365 * 24 * time.Hour)
	rawKey, err := CreateAPIKey(db, *userID, req.Username+" — session key", conf.Settings.Secret, &expiry)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create api key"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"api_key": rawKey})
}

func SignInPost(c *gin.Context) {
	var req SignInRequest
	if err := c.ShouldBind(&req); err != nil {
		c.HTML(http.StatusBadRequest, "signin.tmpl", gin.H{
			"Error": "invalid request body",
		})
		return
	}

	db, err := sql.Open("postgres", conf.Settings.DBUrl)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "signin.tmpl", gin.H{
			"Error": "database unavailable",
		})
		return
	}
	defer db.Close()

	userID, err := Authenticate(db, req.Username, req.Password)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "signin.tmpl", gin.H{
			"Error": "authentication error",
		})
		return
	}
	if userID == nil {
		c.HTML(http.StatusUnauthorized, "signin.tmpl", gin.H{
			"Error": "invalid credentials",
		})
		return
	}

	if err := LogIn(c, db, *userID); err != nil {
		c.HTML(http.StatusInternalServerError, "signin.tmpl", gin.H{
			"Error": "could not create session",
		})
		return
	}

	to := c.Query("to")
	if to == "" {
		c.HTML(http.StatusBadRequest, "signin.tmpl", gin.H{
			"Error": "missing redirect target",
		})
		return
	}

	c.Redirect(http.StatusFound, to)
}

func SignInGet(c *gin.Context) {
	to := c.Query("to")
	if to == "" {
		to = "/machin8/home"
	}

	c.HTML(http.StatusOK, "signin.tmpl", gin.H{
		"To": to,
	})
}
