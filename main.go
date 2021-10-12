package main

import (
	"log"

	"github.com/Ray0427/url-shortener/cache"
	"github.com/Ray0427/url-shortener/config"
	"github.com/Ray0427/url-shortener/controller"
	"github.com/Ray0427/url-shortener/database"
	"github.com/Ray0427/url-shortener/repo"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func initRouter(config config.Config, db *gorm.DB, cache *cache.Cache) *gin.Engine {
	urlRepo := repo.NewUrlRepo(db)
	urlController := controller.NewUrlController(config, urlRepo, cache)
	r := gin.Default()
	r.POST("/api/v1/urls", urlController.PostUrl)
	r.GET("/:url_id", urlController.GetId)
	r.Run()
	return r
}

func main() {
	config := config.InitConfig()
	db := database.InitDatabase(config)
	sqlDB, err := db.DB()
	if err != nil {
		log.Printf("%+v\n", err)
	}
	defer sqlDB.Close()
	cache := cache.InitCache(config)
	initRouter(config, db, cache)
}
