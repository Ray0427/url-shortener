package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func PostUrls(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"id":       "id",
		"shortUrl": "http://localhost/<url_id>",
	})
}
