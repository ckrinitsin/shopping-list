package shopping_list

import (
	"net/http"

	"github.com/ckrinitsin/go-backend/models"
	"github.com/gin-gonic/gin"
)

func LoadElements(c *gin.Context) {

	title := "Shopping List"
	var entries []models.Entry

	err := models.DB.
		Order("checked asc").
		Find(&entries).
		Error

	if err != nil {
		c.String(http.StatusInternalServerError, "Internal Server Error")
		c.Error(err)
		return
	}

	c.HTML(http.StatusOK, "template.html", gin.H{
		"name":    title,
		"entries": entries,
	})
}

func CreateEntry(c *gin.Context) {
	value := c.PostForm("newItem")

	entry := models.Entry{
		Text:    value,
		Checked: false,
	}

	err := models.DB.
		Create(&entry).
		Error

	if err != nil {
		c.String(http.StatusInternalServerError, "Internal Server Error")
		return
	}

	c.Redirect(http.StatusFound, "/")
}

func DeleteEntries(c *gin.Context) {
	models.DB.Delete(&models.Entry{}, "checked = 1")

	c.Redirect(http.StatusFound, "/")
}

func ToggleEntry(c *gin.Context) {
	id := c.PostForm("id")

	var entry models.Entry

	err := models.DB.
		First(&entry, id).
		Error

	if err != nil {
		c.String(http.StatusInternalServerError, "Internal Server Error")
		return
	}

	entry.Checked = !entry.Checked
	err = models.DB.
		Save(&entry).
		Error

	if err != nil {
		c.String(http.StatusInternalServerError, "Internal Server Error")
		return
	}

	c.Redirect(http.StatusFound, "/")
}
