package controller

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Ray0427/url-shortener/cache"
	"github.com/Ray0427/url-shortener/config"
	"github.com/Ray0427/url-shortener/repo"
	"github.com/Ray0427/url-shortener/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type UrlController interface {
	PostUrl(c *gin.Context)
	GetId(c *gin.Context)
}

type urlController struct {
	config  config.Config
	urlRepo repo.UrlRepo
	cache   *cache.Cache
}

//Constructor Function
func NewUrlController(config config.Config, repo repo.UrlRepo, cache *cache.Cache) UrlController {
	return &urlController{
		config:  config,
		urlRepo: repo,
		cache:   cache,
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
	urlDTO, err := uc.urlRepo.CreateUrl(body.Url, body.ExpireAt)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	hashID := utils.Encode(uc.config.HashID.Salt, uc.config.HashID.MinLength, int(urlDTO.ID))
	err = uc.cache.SetUrl(hashID, body)
	if err != nil {
		log.Printf("%+v\n", err)
	}
	s := fmt.Sprintf("http://localhost/%s", hashID)
	c.JSON(http.StatusOK, gin.H{
		"id":       hashID,
		"shortUrl": s,
	})
}

func (uc *urlController) GetId(c *gin.Context) {
	hashId := c.Param("url_id")
	var postUrl PostUrlsParam
	err := uc.cache.GetUrl(hashId, postUrl)
	if err == redis.Nil {
		log.Println("cache not found")
	} else if postUrl.Url == "" {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Not found",
		})
		return
	} else if err != nil {
		log.Printf("%+v\n", err)
	}
	id, err := utils.Decode(uc.config.HashID.Salt, uc.config.HashID.MinLength, hashId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	urlDTO, err := uc.urlRepo.GetUrl(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			uc.cache.SetUrl(hashId, nil)
		}
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.Redirect(http.StatusMovedPermanently, urlDTO.FullUrl)
}
