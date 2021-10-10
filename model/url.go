package model

import (
	"time"

	"gorm.io/gorm"
)

type Url struct {
	gorm.Model
	FullUrl  string
	ExpireAt time.Time
}
