package usecase

import (
	"errors"
	"testing"

	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"avito-shop-test/internal/models"
	mockRepo "avito-shop-test/internal/repository/mock"
	mockToken "avito-shop-test/internal/token"
)

func TestAuthenticate_FindUserByUsername_Error(t *testing.T) {
	mockUserRepo := new(mockRepo.MockUserRepository)
	mockPurchaseRepo := new(mockRepo.MockPurchaseRepository)
	mockTransactionRepo := new(mockRepo.MockTransactionRepository)
	uc := NewUserUsecase(mockUserRepo, mockPurchaseRepo, mockTransactionRepo, mockToken.NewGenerator([]byte("secret")))

	mockUserRepo.On("FindUserByUsername", "testuser").Return(nil, errors.New("database error"))

	token, err := uc.Authenticate("testuser", "password")

	assert.Error(t, err)
	assert.Empty(t, token)
	assert.Equal(t, "database error", err.Error())
	mockUserRepo.AssertExpectations(t)
}

func TestAuthenticate_CreateUser_Error(t *testing.T) {
	mockUserRepo := new(mockRepo.MockUserRepository)
	mockPurchaseRepo := new(mockRepo.MockPurchaseRepository)
	mockTransactionRepo := new(mockRepo.MockTransactionRepository)
	uc := NewUserUsecase(mockUserRepo, mockPurchaseRepo, mockTransactionRepo, mockToken.NewGenerator([]byte("secret")))

	mockUserRepo.On("FindUserByUsername", "testuser").Return(nil, nil)
	mockUserRepo.On("CreateUser", mock.Anything).Return(errors.New("error creating user"))

	token, err := uc.Authenticate("testuser", "password")

	assert.Error(t, err)
	assert.Empty(t, token)
	assert.Equal(t, "error creating user", err.Error())
	mockUserRepo.AssertExpectations(t)
}

func TestAuthenticate_UserCreation_ValIDationError(t *testing.T) {
	mockUserRepo := new(mockRepo.MockUserRepository)
	mockPurchaseRepo := new(mockRepo.MockPurchaseRepository)
	mockTransactionRepo := new(mockRepo.MockTransactionRepository)
	uc := NewUserUsecase(mockUserRepo, mockPurchaseRepo, mockTransactionRepo, mockToken.NewGenerator([]byte("secret")))

	mockUserRepo.On("FindUserByUsername", "testuser").Return(nil, nil)
	mockUserRepo.On("CreateUser", mock.Anything).Return(errors.New("valIDation error"))

	token, err := uc.Authenticate("testuser", "password")

	assert.Error(t, err)
	assert.Empty(t, token)
	assert.Equal(t, "valIDation error", err.Error())

	mockUserRepo.AssertExpectations(t)
}

func TestAuthenticate_EmptyUsernameOrPassword(t *testing.T) {
	mockUserRepo := new(mockRepo.MockUserRepository)
	mockPurchaseRepo := new(mockRepo.MockPurchaseRepository)
	mockTransactionRepo := new(mockRepo.MockTransactionRepository)
	uc := NewUserUsecase(mockUserRepo, mockPurchaseRepo, mockTransactionRepo, mockToken.NewGenerator([]byte("secret")))

	token, err := uc.Authenticate("", "")

	assert.Error(t, err)
	assert.Empty(t, token)
}

func TestAuthenticate_UserNotFound_CreateUser(t *testing.T) {
	mockUserRepo := new(mockRepo.MockUserRepository)
	mockPurchaseRepo := new(mockRepo.MockPurchaseRepository)
	mockTransactionRepo := new(mockRepo.MockTransactionRepository)
	uc := NewUserUsecase(mockUserRepo, mockPurchaseRepo, mockTransactionRepo, mockToken.NewGenerator([]byte("secret")))

	mockUserRepo.On("FindUserByUsername", "testuser").Return(nil, nil)
	mockUserRepo.On("CreateUser", mock.Anything).Return(nil)

	token, err := uc.Authenticate("testuser", "password")

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	mockUserRepo.AssertExpectations(t)
}

func TestAuthenticate_UserFound_WrongPassword(t *testing.T) {
	mockUserRepo := new(mockRepo.MockUserRepository)
	mockPurchaseRepo := new(mockRepo.MockPurchaseRepository)
	mockTransactionRepo := new(mockRepo.MockTransactionRepository)
	uc := NewUserUsecase(mockUserRepo, mockPurchaseRepo, mockTransactionRepo, mockToken.NewGenerator([]byte("secret")))

	user := &models.User{Username: "testuser", Password: "password"}
	mockUserRepo.On("FindUserByUsername", "testuser").Return(user, nil)

	token, err := uc.Authenticate("testuser", "wrongpassword")

	assert.Error(t, err)
	assert.Equal(t, "неавторизован", err.Error())
	assert.Empty(t, token)
	mockUserRepo.AssertExpectations(t)
}

func TestAuthenticate_UserFound_CorrectPassword(t *testing.T) {
	mockUserRepo := new(mockRepo.MockUserRepository)
	mockPurchaseRepo := new(mockRepo.MockPurchaseRepository)
	mockTransactionRepo := new(mockRepo.MockTransactionRepository)
	uc := NewUserUsecase(mockUserRepo, mockPurchaseRepo, mockTransactionRepo, mockToken.NewGenerator([]byte("secret")))

	user := &models.User{Username: "testuser", Password: "password"}
	mockUserRepo.On("FindUserByUsername", "testuser").Return(user, nil)

	token, err := uc.Authenticate("testuser", "password")

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	mockUserRepo.AssertExpectations(t)
}

func TestAuthenticate_GenerateToken_Success(t *testing.T) {
	mockUserRepo := new(mockRepo.MockUserRepository)
	mockPurchaseRepo := new(mockRepo.MockPurchaseRepository)
	mockTransactionRepo := new(mockRepo.MockTransactionRepository)
	uc := NewUserUsecase(mockUserRepo, mockPurchaseRepo, mockTransactionRepo, mockToken.NewGenerator([]byte("secret")))

	user := &models.User{Username: "testuser", Password: "password"}
	mockUserRepo.On("FindUserByUsername", "testuser").Return(user, nil)

	token, err := uc.Authenticate("testuser", "password")

	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	claims := jwt.MapClaims{}
	parsedToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})

	assert.NoError(t, err)
	assert.True(t, parsedToken.Valid)
	assert.Equal(t, user.Username, claims["username"])
}

func TestAuthenticate_TokenGenerationError(t *testing.T) {
	mockUserRepo := new(mockRepo.MockUserRepository)
	mockPurchaseRepo := new(mockRepo.MockPurchaseRepository)
	mockTransactionRepo := new(mockRepo.MockTransactionRepository)
	uc := NewUserUsecase(mockUserRepo, mockPurchaseRepo, mockTransactionRepo, mockToken.NewGenerator([]byte("secret"))).(*userUseCase)
	mockTokenGenerator := new(mockToken.MockTokenGenerator)

	user := &models.User{Username: "testuser", Password: "password"}
	mockUserRepo.On("FindUserByUsername", "testuser").Return(user, nil)

	mockTokenGenerator.On("Generate", user.Username).Return("", errors.New("token generation error"))

	uc.tokenGenerator = mockTokenGenerator

	token, err := uc.Authenticate("testuser", "password")

	assert.Error(t, err)
	assert.Empty(t, token)
	assert.Equal(t, "token generation error", err.Error())

	mockUserRepo.AssertExpectations(t)
	mockTokenGenerator.AssertExpectations(t)
}

func TestGetUserInfo_GetPurchasedItems_Error(t *testing.T) {
	mockUserRepo := new(mockRepo.MockUserRepository)
	mockPurchaseRepo := new(mockRepo.MockPurchaseRepository)
	mockTransactionRepo := new(mockRepo.MockTransactionRepository)
	uc := NewUserUsecase(mockUserRepo, mockPurchaseRepo, mockTransactionRepo, mockToken.NewGenerator([]byte("secret")))

	user := &models.User{ID: "user-ID-1"}

	mockUserRepo.On("FindUserByUsername", "testuser").Return(user, nil)
	mockPurchaseRepo.On("GetPurchasedItems", user.ID).Return(nil, errors.New("error retrieving items"))

	userInfo, err := uc.GetUserInfo("testuser")

	assert.Error(t, err)
	assert.Nil(t, userInfo)
	assert.Equal(t, "error retrieving items", err.Error())
	mockUserRepo.AssertExpectations(t)
	mockPurchaseRepo.AssertExpectations(t)
}

func TestGetUserInfo_Success(t *testing.T) {
	mockUserRepo := new(mockRepo.MockUserRepository)
	mockPurchaseRepo := new(mockRepo.MockPurchaseRepository)
	mockTransactionRepo := new(mockRepo.MockTransactionRepository)
	uc := NewUserUsecase(mockUserRepo, mockPurchaseRepo, mockTransactionRepo, mockToken.NewGenerator([]byte("secret")))

	user := &models.User{Username: "testuser", Balance: 100, ID: "user-ID-1"}
	mockUserRepo.On("FindUserByUsername", "testuser").Return(user, nil)

	mockPurchaseRepo.On("GetPurchasedItems", user.ID).Return([]models.Inventory{
		{UserID: user.ID, ItemType: "item1", Quantity: 1},
	}, nil)

	mockTransactionRepo.On("GetTransactionsHistory", user.ID).Return([]models.CoinTransaction{}, nil)

	userInfo, err := uc.GetUserInfo("testuser")

	assert.NoError(t, err)
	assert.NotNil(t, userInfo)
	assert.Equal(t, user.Balance, userInfo.Coins)
	assert.Len(t, userInfo.Inventory, 1)
	mockUserRepo.AssertExpectations(t)
	mockPurchaseRepo.AssertExpectations(t)
	mockTransactionRepo.AssertExpectations(t)
}

func TestGetUserInfo_UserNotFound(t *testing.T) {
	mockUserRepo := new(mockRepo.MockUserRepository)
	mockPurchaseRepo := new(mockRepo.MockPurchaseRepository)
	mockTransactionRepo := new(mockRepo.MockTransactionRepository)
	uc := NewUserUsecase(mockUserRepo, mockPurchaseRepo, mockTransactionRepo, mockToken.NewGenerator([]byte("secret")))

	mockUserRepo.On("FindUserByUsername", "testuser").Return(nil, errors.New("user not found"))

	userInfo, err := uc.GetUserInfo("testuser")

	assert.Error(t, err)
	assert.Empty(t, userInfo)
	assert.Equal(t, "user not found", err.Error())
	mockUserRepo.AssertExpectations(t)
}

func TestGetUserInfo_GetCoinHistory_Error(t *testing.T) {
	mockUserRepo := new(mockRepo.MockUserRepository)
	mockPurchaseRepo := new(mockRepo.MockPurchaseRepository)
	mockTransactionRepo := new(mockRepo.MockTransactionRepository)
	uc := NewUserUsecase(mockUserRepo, mockPurchaseRepo, mockTransactionRepo, mockToken.NewGenerator([]byte("secret")))

	user := &models.User{ID: "user-ID-1"}

	mockUserRepo.On("FindUserByUsername", "testuser").Return(user, nil)
	mockPurchaseRepo.On("GetPurchasedItems", user.ID).Return([]models.Inventory{}, nil)
	mockTransactionRepo.On("GetTransactionsHistory", user.ID).Return(nil, errors.New("error retrieving coin history"))

	userInfo, err := uc.GetUserInfo("testuser")

	assert.Error(t, err)
	assert.Nil(t, userInfo)
	assert.Equal(t, "error retrieving coin history", err.Error())
	mockUserRepo.AssertExpectations(t)
	mockPurchaseRepo.AssertExpectations(t)
	mockTransactionRepo.AssertExpectations(t)
}

func TestGetCoinHistory_Error(t *testing.T) {
	mockUserRepo := new(mockRepo.MockUserRepository)
	mockPurchaseRepo := new(mockRepo.MockPurchaseRepository)
	mockTransactionRepo := new(mockRepo.MockTransactionRepository)
	uc := NewUserUsecase(mockUserRepo, mockPurchaseRepo, mockTransactionRepo, mockToken.NewGenerator([]byte("secret")))

	user := &models.User{ID: "user-ID-1"}
	mockTransactionRepo.On("GetTransactionsHistory", user.ID).Return(nil, errors.New("error retrieving transactions"))

	coinHistory, err := uc.GetCoinHistory(user.ID)

	assert.Error(t, err)
	assert.Equal(t, "error retrieving transactions", err.Error())
	assert.Empty(t, coinHistory.Received)
	assert.Empty(t, coinHistory.Sent)
}

func TestGetCoinHistory_Success(t *testing.T) {
	mockUserRepo := new(mockRepo.MockUserRepository)
	mockPurchaseRepo := new(mockRepo.MockPurchaseRepository)
	mockTransactionRepo := new(mockRepo.MockTransactionRepository)
	uc := NewUserUsecase(mockUserRepo, mockPurchaseRepo, mockTransactionRepo, mockToken.NewGenerator([]byte("secret")))

	user := &models.User{ID: "user-ID-1"}

	mockTransactionRepo.On("GetTransactionsHistory", user.ID).Return([]models.CoinTransaction{
		{FromUser: "user-ID-1", ToUser: "user-ID-2", Amount: 10},
		{FromUser: "user-ID-2", ToUser: "user-ID-1", Amount: 5},
	}, nil)

	coinHistory, err := uc.GetCoinHistory(user.ID)

	assert.NoError(t, err)
	assert.Len(t, coinHistory.Received, 1)
	assert.Len(t, coinHistory.Sent, 1)
	assert.Equal(t, 10, coinHistory.Sent[0].Amount)
	assert.Equal(t, "user-ID-2", coinHistory.Sent[0].Username)
	assert.Equal(t, 5, coinHistory.Received[0].Amount)
	assert.Equal(t, "user-ID-2", coinHistory.Received[0].Username)
	mockUserRepo.AssertExpectations(t)
	mockTransactionRepo.AssertExpectations(t)
}

func TestGetCoinHistory_NoTransactions(t *testing.T) {
	mockUserRepo := new(mockRepo.MockUserRepository)
	mockPurchaseRepo := new(mockRepo.MockPurchaseRepository)
	mockTransactionRepo := new(mockRepo.MockTransactionRepository)
	uc := NewUserUsecase(mockUserRepo, mockPurchaseRepo, mockTransactionRepo, mockToken.NewGenerator([]byte("secret")))

	user := &models.User{ID: "user-ID-1"}

	mockTransactionRepo.On("GetTransactionsHistory", user.ID).Return([]models.CoinTransaction{}, nil)

	coinHistory, err := uc.GetCoinHistory(user.ID)

	assert.NoError(t, err)
	assert.Len(t, coinHistory.Received, 0)
	assert.Len(t, coinHistory.Sent, 0)
	mockUserRepo.AssertExpectations(t)
	mockTransactionRepo.AssertExpectations(t)
}

func TestGetCoinHistory_NoTransactions_EmptyHistory(t *testing.T) {
	mockUserRepo := new(mockRepo.MockUserRepository)
	mockPurchaseRepo := new(mockRepo.MockPurchaseRepository)
	mockTransactionRepo := new(mockRepo.MockTransactionRepository)
	uc := NewUserUsecase(mockUserRepo, mockPurchaseRepo, mockTransactionRepo, mockToken.NewGenerator([]byte("secret")))

	user := &models.User{ID: "user-ID-1"}

	mockTransactionRepo.On("GetTransactionsHistory", user.ID).Return([]models.CoinTransaction{}, nil)

	coinHistory, err := uc.GetCoinHistory(user.ID)

	assert.NoError(t, err)
	assert.Len(t, coinHistory.Received, 0)
	assert.Len(t, coinHistory.Sent, 0)

	mockUserRepo.AssertExpectations(t)
	mockTransactionRepo.AssertExpectations(t)
}

func TestGetPurchasedItems_Success(t *testing.T) {
	mockUserRepo := new(mockRepo.MockUserRepository)
	mockPurchaseRepo := new(mockRepo.MockPurchaseRepository)
	mockTransactionRepo := new(mockRepo.MockTransactionRepository)
	uc := NewUserUsecase(mockUserRepo, mockPurchaseRepo, mockTransactionRepo, mockToken.NewGenerator([]byte("secret")))

	user := &models.User{ID: "user-ID-1"}

	mockPurchaseRepo.On("GetPurchasedItems", user.ID).Return([]models.Inventory{
		{UserID: user.ID, ItemType: "item1", Quantity: 2},
		{UserID: user.ID, ItemType: "item2", Quantity: 3},
	}, nil)

	items, err := uc.GetPurchasedItems(user.ID)

	assert.NoError(t, err)

	expectedItems := []models.PurchasedItem{
		{ItemName: "item1", Quantity: 2},
		{ItemName: "item2", Quantity: 3},
	}
	assert.ElementsMatch(t, expectedItems, items)

	mockUserRepo.AssertExpectations(t)
	mockPurchaseRepo.AssertExpectations(t)
}

func TestGetPurchasedItems_NoItems(t *testing.T) {
	mockUserRepo := new(mockRepo.MockUserRepository)
	mockPurchaseRepo := new(mockRepo.MockPurchaseRepository)
	mockTransactionRepo := new(mockRepo.MockTransactionRepository)
	uc := NewUserUsecase(mockUserRepo, mockPurchaseRepo, mockTransactionRepo, mockToken.NewGenerator([]byte("secret")))

	user := &models.User{ID: "user-ID-1"}

	mockPurchaseRepo.On("GetPurchasedItems", user.ID).Return([]models.Inventory{}, nil)

	items, err := uc.GetPurchasedItems(user.ID)

	assert.NoError(t, err)
	assert.Empty(t, items)
	mockUserRepo.AssertExpectations(t)
	mockPurchaseRepo.AssertExpectations(t)
}

func TestGetPurchasedItems_Error(t *testing.T) {
	mockUserRepo := new(mockRepo.MockUserRepository)
	mockPurchaseRepo := new(mockRepo.MockPurchaseRepository)
	mockTransactionRepo := new(mockRepo.MockTransactionRepository)
	uc := NewUserUsecase(mockUserRepo, mockPurchaseRepo, mockTransactionRepo, mockToken.NewGenerator([]byte("secret")))

	user := &models.User{ID: "user-ID-1"}

	mockPurchaseRepo.On("GetPurchasedItems", user.ID).Return(nil, errors.New("error retrieving items"))

	items, err := uc.GetPurchasedItems(user.ID)

	assert.Error(t, err)
	assert.Nil(t, items)
	assert.Equal(t, "error retrieving items", err.Error())
	mockUserRepo.AssertExpectations(t)
	mockPurchaseRepo.AssertExpectations(t)
}

func TestGetPurchasedItems_InvalIDItems(t *testing.T) {
	mockUserRepo := new(mockRepo.MockUserRepository)
	mockPurchaseRepo := new(mockRepo.MockPurchaseRepository)
	mockTransactionRepo := new(mockRepo.MockTransactionRepository)
	uc := NewUserUsecase(mockUserRepo, mockPurchaseRepo, mockTransactionRepo, mockToken.NewGenerator([]byte("secret")))

	user := &models.User{ID: "user-ID-1"}

	mockPurchaseRepo.On("GetPurchasedItems", user.ID).Return([]models.Inventory{
		{UserID: user.ID, ItemType: "", Quantity: 0},
	}, nil)

	items, err := uc.GetPurchasedItems(user.ID)

	assert.NoError(t, err)
	assert.Len(t, items, 1)
	assert.Empty(t, items[0].ItemName)
	assert.Equal(t, 0, items[0].Quantity)
	mockUserRepo.AssertExpectations(t)
	mockPurchaseRepo.AssertExpectations(t)
}

func TestGetPurchasedItems_DuplicateItems(t *testing.T) {
	mockUserRepo := new(mockRepo.MockUserRepository)
	mockPurchaseRepo := new(mockRepo.MockPurchaseRepository)
	mockTransactionRepo := new(mockRepo.MockTransactionRepository)
	uc := NewUserUsecase(mockUserRepo, mockPurchaseRepo, mockTransactionRepo, mockToken.NewGenerator([]byte("secret")))

	user := &models.User{ID: "user-ID-1"}

	mockPurchaseRepo.On("GetPurchasedItems", user.ID).Return([]models.Inventory{
		{UserID: user.ID, ItemType: "item1", Quantity: 2},
		{UserID: user.ID, ItemType: "item1", Quantity: 3},
	}, nil)

	items, err := uc.GetPurchasedItems(user.ID)

	assert.NoError(t, err)
	assert.Len(t, items, 1)
	assert.Equal(t, "item1", items[0].ItemName)
	assert.Equal(t, 5, items[0].Quantity)
	mockUserRepo.AssertExpectations(t)
	mockPurchaseRepo.AssertExpectations(t)
}

func TestGetPurchasedItems_OneItem(t *testing.T) {
	mockUserRepo := new(mockRepo.MockUserRepository)
	mockPurchaseRepo := new(mockRepo.MockPurchaseRepository)
	mockTransactionRepo := new(mockRepo.MockTransactionRepository)
	uc := NewUserUsecase(mockUserRepo, mockPurchaseRepo, mockTransactionRepo, mockToken.NewGenerator([]byte("secret")))

	user := &models.User{ID: "user-ID-1"}

	mockPurchaseRepo.On("GetPurchasedItems", user.ID).Return([]models.Inventory{
		{UserID: user.ID, ItemType: "item1", Quantity: 1},
	}, nil)

	items, err := uc.GetPurchasedItems(user.ID)

	assert.NoError(t, err)
	assert.Len(t, items, 1)
	assert.Equal(t, "item1", items[0].ItemName)
	assert.Equal(t, 1, items[0].Quantity)

	mockUserRepo.AssertExpectations(t)
	mockPurchaseRepo.AssertExpectations(t)
}
