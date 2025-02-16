package usecase

import "avito-shop-test/internal/models"

type CoinTransactionRepository interface {
	RecordTransaction(transaction *models.CoinTransaction) error
	GetTransactionsHistory(userID string) ([]models.CoinTransaction, error)
}

type CoinTransactionUseCase interface {
	SendCoins(fromUser, toUser string, amount int) error
}

type UserRepository interface {
	FindUserByUsername(username string) (*models.User, error)
	CreateUser(user *models.User) error
	UpdateUserBalance(username string, amount int) error
	GetUserByUserID(userID string) (*models.User, error)
}

type UserUseCase interface {
	Authenticate(username, password string) (string, error)
	GetUserInfo(username string) (*models.UserInfo, error)
	GetCoinHistory(userID string) (models.CoinHistory, error)
	GetPurchasedItems(userID string) ([]models.PurchasedItem, error)
}

type PurchaseRepository interface {
	RecordPurchase(inventory *models.Inventory) error
	GetPurchasedItems(userID string) ([]models.Inventory, error)
}

type PurchaseUseCase interface {
	BuyItem(username string, itemName string) error
}

type StoreRepository interface {
	GetItemByName(name string) (*models.Product, error)
}

type StoreUseCase interface {
	LoadItems(filename string) error
}

type TokenGenerator interface {
	Generate(username string) (string, error)
}
