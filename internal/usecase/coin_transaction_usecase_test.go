package usecase

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"avito-shop-test/internal/models"
	mockRepo "avito-shop-test/internal/repository/mock"
)

func TestSendCoins_Success(t *testing.T) {
	mockUserRepo := new(mockRepo.MockUserRepository)
	mockTransactionRepo := new(mockRepo.MockTransactionRepository)
	uc := NewCoinTransactionUseCase(mockTransactionRepo, mockUserRepo)

	userFrom := &models.User{ID: "user1", Balance: 100}
	userTo := &models.User{ID: "user2", Balance: 50}

	mockUserRepo.On("FindUserByUsername", "user1").Return(userFrom, nil)
	mockUserRepo.On("FindUserByUsername", "user2").Return(userTo, nil)
	mockTransactionRepo.On("RecordTransaction", mock.Anything).Return(nil)
	mockUserRepo.On("UpdateUserBalance", "user1", -50).Return(nil)
	mockUserRepo.On("UpdateUserBalance", "user2", 50).Return(nil)

	err := uc.SendCoins("user1", "user2", 50)

	assert.NoError(t, err)
	mockUserRepo.AssertExpectations(t)
	mockTransactionRepo.AssertExpectations(t)
}

func TestSendCoins_UserNotFound(t *testing.T) {
	mockUserRepo := new(mockRepo.MockUserRepository)
	mockTransactionRepo := new(mockRepo.MockTransactionRepository)
	uc := NewCoinTransactionUseCase(mockTransactionRepo, mockUserRepo)

	mockUserRepo.On("FindUserByUsername", "user1").Return(nil, nil)

	err := uc.SendCoins("user1", "user2", 50)

	assert.Error(t, err)
	assert.Equal(t, "отправитель не найден", err.Error())
	mockUserRepo.AssertExpectations(t)
}

func TestSendCoins_ReceiverNotFound(t *testing.T) {
	mockUserRepo := new(mockRepo.MockUserRepository)
	mockTransactionRepo := new(mockRepo.MockTransactionRepository)
	uc := NewCoinTransactionUseCase(mockTransactionRepo, mockUserRepo)

	userFrom := &models.User{ID: "user1", Balance: 100}
	mockUserRepo.On("FindUserByUsername", "user1").Return(userFrom, nil)
	mockUserRepo.On("FindUserByUsername", "user2").Return(nil, nil)

	err := uc.SendCoins("user1", "user2", 50)

	assert.Error(t, err)
	assert.Equal(t, "получатель не найден", err.Error())
	mockUserRepo.AssertExpectations(t)
}

func TestSendCoins_InsufficientBalance(t *testing.T) {
	mockUserRepo := new(mockRepo.MockUserRepository)
	mockTransactionRepo := new(mockRepo.MockTransactionRepository)
	uc := NewCoinTransactionUseCase(mockTransactionRepo, mockUserRepo)

	userFrom := &models.User{ID: "user1", Balance: 30}
	userTo := &models.User{ID: "user2", Balance: 50}

	mockUserRepo.On("FindUserByUsername", "user1").Return(userFrom, nil)
	mockUserRepo.On("FindUserByUsername", "user2").Return(userTo, nil)

	err := uc.SendCoins("user1", "user2", 50)

	assert.Error(t, err)
	assert.Equal(t, "недостаточно монет для отправки", err.Error())
	mockUserRepo.AssertExpectations(t)
}

func TestSendCoins_RecordTransactionError(t *testing.T) {
	mockUserRepo := new(mockRepo.MockUserRepository)
	mockTransactionRepo := new(mockRepo.MockTransactionRepository)
	uc := NewCoinTransactionUseCase(mockTransactionRepo, mockUserRepo)

	userFrom := &models.User{ID: "user1", Balance: 100}
	userTo := &models.User{ID: "user2", Balance: 50}

	mockUserRepo.On("FindUserByUsername", "user1").Return(userFrom, nil)
	mockUserRepo.On("FindUserByUsername", "user2").Return(userTo, nil)
	mockTransactionRepo.On("RecordTransaction", mock.Anything).Return(errors.New("database error"))

	err := uc.SendCoins("user1", "user2", 50)

	assert.Error(t, err)
	assert.Equal(t, "database error", err.Error())
	mockUserRepo.AssertExpectations(t)
	mockTransactionRepo.AssertExpectations(t)
}

func TestSendCoins_UpdateSenderBalanceError(t *testing.T) {
	mockUserRepo := new(mockRepo.MockUserRepository)
	mockTransactionRepo := new(mockRepo.MockTransactionRepository)
	uc := NewCoinTransactionUseCase(mockTransactionRepo, mockUserRepo)

	userFrom := &models.User{ID: "user1", Balance: 100}
	userTo := &models.User{ID: "user2", Balance: 50}

	mockUserRepo.On("FindUserByUsername", "user1").Return(userFrom, nil)
	mockUserRepo.On("FindUserByUsername", "user2").Return(userTo, nil)
	mockTransactionRepo.On("RecordTransaction", mock.Anything).Return(nil)
	mockUserRepo.On("UpdateUserBalance", "user1", -50).Return(errors.New("update balance error"))

	err := uc.SendCoins("user1", "user2", 50)

	assert.Error(t, err)
	assert.Equal(t, "update balance error", err.Error())
	mockUserRepo.AssertExpectations(t)
	mockTransactionRepo.AssertExpectations(t)
}

func TestSendCoins_UpdateReceiverBalanceError(t *testing.T) {
	mockUserRepo := new(mockRepo.MockUserRepository)
	mockTransactionRepo := new(mockRepo.MockTransactionRepository)
	uc := NewCoinTransactionUseCase(mockTransactionRepo, mockUserRepo)

	userFrom := &models.User{ID: "user1", Balance: 100}
	userTo := &models.User{ID: "user2", Balance: 50}

	mockUserRepo.On("FindUserByUsername", "user1").Return(userFrom, nil)
	mockUserRepo.On("FindUserByUsername", "user2").Return(userTo, nil)
	mockTransactionRepo.On("RecordTransaction", mock.Anything).Return(nil)
	mockUserRepo.On("UpdateUserBalance", "user1", -50).Return(nil)
	mockUserRepo.On("UpdateUserBalance", "user2", 50).Return(errors.New("update balance error"))

	err := uc.SendCoins("user1", "user2", 50)

	assert.Error(t, err)
	assert.Equal(t, "update balance error", err.Error())
	mockUserRepo.AssertExpectations(t)
	mockTransactionRepo.AssertExpectations(t)
}
