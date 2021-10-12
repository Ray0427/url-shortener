package model

import (
	"time"

	"gorm.io/gorm"
)

type Url struct {
	gorm.Model
	FullUrl  string    `json:"fullUrl"`
	ExpireAt time.Time `json:"expireAt"`
}
