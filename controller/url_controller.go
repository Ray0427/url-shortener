package controller

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Ray0427/url-shortener/repo"
	"github.com/Ray0427/url-shortener/utils"
	"github.com/gin-gonic/gin"
)

type UrlController interface {
	PostUrl(c *gin.Context)
	GetId(c *gin.Context)
}

type urlController struct {
	urlRepo repo.UrlRepoInterface
}

//Constructor Function
func NewUrlController(repo repo.UrlRepoInterface) UrlController {
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
	if !utils.CheckUrl(body.Url) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid url"})
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
