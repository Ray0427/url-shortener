package model

import (
	"time"
)

type Url struct {
	ID       uint      `gorm:"primarykey"`
	FullUrl  string    `json:"fullUrl"`
	ExpireAt time.Time `json:"expireAt"`
}
