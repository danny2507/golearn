package services

import (
	"golearn/product/models"
	"gorm.io/gorm"
)

type SharedPostgreService struct {
	Db *gorm.DB
}

func NewSharedPostgreService(db *gorm.DB) *SharedPostgreService {
	return &SharedPostgreService{Db: db}
}
func (sps *SharedPostgreService) GetActiveToken(token string) (*models.ActiveToken, error) {
	var activeToken models.ActiveToken
	result := sps.Db.Where(&models.ActiveToken{Token: token}).First(&activeToken)
	if result.Error != nil {
		return nil, result.Error
	}
	return &activeToken, nil
}
func (sps *SharedPostgreService) AddActiveToken(token string) (*models.ActiveToken, error) {
	tokenRecord := &models.ActiveToken{Token: token}
	result := sps.Db.Create(tokenRecord)
	if result.Error != nil {
		return nil, result.Error
	}
	return tokenRecord, nil
}
