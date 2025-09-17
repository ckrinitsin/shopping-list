package shopping_list

import (
	"net/http"

	"github.com/ckrinitsin/shopping-list/models"
	"github.com/gin-gonic/gin"
)


func LoadElements(c *gin.Context) {
	var entries []models.Entry

	any_list, ok := c.Get("current_list")
	if !ok {
		c.String(http.StatusInternalServerError, "Internal Server Error")
		return
	}

	list, ok := any_list.(models.List)
	if !ok {
		c.String(http.StatusInternalServerError, "Internal Server Error")
		return
	}

	err := models.DB.
		Where("list_name = ?", list.Name).
		Order("checked asc").
		Find(&entries).
		Error

	if err != nil {
		c.String(http.StatusInternalServerError, "Internal Server Error")
		c.Error(err)
		return
	}

	c.HTML(http.StatusOK, "template.html", gin.H{
		"name":      list.Name,
		"entries":   entries,
		"base_path": models.BasePath(),
	})
}

func CreateEntry(c *gin.Context) {
	value := c.PostForm("newItem")

	any_list, ok := c.Get("current_list")
	if !ok {
		c.String(http.StatusInternalServerError, "Internal Server Error")
		return
	}

	list, ok := any_list.(models.List)
	if !ok {
		c.String(http.StatusInternalServerError, "Internal Server Error")
		return
	}

	entry := models.Entry{
		Text:    value,
		Checked: false,
		ListName: list.Name,
	}

	err := models.DB.
		Create(&entry).
		Error

	if err != nil {
		c.String(http.StatusInternalServerError, "Internal Server Error")
		return
	}

	c.Redirect(http.StatusFound, models.BasePath() + "/")
}

func DeleteEntries(c *gin.Context) {
	any_list, ok := c.Get("current_list")
	if !ok {
		c.String(http.StatusInternalServerError, "Internal Server Error")
		return
	}

	list, ok := any_list.(models.List)
	if !ok {
		c.String(http.StatusInternalServerError, "Internal Server Error")
		return
	}

	err := models.DB.
		Where("list_name = ?", list.Name).
		Delete(&models.Entry{}, "checked = 1").
		Error

	if err != nil {
		c.String(http.StatusInternalServerError, "Internal Server Error")
		return
	}

	c.Redirect(http.StatusFound, models.BasePath() + "/")
}

func ToggleEntry(c *gin.Context) {
	id := c.PostForm("id")
	checked := c.PostForm("checked")

	any_list, ok := c.Get("current_list")
	if !ok {
		c.String(http.StatusInternalServerError, "Internal Server Error")
		return
	}

	list, ok := any_list.(models.List)
	if !ok {
		c.String(http.StatusInternalServerError, "Internal Server Error")
		return
	}

	var entry models.Entry

	err := models.DB.
		Where("list_name = ?", list.Name).
		First(&entry, id).
		Error

	if err != nil {
		c.String(http.StatusInternalServerError, "Internal Server Error")
		return
	}

	entry.Checked = checked[0] == 'f'
	err = models.DB.
		Save(&entry).
		Error

	if err != nil {
		c.String(http.StatusInternalServerError, "Internal Server Error")
		return
	}

	c.Redirect(http.StatusFound, models.BasePath() + "/")
}
