package repository

import (
	"gorm.io/gorm"

	"avito-shop-test/internal/models"
)

type PurchaseRepository interface {
	RecordPurchase(inventory *models.Inventory) error
	GetPurchasedItems(userID string) ([]models.Inventory, error)
}

type purchaseRepository struct {
	db *gorm.DB
}

func NewPurchaseRepository(db *gorm.DB) PurchaseRepository {
	return &purchaseRepository{db: db}
}

func (r *purchaseRepository) RecordPurchase(inventory *models.Inventory) error {
	if err := r.db.Create(&inventory).Error; err != nil {
		return err
	}
	return nil
}

func (r *purchaseRepository) GetPurchasedItems(userID string) ([]models.Inventory, error) {
	var purchases []models.Inventory
	err := r.db.Where("user_ID = ?", userID).Find(&purchases).Error
	if err != nil {
		return nil, err
	}
	return purchases, nil
}
