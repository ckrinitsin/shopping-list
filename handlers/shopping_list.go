package shopping_list

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func LoadElements(c *gin.Context) {
	c.HTML(http.StatusOK, "template.html", gin.H{
		"Name": "Shopping List",
	})
}
