package services

import (
	"errors"
	"fmt"
	"golearn/product/config"
	"golearn/product/models"
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
	db.AutoMigrate(&models.Product{})

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
	sdb.AutoMigrate(&models.ActiveToken{})

	log.Println("Database connected and initialized successfully.")

	return &PostgreService{Db: db, Sdb: sdb}
}

func (ps *PostgreService) AddProduct(name, description string, price float64, userID int) (*models.Product, error) {
	product := &models.Product{
		Name:        name,
		Description: description,
		Price:       price,
		UserID:      userID,
	}
	result := ps.Db.Create(product)
	if result.Error != nil {
		return nil, result.Error
	}
	return product, nil
}

func (ps *PostgreService) GetProduct(id int) (*models.Product, error) {
	var product models.Product
	result := ps.Db.First(&product, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil // Product not found
		}
		return nil, result.Error
	}
	return &product, nil
}

func (ps *PostgreService) ListProducts() ([]models.Product, error) {
	var products []models.Product
	result := ps.Db.Find(&products)
	if result.Error != nil {
		return nil, result.Error
	}
	return products, nil
}

func (ps *PostgreService) UpdateProduct(id int, name, description string, price float64) error {
	result := ps.Db.Model(&models.Product{}).Where("id = ?", id).Updates(models.Product{
		Name:        name,
		Description: description,
		Price:       price,
	})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (ps *PostgreService) DeleteProduct(id int) error {
	result := ps.Db.Delete(&models.Product{}, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (ps *PostgreService) GetActiveToken(token string) (*models.ActiveToken, error) {
	var activeToken models.ActiveToken
	result := ps.Sdb.Where(&models.ActiveToken{Token: token}).First(&activeToken)
	if result.Error != nil {
		return nil, result.Error
	}
	return &activeToken, nil
}
func (ps *PostgreService) AddActiveToken(token string) (*models.ActiveToken, error) {
	tokenRecord := &models.ActiveToken{Token: token}
	result := ps.Sdb.Create(tokenRecord)
	if result.Error != nil {
		return nil, result.Error
	}
	return tokenRecord, nil
}

func (ps *PostgreService) insertMockData() {
	// Check if there are any products in the database
	var productCount int64
	ps.Db.Model(&models.Product{}).Count(&productCount)
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
		if err := ps.Db.Create(&product).Error; err != nil {
			log.Printf("Failed to insert product %s: %v", product.Name, err)
		}
	}

	log.Println("Mock data inserted successfully.")
}
