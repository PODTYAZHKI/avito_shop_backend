package usecase

import (
	"github.com/stretchr/testify/mock"

	"avito-shop-test/internal/models"
)

type MockUserUC struct {
	mock.Mock
}

func (m *MockUserUC) Authenticate(username, password string) (string, error) {
	args := m.Called(username, password)

	return args.String(0), args.Error(1)
}

func (m *MockUserUC) GetUserInfo(username string) (*models.UserInfo, error) {
	args := m.Called(username)

	return args.Get(0).(*models.UserInfo), args.Error(1)
}

func (m *MockUserUC) GetCoinHistory(userID string) (models.CoinHistory, error) {
	args := m.Called(userID)

	return args.Get(0).(models.CoinHistory), args.Error(1)
}

func (m *MockUserUC) GetPurchasedItems(userID string) ([]models.PurchasedItem, error) {
	args := m.Called(userID)

	return args.Get(0).([]models.PurchasedItem), args.Error(1)
}
