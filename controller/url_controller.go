package controller

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Ray0427/url-shortener/config"
	"github.com/Ray0427/url-shortener/repo"
	"github.com/Ray0427/url-shortener/utils"
	"github.com/gin-gonic/gin"
)

type UrlController interface {
	PostUrl(c *gin.Context)
	GetId(c *gin.Context)
}

type urlController struct {
	urlRepo repo.UrlRepo
	config  config.Config
}

//Constructor Function
func NewUrlController(repo repo.UrlRepo, config config.Config) UrlController {
	return &urlController{
		urlRepo: repo,
		config:  config,
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
	hashID := utils.Encode(uc.config.HashID.Salt, uc.config.HashID.MinLength, int(urlDTO.ID))
	s := fmt.Sprintf("http://localhost/%s", hashID)
	c.JSON(http.StatusOK, gin.H{
		"id":       hashID,
		"shortUrl": s,
	})
}

func (uc *urlController) GetId(c *gin.Context) {
	hashId := c.Param("url_id")
	id := utils.Decode(uc.config.HashID.Salt, uc.config.HashID.MinLength, hashId)
	urlDTO := uc.urlRepo.GetUrl(uint(id))
	c.Redirect(http.StatusMovedPermanently, urlDTO.FullUrl)
}
