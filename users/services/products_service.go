package services

import (
	"errors"

	"golearn/models"
	"gorm.io/gorm"
)

func AddProduct(db *gorm.DB, name, description string, price float64, userID int) (*models.Product, error) {
	product := &models.Product{
		Name:        name,
		Description: description,
		Price:       price,
		UserID:      userID,
	}
	result := db.Create(product)
	if result.Error != nil {
		return nil, result.Error
	}
	return product, nil
}

func GetProduct(db *gorm.DB, id int) (*models.Product, error) {
	var product models.Product
	result := db.First(&product, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil // Product not found
		}
		return nil, result.Error
	}
	return &product, nil
}

func ListProducts(db *gorm.DB) ([]models.Product, error) {
	var products []models.Product
	result := db.Find(&products)
	if result.Error != nil {
		return nil, result.Error
	}
	return products, nil
}

func UpdateProduct(db *gorm.DB, id int, name, description string, price float64) error {
	result := db.Model(&models.Product{}).Where("id = ?", id).Updates(models.Product{
		Name:        name,
		Description: description,
		Price:       price,
	})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func DeleteProduct(db *gorm.DB, id int) error {
	result := db.Delete(&models.Product{}, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
