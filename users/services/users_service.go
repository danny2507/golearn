package services

import (
	"errors"
	"golearn/users/models"
	"golearn/users/utils"
	"gorm.io/gorm"
)

// UserService struct encapsulates the database connection
type UserService struct {
	DB *gorm.DB
}

// NewUserService creates a new instance of UserService with the given DB connection
func NewUserService(db *gorm.DB) *UserService {
	return &UserService{DB: db}
}

// RegisterUser registers a new user in the database
func (s *UserService) RegisterUser(username, email, password string) error {
	hashedPassword := utils.HashPassword(password)
	user := models.User{
		Username: username,
		Email:    email,
		Password: hashedPassword,
	}
	result := s.DB.Create(&user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// LoginUser authenticates a user and returns their ID if successful
func (s *UserService) LoginUser(email, password string) (int, error) {
	var user models.User
	result := s.DB.Where("email = ?", email).First(&user)
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
