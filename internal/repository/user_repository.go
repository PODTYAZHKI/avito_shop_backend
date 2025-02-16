package repository

import (
	"gorm.io/gorm"

	"github.com/pkg/errors"

	"avito-shop-test/internal/models"
)

type UserRepository interface {
	FindUserByUsername(username string) (*models.User, error)
	CreateUser(user *models.User) error
	UpdateUserBalance(username string, amount int) error
	GetUserByUserID(userID string) (*models.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUSerRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (userDb *userRepository) FindUserByUsername(username string) (*models.User, error) {
	user := models.User{}
	tx := userDb.db.Table("users").Where("username = ?", username).Take(&user)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, errors.Wrap(tx.Error, "database error (table users)")
	}
	return &user, nil
}

func (userDb *userRepository) GetUserByUserID(userID string) (*models.User, error) {
	user := models.User{}
	tx := userDb.db.Table("users").Where("ID =?", userID).Take(&user)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, errors.Wrap(tx.Error, "database error (table users)")
	}
	return &user, nil
}

func (userDb *userRepository) CreateUser(user *models.User) error {
	tx := userDb.db.Create(user)
	if tx.Error != nil {
		return errors.Wrap(tx.Error, "database error (table users)")
	}
	return nil
}

func (r *userRepository) UpdateUserBalance(username string, amount int) error {
	return r.db.Model(&models.User{}).Where("username = ?", username).Update("Balance", gorm.Expr("Balance + ?", amount)).Error
}
