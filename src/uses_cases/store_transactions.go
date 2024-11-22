package uses_cases

import (
	"github.com/rodrinoblega/stori/src/entities"
)

type StoreTransactionsUseCase struct {
	database Database
}

func NewStoreTransactionsUseCase(database Database) *StoreTransactionsUseCase {
	return &StoreTransactionsUseCase{
		database: database,
	}
}

func (s *StoreTransactionsUseCase) Execute(transactions entities.Transactions) error {
	err := s.database.StoreTransactions(transactions)
	if err != nil {
		return err
	}

	return nil
}
