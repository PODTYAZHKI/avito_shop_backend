package usecase

import (
	"errors"

	"avito-shop-test/internal/models"
)

type coinTransactionUseCase struct {
	coinTransactionRepo CoinTransactionRepository
	userRepo            UserRepository
}

func NewCoinTransactionUseCase(coinTransactionRepo CoinTransactionRepository, userRepo UserRepository) CoinTransactionUseCase {
	return &coinTransactionUseCase{
		coinTransactionRepo: coinTransactionRepo,
		userRepo:            userRepo,
	}
}

func (uc *coinTransactionUseCase) SendCoins(fromUser, toUser string, amount int) error {

	userFrom, err := uc.userRepo.FindUserByUsername(fromUser)
	if err != nil || userFrom == nil {
		return errors.New("отправитель не найден")
	}

	userTo, err := uc.userRepo.FindUserByUsername(toUser)
	if err != nil || userTo == nil {
		return errors.New("получатель не найден")
	}

	if userFrom.Balance < amount {
		return errors.New("недостаточно монет для отправки")
	}

	transaction := &models.CoinTransaction{
		FromUser: userFrom.ID,
		ToUser:   userTo.ID,
		Amount:   amount,
	}

	err = uc.coinTransactionRepo.RecordTransaction(transaction)
	if err != nil {
		return err
	}

	if err := uc.userRepo.UpdateUserBalance(fromUser, -amount); err != nil {
		return err
	}
	if err := uc.userRepo.UpdateUserBalance(toUser, amount); err != nil {
		return err
	}

	return nil
}
