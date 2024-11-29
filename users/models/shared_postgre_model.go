package models

import "gorm.io/gorm"

type ActiveToken struct {
	gorm.Model
	Token string `gorm:"not null;index:idx_active_token_token,unique"`
}
