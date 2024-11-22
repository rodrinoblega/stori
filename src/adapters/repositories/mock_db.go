package repositories

import "github.com/rodrinoblega/stori/src/entities"

type MockTestDB struct{}

func NewMockTestDB() *MockTestDB {
	return &MockTestDB{}
}

func (n *MockTestDB) StoreTransactions(_ entities.Transactions) error {
	return nil
}

func (d *MockTestDB) GetAccountById(accountID int) (*entities.Account, error) {
	return &entities.Account{}, nil
}
