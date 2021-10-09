package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetId(c *gin.Context) {
	c.Redirect(http.StatusMovedPermanently, "http://localhost/<url_id>")
}
