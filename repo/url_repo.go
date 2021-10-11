package repo

import (
	"time"

	"github.com/Ray0427/url-shortener/model"
	"gorm.io/gorm"
)

type UrlRepo struct {
	DB *gorm.DB
}

func NewUrlRepo(DB *gorm.DB) UrlRepo {
	return UrlRepo{DB: DB}
}

func (ur *UrlRepo) CreateUrl(url string, expireAt time.Time) (model.Url, error) {
	u := model.Url{FullUrl: url, ExpireAt: expireAt}
	err := ur.DB.Create(&u).Error
	return u, err
}

func (ur *UrlRepo) GetUrl(id uint) (model.Url, error) {
	var u model.Url
	err := ur.DB.First(&u, id).Error
	return u, err
}
