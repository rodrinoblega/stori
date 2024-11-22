package repositories

import "github.com/rodrinoblega/stori/src/entities"

type DummyDB struct{}

func NewDummyDB() *DummyDB {
	return &DummyDB{}
}

func (n *DummyDB) StoreTransactions(_ entities.Transactions) error {
	return nil
}

func (d *DummyDB) GetAccountById(_ int) (*entities.Account, error) {
	return &entities.Account{}, nil
}
