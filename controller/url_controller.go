package controller

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Ray0427/url-shortener/cache"
	"github.com/Ray0427/url-shortener/config"
	"github.com/Ray0427/url-shortener/repo"
	"github.com/gin-gonic/gin"
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
	hashID, err := uc.urlRepo.CreateUrl(body.Url, body.ExpireAt)
	if err != nil {
		switch err.(type) {
		case *repo.BadRequestError:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	s := fmt.Sprintf("http://localhost/%s", hashID)
	c.JSON(http.StatusOK, gin.H{
		"id":       hashID,
		"shortUrl": s,
	})
}

func (uc *urlController) GetId(c *gin.Context) {
	hashId := c.Param("url_id")

	url, err := uc.urlRepo.GetUrl(hashId)
	if err != nil {
		switch err.(type) {
		case *repo.BadRequestError:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case *repo.NotFoundError:
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case *repo.InternalServerError:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.Redirect(http.StatusMovedPermanently, url)
}
