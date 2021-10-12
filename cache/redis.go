package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Ray0427/url-shortener/config"
	"github.com/go-redis/redis/v8"
)

type Cache struct {
	client *redis.Client
	ctx    context.Context
}

func InitCache(config config.Config) *Cache {
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.Redis.Address + ":" + config.Redis.Port,
		Password: config.Redis.Password,
		DB:       0, // use default DB
	})
	ctx := context.Background()
	return &Cache{
		client: rdb,
		ctx:    ctx,
	}
}

func (c *Cache) SetUrl(hashId string, url interface{}) error {
	val, err := json.Marshal(url)
	if err != nil {
		return err
	}
	err = c.client.Set(c.ctx, "HASH_ID:"+hashId, val, time.Hour).Err()
	if err != nil {
		return err
	}
	return nil
}

func (c *Cache) GetUrl(hashId string, url interface{}) error {
	cacheVal, err := c.client.Get(c.ctx, "HASH_ID:"+hashId).Result()
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(cacheVal), &url)
	if err != nil {
		return err
	}
	return nil
}
