package config

import (
	"fmt"
	"log"
	"os"

	"golearn/models"
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
	db.AutoMigrate(&models.User{}, &models.Product{})

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

	// Insert mock products
	products := []models.Product{
		{Name: "Laptop", Description: "High-performance laptop for gaming and work", Price: 1200.00, UserID: 1},
		{Name: "Smartphone", Description: "Latest model smartphone with all features", Price: 800.00, UserID: 2},
		{Name: "Headphones", Description: "Noise-canceling headphones", Price: 150.00, UserID: 1},
		{Name: "Tablet", Description: "Lightweight tablet for on-the-go productivity", Price: 300.00, UserID: 3},
		{Name: "Monitor", Description: "4K UHD monitor for immersive viewing", Price: 400.00, UserID: 2},
		{Name: "Keyboard", Description: "Mechanical keyboard with RGB lighting", Price: 90.00, UserID: 1},
	}

	for _, product := range products {
		if err := db.Create(&product).Error; err != nil {
			log.Printf("Failed to insert product %s: %v", product.Name, err)
		}
	}

	log.Println("Mock data inserted successfully.")
}
