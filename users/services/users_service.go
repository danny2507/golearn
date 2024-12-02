package services

import (
	"errors"
	"fmt"
	"golearn/users/config"
	"golearn/users/models"
	"golearn/users/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

type PostgreService struct {
	Db  *gorm.DB
	Sdb *gorm.DB
}

func NewPostgreService(config *config.Config) *PostgreService {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.PostgreSQL.Host,
		config.PostgreSQL.Port,
		config.PostgreSQL.User,
		config.PostgreSQL.Password,
		config.PostgreSQL.Database)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto-migrate ensures schema matches the models
	db.AutoMigrate(&models.User{})

	dsn = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.PostgreSQL.Host,
		config.PostgreSQL.Port,
		config.PostgreSQL.User,
		config.PostgreSQL.Password,
		config.PostgreSQL.SharedDatabase)

	sdb, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Println("Database connected and initialized successfully.")

	return &PostgreService{Db: db, Sdb: sdb}
}

// RegisterUser registers a new user in the database
func (s *PostgreService) RegisterUser(username, email, password string) error {
	hashedPassword := utils.HashPassword(password)
	user := models.User{
		Username: username,
		Email:    email,
		Password: hashedPassword,
	}
	result := s.Db.Create(&user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// LoginUser authenticates a user and returns their ID if successful
func (s *PostgreService) LoginUser(email, password string) (int, error) {
	var user models.User
	result := s.Db.Where("email = ?", email).First(&user)
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
