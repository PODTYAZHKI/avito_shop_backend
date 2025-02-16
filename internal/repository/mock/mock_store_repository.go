package repository

import (
	"github.com/stretchr/testify/mock"

	"avito-shop-test/internal/models"
)

type MockStoreRepository struct {
	mock.Mock
}

func (m *MockStoreRepository) GetItemByName(name string) (*models.Product, error) {
	args := m.Called(name)

	if product, ok := args.Get(0).(*models.Product); ok {
		return product, args.Error(1)
	}

	return nil, args.Error(1)
}
