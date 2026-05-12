package home

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"

	"machin8/core/conf"
	"machin8/handling/auth"
	"machin8/handling/home/analist"
)

func Routing(router *gin.RouterGroup) {
	ui := router.Group("/home")
	ui.Use(auth.LogInRequired())
	ui.GET("", HomeGet)
	ui.GET("/:kind", HomeKindGet)
	analist.Routing(ui.Group("/analist"))
}

func HomeGet(c *gin.Context) {
	userID := c.MustGet("user_id").(int64)

	db, err := sql.Open("postgres", conf.Settings.DBUrl)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "home/error.tmpl", gin.H{
			"Error": "database unavailable",
		})
		return
	}
	defer db.Close()

	var kind string
	err = db.QueryRow(`SELECT role::text FROM auth."user" WHERE id = $1`, userID).Scan(&kind)
	if err == sql.ErrNoRows {
		c.HTML(http.StatusNotFound, "home/error.tmpl", gin.H{
			"Error": "user not found",
		})
		return
	}
	if err != nil {
		c.HTML(http.StatusInternalServerError, "home/error.tmpl", gin.H{
			"Error": "database error",
		})
		return
	}

	c.Redirect(http.StatusFound, "/machin8/home/"+kind)
}

func HomeKindGet(c *gin.Context) {
	kind := c.Param("kind")
	c.HTML(http.StatusOK, "home/"+kind+"/index.tmpl", gin.H{
		"Kind": kind,
	})
}
