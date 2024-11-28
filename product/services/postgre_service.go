package services

import (
	"errors"

	"golearn/product/models"
	"gorm.io/gorm"
)

type PostgreService struct {
	Db *gorm.DB
}

func NewPostgreService(db *gorm.DB) *PostgreService {
	return &PostgreService{Db: db}
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
