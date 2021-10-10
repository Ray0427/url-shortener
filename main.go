package main

import (
	"fmt"
	"log"

	"github.com/Ray0427/url-shortener/controller"
	"github.com/Ray0427/url-shortener/model"
	"github.com/Ray0427/url-shortener/repo"
	"github.com/caarlos0/env/v6"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Config struct {
	Database struct {
		USERNAME string `env:"DB_USERNAME"`
		PASSWORD string `env:"DB_PASSWORD"`
		NETWORK  string `env:"DB_NETWORK" envDefault:"tcp"`
		SERVER   string `env:"DB_ADDRESS"`
		PORT     string `env:"DB_PORT" envDefault:"3306"`
		DATABASE string `env:"DB_DATABASE"`
	}
	Server struct {
		PORT string `env:"PORT"`
	}
}

func initConfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	config := Config{}
	if err := env.Parse(&config); err != nil {
		log.Fatalf("%+v\n", err)
	}
	// log.Printf("%+v\n", config)
	return config
}

func initRouter(db *gorm.DB) *gin.Engine {
	urlRepo := repo.NewUrlRepo(db)
	urlController := controller.NewUrlController(urlRepo)
	r := gin.Default()
	r.POST("/api/v1/urls", urlController.PostUrl)
	r.GET("/:url_id", urlController.GetId)
	r.Run()
	return r
}

func initDatabase(config Config) *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@%s(%s:%s)/%s?charset=utf8mb4&parseTime=True", config.Database.USERNAME, config.Database.PASSWORD, config.Database.NETWORK, config.Database.SERVER, config.Database.PORT, config.Database.DATABASE)
	// log.Println(dsn)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Error connect Database with gorm")
	}
	db.AutoMigrate(&model.Url{})
	return db
}

func main() {
	config := initConfig()
	db := initDatabase(config)
	sqlDB, err := db.DB()
	if err != nil {
		log.Printf("%+v\n", err)
	}
	defer sqlDB.Close()
	initRouter(db)
}
