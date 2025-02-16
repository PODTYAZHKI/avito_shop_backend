package repository

import (
	"github.com/stretchr/testify/mock"

	"avito-shop-test/internal/models"
)

type MockTransactionRepository struct {
	mock.Mock
}

func (m *MockTransactionRepository) RecordTransaction(transaction *models.CoinTransaction) error {
	args := m.Called(transaction)
	return args.Error(0)
}

func (m *MockTransactionRepository) GetTransactionsHistory(userID string) ([]models.CoinTransaction, error) {
	args := m.Called(userID)

	if transactions, ok := args.Get(0).([]models.CoinTransaction); ok {
		return transactions, args.Error(1)
	}

	return nil, args.Error(1)
}
