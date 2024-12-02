package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"uniqueIndex;not null"`
	Email    string `gorm:"uniqueIndex;not null"`
	Password string `gorm:"not null"`
}
type ActiveToken struct {
	gorm.Model
	Token string `gorm:"not null;index:idx_active_token_token,unique"`
}
