package main

import (
	"github.com/Ray0427/url-shortener/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.POST("/api/v1/urls", handlers.PostUrls)
	r.GET("/:url_id", handlers.GetId)
	r.Run()
}
