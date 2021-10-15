package repo

import (
	"errors"
	"log"
	"time"

	"github.com/Ray0427/url-shortener/cache"
	"github.com/Ray0427/url-shortener/config"
	"github.com/Ray0427/url-shortener/dao"
	"github.com/Ray0427/url-shortener/utils"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type UrlRepo struct {
	Dao    dao.UrlDaoInterface
	cache  cache.CacheInterface
	config config.Config
}

func NewUrlRepo(dao dao.UrlDaoInterface, cache cache.CacheInterface, config config.Config) *UrlRepo {
	return &UrlRepo{
		Dao:    dao,
		cache:  cache,
		config: config,
	}
}

type NotFoundError struct {
	Message string
}

func (e *NotFoundError) Error() string {
	return "Not Found"
}

type BadRequestError struct {
	Message string
}

func (e *BadRequestError) Error() string {
	return e.Message
}

type InternalServerError struct {
	Message string
}

func (e *InternalServerError) Error() string {
	return e.Message
}

func (ur *UrlRepo) CreateUrl(url string, expireAt time.Time) (string, error) {
	urlDTO, err := ur.Dao.Create(url, expireAt)
	if err != nil {
		log.Printf("%+v\n", err)
		return "", &InternalServerError{
			Message: "DB Error",
		}
	}
	hashID, err := utils.Encode(ur.config.HashID.Salt, ur.config.HashID.MinLength, int(urlDTO.ID))
	if err != nil {
		log.Printf("%+v\n", err)
		return "", &BadRequestError{
			Message: "Invalid HashId",
		}
	}
	err = ur.cache.SetUrl(hashID, cache.Url{
		Url:      url,
		ExpireAt: expireAt,
	})
	if err != nil {
		log.Printf("%+v\n", err)
		return "", &InternalServerError{
			Message: "Cache Error",
		}
	}
	return hashID, err
}

func (ur *UrlRepo) GetUrl(hashID string) (string, error) {
	var url cache.Url
	err := ur.cache.GetUrl(hashID, &url)
	if err == redis.Nil {
		log.Println("cache not hit")
	} else if _, isEmptyError := err.(*cache.EmptyError); isEmptyError {
		return "", &NotFoundError{}
	} else if err != nil {
		log.Printf("%+v\n", err)
	} else if err == nil {
		if url.ExpireAt.Before(time.Now()) {
			return "", &NotFoundError{
				Message: "Link Expired",
			}
		}
		return url.Url, nil
	}
	id, err := utils.Decode(ur.config.HashID.Salt, ur.config.HashID.MinLength, hashID)
	if err != nil {
		return "", &BadRequestError{
			Message: "Invalid HashId",
		}
	}
	urlDTO, err := ur.Dao.Get(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ur.cache.SetUrl(hashID, nil)
			return "", &NotFoundError{}
		}
		return "", &InternalServerError{
			Message: "DB Error",
		}
	}
	err = ur.cache.SetUrl(hashID, cache.Url{
		Url:      urlDTO.FullUrl,
		ExpireAt: urlDTO.ExpireAt,
	})
	if err != nil {
		log.Printf("%+v\n", err)
	}
	if urlDTO.ExpireAt.Before(time.Now()) {
		return "", &NotFoundError{}
	}
	return urlDTO.FullUrl, err
}
