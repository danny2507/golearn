package services

import (
	"errors"

	"golearn/models"
	"golearn/utils"
	"gorm.io/gorm"
)

func RegisterUser(db *gorm.DB, username, email, password string) error {
	hashedPassword := utils.HashPassword(password)
	user := models.User{
		Username: username,
		Email:    email,
		Password: hashedPassword,
	}
	result := db.Create(&user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func LoginUser(db *gorm.DB, email, password string) (int, error) {
	var user models.User
	result := db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return 0, errors.New("username/password invalid")
		}
		return 0, result.Error
	}

	if !utils.VerifyPasswordHash(password, user.Password) {
		return 0, errors.New("username/password invalid")
	}

	return int(user.ID), nil
}
