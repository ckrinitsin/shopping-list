package main

import (
	"embed"
	"html/template"
	"net/http"

	"github.com/ckrinitsin/shopping-list/handlers"
	"github.com/ckrinitsin/shopping-list/models"
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

	r.GET("/", shopping_list.LoadElements)
	r.POST("/create", shopping_list.CreateEntry)
	r.POST("/delete", shopping_list.DeleteEntries)
	r.POST("/toggle", shopping_list.ToggleEntry)

	r.Run()
}
