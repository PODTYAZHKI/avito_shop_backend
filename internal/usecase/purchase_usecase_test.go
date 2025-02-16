package usecase

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"avito-shop-test/internal/models"
	mockRepo "avito-shop-test/internal/repository/mock"
)

func TestBuyItem_Success(t *testing.T) {
	mockUserRepo := new(mockRepo.MockUserRepository)
	mockPurchaseRepo := new(mockRepo.MockPurchaseRepository)
	mockStoreRepo := new(mockRepo.MockStoreRepository)
	uc := NewPurchaseUseCase(mockPurchaseRepo, mockUserRepo, mockStoreRepo)

	user := &models.User{ID: "user1", Balance: 100}
	product := &models.Product{Name: "item1", Price: 50}
	mockUserRepo.On("FindUserByUsername", "user1").Return(user, nil)
	mockStoreRepo.On("GetItemByName", "item1").Return(product, nil)

	inventory := &models.Inventory{
		UserID:   user.ID,
		ItemType: "item1",
		Quantity: 1,
	}
	mockUserRepo.On("UpdateUserBalance", "user1", -50).Return(nil)
	mockPurchaseRepo.On("RecordPurchase", inventory).Return(nil)

	err := uc.BuyItem("user1", "item1")

	assert.NoError(t, err)

	mockUserRepo.AssertExpectations(t)
	mockPurchaseRepo.AssertExpectations(t)
}

func TestBuyItem_UserNotFound(t *testing.T) {
	mockUserRepo := new(mockRepo.MockUserRepository)
	mockPurchaseRepo := new(mockRepo.MockPurchaseRepository)
	mockStoreRepo := new(mockRepo.MockStoreRepository)
	uc := NewPurchaseUseCase(mockPurchaseRepo, mockUserRepo, mockStoreRepo)

	mockUserRepo.On("FindUserByUsername", "user1").Return(nil, nil)

	err := uc.BuyItem("user1", "item1")

	assert.Error(t, err)
	assert.Equal(t, "пользователь не найден", err.Error())
	mockUserRepo.AssertExpectations(t)
}

func TestBuyItem_ProductNotFound(t *testing.T) {
	mockUserRepo := new(mockRepo.MockUserRepository)
	mockPurchaseRepo := new(mockRepo.MockPurchaseRepository)
	mockStoreRepo := new(mockRepo.MockStoreRepository)
	uc := NewPurchaseUseCase(mockPurchaseRepo, mockUserRepo, mockStoreRepo)

	user := &models.User{ID: "user1", Balance: 100}
	mockUserRepo.On("FindUserByUsername", "user1").Return(user, nil)
	mockStoreRepo.On("GetItemByName", "itemNotExists").Return(nil, errors.New("database error"))

	err := uc.BuyItem("user1", "itemNotExists")

	assert.Error(t, err)
	assert.Equal(t, "товар не найден", err.Error())
	mockUserRepo.AssertExpectations(t)
}

func TestBuyItem_InsufficientBalance(t *testing.T) {
	mockUserRepo := new(mockRepo.MockUserRepository)
	mockPurchaseRepo := new(mockRepo.MockPurchaseRepository)
	mockStoreRepo := new(mockRepo.MockStoreRepository)
	uc := NewPurchaseUseCase(mockPurchaseRepo, mockUserRepo, mockStoreRepo)

	user := &models.User{ID: "user1", Balance: 30}
	mockUserRepo.On("FindUserByUsername", "user1").Return(user, nil)
	product := &models.Product{Name: "item1", Price: 50}
	mockStoreRepo.On("GetItemByName", "item1").Return(product, nil)

	err := uc.BuyItem("user1", "item1")

	assert.Error(t, err)
	assert.Equal(t, "недостаточно монет для покупки", err.Error())
	mockUserRepo.AssertExpectations(t)
}

func TestBuyItem_UpdateBalanceError(t *testing.T) {
	mockUserRepo := new(mockRepo.MockUserRepository)
	mockPurchaseRepo := new(mockRepo.MockPurchaseRepository)
	mockStoreRepo := new(mockRepo.MockStoreRepository)
	uc := NewPurchaseUseCase(mockPurchaseRepo, mockUserRepo, mockStoreRepo)

	user := &models.User{ID: "user1", Balance: 100}
	mockUserRepo.On("FindUserByUsername", "user1").Return(user, nil)
	mockUserRepo.On("UpdateUserBalance", "user1", -50).Return(errors.New("update balance error"))

	product := &models.Product{Name: "item1", Price: 50}
	mockStoreRepo.On("GetItemByName", "item1").Return(product, nil)

	err := uc.BuyItem("user1", "item1")

	assert.Error(t, err)
	assert.Equal(t, "update balance error", err.Error())
	mockUserRepo.AssertExpectations(t)
}

func TestBuyItem_RecordPurchaseError(t *testing.T) {
	mockUserRepo := new(mockRepo.MockUserRepository)
	mockPurchaseRepo := new(mockRepo.MockPurchaseRepository)
	mockStoreRepo := new(mockRepo.MockStoreRepository)
	uc := NewPurchaseUseCase(mockPurchaseRepo, mockUserRepo, mockStoreRepo)

	user := &models.User{ID: "user1", Balance: 100}
	mockUserRepo.On("FindUserByUsername", "user1").Return(user, nil)
	mockUserRepo.On("UpdateUserBalance", "user1", -50).Return(nil)

	product := &models.Product{Name: "item1", Price: 50}
	mockStoreRepo.On("GetItemByName", "item1").Return(product, nil)

	inventory := &models.Inventory{
		UserID:   user.ID,
		ItemType: "item1",
		Quantity: 1,
	}
	mockPurchaseRepo.On("RecordPurchase", inventory).Return(errors.New("record purchase error"))

	err := uc.BuyItem("user1", "item1")

	assert.Error(t, err)
	assert.Equal(t, "record purchase error", err.Error())
	mockUserRepo.AssertExpectations(t)
}
