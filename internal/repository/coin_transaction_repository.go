package repository

import (
	"gorm.io/gorm"

	"avito-shop-test/internal/models"
)

type CoinTransactionRepository interface {
	RecordTransaction(transaction *models.CoinTransaction) error
	GetTransactionsHistory(userID string) ([]models.CoinTransaction, error)
}

type coinTransactionRepository struct {
	db *gorm.DB
}

func NewCoinTransactionRepository(db *gorm.DB) CoinTransactionRepository {
	return &coinTransactionRepository{db: db}
}

func (r *coinTransactionRepository) RecordTransaction(transaction *models.CoinTransaction) error {
	return r.db.Create(transaction).Error
}

func (r *coinTransactionRepository) GetTransactionsHistory(userID string) ([]models.CoinTransaction, error) {
	var transactions []models.CoinTransaction
	err := r.db.Where("from_user_id = ? OR to_user_id = ?", userID, userID).Find(&transactions).Error
	return transactions, err
}
