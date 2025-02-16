package repository

import (
	"github.com/stretchr/testify/mock"

	"avito-shop-test/internal/models"
)

type MockUserRepository struct {
	mock.Mock
}

// GetUserByUserID implements usecase.UserRepository.
func (m *MockUserRepository) GetUserByUserID(userID string) (*models.User, error) {
	panic("unimplemented")
}

func (m *MockUserRepository) FindUserByUsername(username string) (*models.User, error) {
	args := m.Called(username)

	if user, ok := args.Get(0).(*models.User); ok {
		return user, args.Error(1)
	}

	return nil, args.Error(1)
}

func (m *MockPurchaseRepository) GetUserByUserID(userID string) (*models.User, error) {
	args := m.Called(userID)

	if user, ok := args.Get(0).(*models.User); ok {
		return user, args.Error(1)
	}

	return nil, args.Error(1)
}

func (m *MockUserRepository) CreateUser(user *models.User) error {
	return m.Called(user).Error(0)
}

func (m *MockUserRepository) UpdateUserBalance(username string, amount int) error {
	return m.Called(username, amount).Error(0)
}
