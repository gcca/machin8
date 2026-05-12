package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"machin8/core/conf"
	"machin8/handling/auth"
	"machin8/handling/home"
	"machin8/handling/ngin"
)

func main() {
	if err := conf.InitSettings(); err != nil {
		log.Fatalf("config: %v", err)
	}

	r := gin.Default()
	r.LoadHTMLFiles(
		"templates/auth/signin.tmpl",
		"templates/home/error.tmpl",
		"templates/home/analist/index.tmpl",
		"templates/home/analist/home.tmpl",
		"templates/home/analist/notebook/index.tmpl",
		"templates/home/analist/notebook/create.tmpl",
		"templates/home/analist/notebook/list.tmpl",
		"templates/home/analist/notebook/list/table.tmpl",
		"templates/home/analist/notebook/list/table-error.tmpl",
		"templates/home/manager/index.tmpl",
		"templates/home/staff/index.tmpl",
		"templates/home/superuser/index.tmpl",
	)
	machin8 := r.Group("/machin8")
	_ = r.SetTrustedProxies(nil)

	machin8.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/machin8/auth/signin")
	})

	auth.Routing(machin8)
	home.Routing(machin8)
	ngin.Routing(machin8)

	r.Run(":" + conf.Settings.Port)
}
