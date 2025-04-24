package main

import (
	"embed"
	"html/template"
	"net/http"
	"os"

	"github.com/ckrinitsin/shopping-list/authenticate"
	"github.com/ckrinitsin/shopping-list/handlers"
	"github.com/ckrinitsin/shopping-list/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

//go:embed templates/*
var templatesFS embed.FS

func main() {
	r := gin.Default()

	models.ConnectDatabase()

	tmpl := template.Must(template.ParseFS(templatesFS, "templates/*"))
	r.SetHTMLTemplate(tmpl)

	store := cookie.NewStore([]byte(os.Getenv("SECRET")))
	r.Use(sessions.Sessions("session", store))

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"health-check": "passed",
		})
	})

	r.POST("/login", authenticate.LoginPOST)
	r.GET("/login", authenticate.LoginGET)
	r.POST("/register", authenticate.RegisterPOST)
	r.GET("/register", authenticate.RegisterGET)
	r.POST("/logout", authenticate.Logout)

	r.GET("/", authenticate.CheckAuth, shopping_list.LoadElements)
	r.POST("/create", authenticate.CheckAuth, shopping_list.CreateEntry)
	r.POST("/delete", authenticate.CheckAuth, shopping_list.DeleteEntries)
	r.POST("/toggle", authenticate.CheckAuth, shopping_list.ToggleEntry)

	r.Run()
}
