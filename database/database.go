package database

import (
	"fmt"
	"log"

	"github.com/Ray0427/url-shortener/config"
	"github.com/Ray0427/url-shortener/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDatabase(config config.Config) *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@%s(%s:%s)/%s?charset=utf8mb4&parseTime=True", config.Database.USERNAME, config.Database.PASSWORD, config.Database.NETWORK, config.Database.SERVER, config.Database.PORT, config.Database.DATABASE)
	// log.Println(dsn)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Error connect Database with gorm")
	}
	db.AutoMigrate(&model.Url{})
	return db
}
