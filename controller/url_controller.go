package controller

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Ray0427/url-shortener/repo"
	"github.com/gin-gonic/gin"
)

type UrlController interface {
	PostUrl(c *gin.Context)
	GetId(c *gin.Context)
}

type urlController struct {
	urlRepo repo.UrlRepo
}

//Constructor Function
func NewUrlController(repo repo.UrlRepo) UrlController {
	return &urlController{
		urlRepo: repo,
	}
}

type PostUrlsParam struct {
	Url      string    `json:"url" binding:"required"`
	ExpireAt time.Time `json:"expireAt" binding:"required"`
}

func (uc *urlController) PostUrl(c *gin.Context) {
	var body PostUrlsParam
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	urlDTO := uc.urlRepo.CreateUrl(body.Url, body.ExpireAt)
	s := fmt.Sprintf("http://localhost/%d", urlDTO.ID)
	c.JSON(http.StatusOK, gin.H{
		"id":       urlDTO.ID,
		"shortUrl": s,
	})
}

func (uc *urlController) GetId(c *gin.Context) {
	c.Redirect(http.StatusMovedPermanently, "http://localhost/<url_id>")
}
