package repository

import (
	"errors"

	"github.com/jonathantvrs/cvisa/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AccountLimitRepository struct {
	DB *gorm.DB
}

var (
	ErrAccountNotFound   = errors.New("Account not found!")
	ErrInsufficientFunds = errors.New("Transaction failed: insufficient funds!")
)

func (r *AccountLimitRepository) CreateAccountLimit(acc_limit *model.AccountLimit) error {
	return r.DB.Create(acc_limit).Error
}

func (r *AccountLimitRepository) GetLimitByAccountId(acc_id uint) (*model.AccountLimit, error) {
	var acc_limit model.AccountLimit
	if err := r.DB.First(&acc_limit, acc_id).Error; err != nil {
		return nil, err
	}
	return &acc_limit, nil
}

func (r *AccountLimitRepository) UpdateAccountLimit(acc_id uint, amount int64) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		var acc_limit model.AccountLimit
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&acc_limit, acc_id).Error; err != nil {
			return ErrAccountNotFound
		}

		newLimit := acc_limit.AvailableCreditLimit + amount

		if newLimit < 0 {
			return ErrInsufficientFunds
		}

		acc_limit.AvailableCreditLimit = newLimit
		return tx.Save(&acc_limit).Error
	})
}
