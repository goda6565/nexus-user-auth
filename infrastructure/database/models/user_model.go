package models

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	ObjID           string `gorm:"type:uuid;uniqueIndex;not null"` // 外部識別用のUUID
	Email           string `gorm:"size:255;uniqueIndex;not null"`
	Password        string `gorm:"size:255;not null"`
	Username        string `gorm:"size:255;not null"`
	AvatarURL       string `gorm:"size:255"`
	EmailVerifiedAt *time.Time
	LastLoginAt     *time.Time
	Role            string `gorm:"size:50;default:'user';not null"`
}
