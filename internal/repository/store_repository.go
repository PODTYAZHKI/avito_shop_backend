package repository

import (
	"gorm.io/gorm"

	"avito-shop-test/internal/models"
)

type StoreRepository interface {
	GetItemByName(name string) (*models.Product, error)
}

type storeRepository struct {
	db *gorm.DB
}

func NewStoreRepository(db *gorm.DB) StoreRepository {
	return &storeRepository{db: db}
}

func (r *storeRepository) GetItemByName(name string) (*models.Product, error) {
	var item models.Product
	err := r.db.Where("name = ?", name).First(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}
