package repository

import (
	"github.com/stretchr/testify/mock"

	"avito-shop-test/internal/models"
)

type MockPurchaseRepository struct {
	mock.Mock
}

func (m *MockPurchaseRepository) RecordPurchase(inventory *models.Inventory) error {
	args := m.Called(inventory)
	return args.Error(0)
}

func (m *MockPurchaseRepository) GetPurchasedItems(userID string) ([]models.Inventory, error) {
	args := m.Called(userID)

	if inventory, ok := args.Get(0).([]models.Inventory); ok {
		return inventory, args.Error(1)
	}

	return nil, args.Error(1)
}
