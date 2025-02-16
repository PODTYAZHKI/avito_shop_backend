package usecase

import (
	"errors"

	"avito-shop-test/internal/models"
)

type userUseCase struct {
	userRepo            UserRepository
	purchaseRepo        PurchaseRepository
	coinTransactionRepo CoinTransactionRepository
	tokenGenerator      TokenGenerator
}

func NewUserUsecase(userRepo UserRepository, purchaseRepo PurchaseRepository, coinTransactionRepo CoinTransactionRepository, tokenGenerator TokenGenerator) UserUseCase {
	return &userUseCase{
		userRepo:            userRepo,
		purchaseRepo:        purchaseRepo,
		coinTransactionRepo: coinTransactionRepo,
		tokenGenerator:      tokenGenerator,
	}
}

func (uc *userUseCase) Authenticate(username, password string) (string, error) {
	if username == "" || password == "" {
		return "", errors.New("username and password cannot be empty")
	}
	user, err := uc.userRepo.FindUserByUsername(username)
	if err != nil {
		return "", err
	}

	if user == nil {

		newUser := &models.User{Username: username, Password: password, Balance: 1000}
		if err := uc.userRepo.CreateUser(newUser); err != nil {
			return "", err
		}
		user = newUser
	} else {

		if user.Password != password {
			return "", errors.New("неавторизован")
		}
	}

	tokenString, err := uc.tokenGenerator.Generate(user.Username)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (uc *userUseCase) GetCoinHistory(userID string) (models.CoinHistory, error) {
	transactions, err := uc.coinTransactionRepo.GetTransactionsHistory(userID)

	if err != nil {
		return models.CoinHistory{}, err
	}
	history := models.CoinHistory{}

	for _, transaction := range transactions {
		if transaction.FromUser == userID {
			user, _ := uc.userRepo.GetUserByUserID(transaction.ToUser)

			history.Sent = append(history.Sent, models.CoinTransactionInfo{
				Amount:   transaction.Amount,
				Username: user.Username,
			})
			continue
		}
		user, _ := uc.userRepo.GetUserByUserID(transaction.FromUser)
		history.Received = append(history.Received, models.CoinTransactionInfo{
			Amount:   transaction.Amount,
			Username: user.Username,
		})

	}

	return history, nil
}

func (uc *userUseCase) GetPurchasedItems(userID string) ([]models.PurchasedItem, error) {
	purchases, err := uc.purchaseRepo.GetPurchasedItems(userID)
	if err != nil {
		return nil, err
	}

	itemMap := make(map[string]*models.PurchasedItem)

	for _, purchase := range purchases {
		if item, exists := itemMap[purchase.ItemType]; exists {

			item.Quantity += purchase.Quantity
		} else {

			itemMap[purchase.ItemType] = &models.PurchasedItem{
				ItemName: purchase.ItemType,
				Quantity: purchase.Quantity,
			}
		}
	}

	var items []models.PurchasedItem
	for _, item := range itemMap {
		items = append(items, *item)
	}

	return items, nil
}

func (uc *userUseCase) GetUserInfo(username string) (*models.UserInfo, error) {
	user, err := uc.userRepo.FindUserByUsername(username)
	if err != nil {
		return nil, err
	}

	inventory, err := uc.GetPurchasedItems(user.ID)
	if err != nil {
		return nil, err
	}

	coinHistory, err := uc.GetCoinHistory(user.ID)
	if err != nil {
		return nil, err
	}

	return &models.UserInfo{
		Coins:       user.Balance,
		Inventory:   inventory,
		CoinHistory: coinHistory,
	}, nil
}
