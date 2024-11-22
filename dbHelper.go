package main

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

type User struct {
	gorm.Model
	Username string `gorm:"uniqueIndex;not null"`
	Email    string `gorm:"uniqueIndex;not null"`
	Password string `gorm:"not null"`
	Products []Product
}

type Product struct {
	gorm.Model
	Name        string `gorm:"not null"`
	Description string
	Price       float64 `gorm:"not null;check:price >= 0"`
	UserID      int
}

// CreateMockData creates mock users and products in the database
func CreateMockData() error {
	// Create mock users
	users := []User{
		{Username: "alice", Email: "alice@example.com", Password: "$2a$10$pD8kOZPLyA7.b7/Afd1sRedxQzT8v.PJrRNJc4pKwEKI.SvM8LzkW"},
		{Username: "bob", Email: "bob@example.com", Password: "$2a$10$pD8kOZPLyA7.b7/Afd1sRedxQzT8v.PJrRNJc4pKwEKI.SvM8LzkW"},
		{Username: "charlie", Email: "charlie@example.com", Password: "$2a$10$pD8kOZPLyA7.b7/Afd1sRedxQzT8v.PJrRNJc4pKwEKI.SvM8LzkW"},
	}

	for _, user := range users {
		result := db.Create(&user)
		if result.Error != nil {
			return fmt.Errorf("could not create mock user: %v", result.Error)
		}
	}

	// Create mock products
	products := []Product{
		{Name: "Laptop", Description: "High-performance laptop for gaming and work", Price: 1200.00, UserID: 1},
		{Name: "Smartphone", Description: "Latest model smartphone with all features", Price: 800.00, UserID: 2},
		{Name: "Headphones", Description: "Noise-canceling headphones", Price: 150.00, UserID: 1},
		{Name: "Tablet", Description: "Lightweight tablet for on-the-go productivity", Price: 300.00, UserID: 3},
		{Name: "Monitor", Description: "4K UHD monitor for immersive viewing", Price: 400.00, UserID: 2},
		{Name: "Keyboard", Description: "Mechanical keyboard with RGB lighting", Price: 90.00, UserID: 1},
	}

	for _, product := range products {
		result := db.Create(&product)
		if result.Error != nil {
			return fmt.Errorf("could not create mock product: %v", result.Error)
		}
	}

	return nil
}

// InitDB initializes the database connection and sets up tables
func InitDB() *gorm.DB {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"))

	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto Migrate the schema
	db.AutoMigrate(&User{}, &Product{})
	var count int64
	db.Model(&User{}).Count(&count)
	if count == 0 {
		// If the database is empty, create mock data
		err := CreateMockData()
		if err != nil {
			log.Printf("Failed to create mock data: %v", err)
		} else {
			log.Println("Mock data created successfully")
		}
	}
	return db
}

// CloseDB closes the database connection
func CloseDB() {
	sqlDB, _ := db.DB()
	if sqlDB != nil {
		sqlDB.Close()
	}
}

// AddProduct inserts a new product into the database
func AddProduct(name, description string, price float64, userID int) (*Product, error) {
	product := &Product{
		Name:        name,
		Description: description,
		Price:       price,
		UserID:      userID,
	}
	result := db.Create(product)
	if result.Error != nil {
		return nil, fmt.Errorf("could not insert product: %v", result.Error)
	}
	return product, nil
}

// GetProduct retrieves a single product by ID
func GetProduct(id int) (*Product, error) {
	var product Product
	result := db.First(&product, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil // Product not found
		}
		return nil, fmt.Errorf("could not get product: %v", result.Error)
	}
	return &product, nil
}

// ListProducts retrieves all products from the database
func ListProducts() ([]Product, error) {
	var products []Product
	result := db.Find(&products)
	if result.Error != nil {
		return nil, fmt.Errorf("could not get products: %v", result.Error)
	}
	return products, nil
}

// UpdateProduct updates an existing product in the database
func UpdateProduct(id int, name, description string, price float64) error {
	result := db.Model(&Product{}).Where("id = ?", id).Updates(Product{
		Name:        name,
		Description: description,
		Price:       price,
	})
	if result.Error != nil {
		return fmt.Errorf("could not update product: %v", result.Error)
	}
	return nil
}

// DeleteProduct removes a product from the database
func DeleteProduct(id int) error {
	result := db.Delete(&Product{}, id)
	if result.Error != nil {
		return fmt.Errorf("could not delete product: %v", result.Error)
	}
	return nil
}

// LoginUser checks user credentials and returns user ID if valid
func LoginUser(email, password string) (int, error) {
	var user User
	result := db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return 0, fmt.Errorf("username/password invalid")
		}
		return 0, fmt.Errorf("error retrieving information: %v", result.Error)
	}

	if !VerifyPasswordHash(password, user.Password) {
		return 0, fmt.Errorf("username/password invalid")
	}

	return int(user.ID), nil
}

// RegisterUser creates a new user in the database
func RegisterUser(username, email, password string) error {
	hashedPassword := HashPassword(password)
	user := User{
		Username: username,
		Email:    email,
		Password: hashedPassword,
	}
	result := db.Create(&user)
	if result.Error != nil {
		return fmt.Errorf("could not insert user: %v", result.Error)
	}
	return nil
}
