package main

import (
	"embed"
	"html/template"
	"net/http"

	"github.com/ckrinitsin/go-backend/handlers"
	"github.com/ckrinitsin/go-backend/models"
	"github.com/gin-gonic/gin"
)

//go:embed templates/*
var templatesFS embed.FS

func main() {
	r := gin.Default()

	models.ConnectDatabase()

	tmpl := template.Must(template.ParseFS(templatesFS, "templates/*"))
	r.SetHTMLTemplate(tmpl)

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"health-check": "passed",
		})
	})

	r.GET("/shopping", shopping_list.LoadElements)

	r.Run()
}
