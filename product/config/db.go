package config

import (
	"fmt"
	"log"
	"os"

	"golearn/product/models"
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
	db.AutoMigrate(&models.Product{})

	// Insert mock data if tables are empty
	insertMockData(db)

	log.Println("Database connected and initialized successfully.")
	return db
}

// insertMockData inserts mock users and products into the database
func insertMockData(db *gorm.DB) {
	// Check if there are any products in the database
	var productCount int64
	db.Model(&models.Product{}).Count(&productCount)
	if productCount > 0 {
		log.Println("Mock data already exists. Skipping insertion.")
		return
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
