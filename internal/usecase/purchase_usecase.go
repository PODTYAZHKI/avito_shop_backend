package usecase

import (
	"errors"

	"avito-shop-test/internal/models"
)

type purchaseUseCase struct {
	purchaseRepo PurchaseRepository
	userRepo     UserRepository
	storeRepo    StoreRepository
}

func NewPurchaseUseCase(purchaseRepo PurchaseRepository, userRepo UserRepository, storeRepo StoreRepository) PurchaseUseCase {
	return &purchaseUseCase{
		purchaseRepo: purchaseRepo,
		userRepo:     userRepo,
		storeRepo:    storeRepo,
	}
}

func (uc *purchaseUseCase) BuyItem(username string, itemName string) error {
	user, err := uc.userRepo.FindUserByUsername(username)
	if err != nil || user == nil {
		return errors.New("пользователь не найден")
	}

	product, err := uc.storeRepo.GetItemByName(itemName)
	if err != nil {
		return errors.New("товар не найден")
	}

	if user.Balance < product.Price {
		return errors.New("недостаточно монет для покупки")
	}

	inventory := &models.Inventory{
		UserID:   user.ID,
		ItemType: product.Name,
		Quantity: 1,
	}

	if err := uc.userRepo.UpdateUserBalance(username, -product.Price); err != nil {
		return err
	}

	if err := uc.purchaseRepo.RecordPurchase(inventory); err != nil {
		return err
	}

	return nil
}
