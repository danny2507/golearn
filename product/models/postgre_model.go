package models

import (
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Name        string `gorm:"not null"`
	Description string
	Price       float64 `gorm:"not null;check:price >= 0"`
	UserID      int
}
type ActiveToken struct {
	gorm.Model
	Token string `gorm:"not null;index:idx_active_token_token,unique"`
}
