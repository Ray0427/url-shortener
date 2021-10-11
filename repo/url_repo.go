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

func (ur *UrlRepo) CreateUrl(url string, expireAt time.Time) model.Url {
	u := model.Url{FullUrl: url, ExpireAt: expireAt}
	ur.DB.Create(&u)
	return u
}

func (ur *UrlRepo) GetUrl(id uint) model.Url {
	var u model.Url
	ur.DB.First(&u, id)
	return u
}
