package repositories

import (
	"fmt"
	"github.com/rodrinoblega/stori/src/entities"
	"gorm.io/gorm"
)

type Database struct {
	*gorm.DB
}

func (d *Database) StoreTransactions(transactions entities.Transactions) error {
	result := d.Create(transactions)
	if result.Error != nil {
		return fmt.Errorf("failed to store transactions: %w", result.Error)
	}
	return nil
}

func (d *Database) GetAccountById(accountID int) (*entities.Account, error) {
	var account entities.Account
	if err := d.First(&account, accountID).Error; err != nil {
		return nil, err
	}
	return &account, nil
}
