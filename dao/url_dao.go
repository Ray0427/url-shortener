package dao

import (
	"time"

	"github.com/Ray0427/url-shortener/model"
	"gorm.io/gorm"
)

type UrlDao struct {
	DB *gorm.DB
}

func NewUrlDao(DB *gorm.DB) UrlDao {
	return UrlDao{
		DB: DB,
	}
}

func (d *UrlDao) Create(url string, expireAt time.Time) (model.Url, error) {
	u := model.Url{FullUrl: url, ExpireAt: expireAt}
	err := d.DB.Create(&u).Error
	return u, err
}

func (d *UrlDao) Get(id uint) (model.Url, error) {
	var u model.Url
	err := d.DB.First(&u, id).Error
	return u, err
}
