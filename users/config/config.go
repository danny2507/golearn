package config

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"log"
	"os"

	"golearn/users/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// InitDB initializes the database connection and inserts mock data
func InitDB() *gorm.DB {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"))

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto-migrate ensures schema matches the models
	db.AutoMigrate(&models.User{})

	// Insert mock data if tables are empty
	insertMockData(db)

	log.Println("Database connected and initialized successfully.")
	return db
}

// insertMockData inserts mock users and products into the database
func insertMockData(db *gorm.DB) {
	// Check if there are any users in the database
	var userCount int64
	db.Model(&models.User{}).Count(&userCount)
	if userCount > 0 {
		log.Println("Mock data already exists. Skipping insertion.")
		return
	}

	// Insert mock users
	users := []models.User{
		{Username: "alice", Email: "alice@example.com", Password: "$2a$10$pD8kOZPLyA7.b7/Afd1sRedxQzT8v.PJrRNJc4pKwEKI.SvM8LzkW"},
		{Username: "bob", Email: "bob@example.com", Password: "$2a$10$pD8kOZPLyA7.b7/Afd1sRedxQzT8v.PJrRNJc4pKwEKI.SvM8LzkW"},
		{Username: "charlie", Email: "charlie@example.com", Password: "$2a$10$pD8kOZPLyA7.b7/Afd1sRedxQzT8v.PJrRNJc4pKwEKI.SvM8LzkW"},
	}

	for _, user := range users {
		if err := db.Create(&user).Error; err != nil {
			log.Printf("Failed to insert user %s: %v", user.Username, err)
		}
	}

	log.Println("Mock data inserted successfully.")
}
func InitSharedDB() *gorm.DB {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("SHARED_DB_HOST"),
		os.Getenv("SHARED_DB_PORT"),
		os.Getenv("SHARED_DB_USER"),
		os.Getenv("SHARED_DB_PASSWORD"),
		os.Getenv("SHARED_DB_NAME"))

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to shared database: %v", err)
	}

	// Auto-migrate ensures schema matches the models
	db.AutoMigrate(&models.ActiveToken{})

	log.Println("Shared Database connected and initialized successfully.")
	return db
}

type Config struct {
	PostgreSQL struct {
		Host           string `envconfig:"DB_HOST"`
		Port           int    `envconfig:"DB_PORT"`
		User           string `envconfig:"DB_USER"`
		Password       string `envconfig:"DB_PASSWORD"`
		Database       string `envconfig:"DB_NAME"`
		SharedDatabase string `envconfig:"SHARED_DB_NAME"`
	}

	SignatureRequest struct {
		SecretKey string `envconfig:"SIGNATURE_REQUEST_SECRET_KEY"`
	}
	Debug struct {
		IsProfiling bool `envconfig:"DEBUG_IS_PROFILING"`
	}
}

func LoadConfig() *Config {

	var config Config

	// Override with environment variables
	err := envconfig.Process("", &config)
	if err != nil {
		fmt.Printf("Error processing environment variables: %s\n", err)
		return nil
	}

	// Use config as needed
	fmt.Println("Signature Request Secret Key:", config.SignatureRequest.SecretKey)
	fmt.Println("Debug Is Profiling:", config.Debug.IsProfiling)

	return &config
}