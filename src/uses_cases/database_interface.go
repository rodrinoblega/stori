package uses_cases

import (
	"github.com/rodrinoblega/stori/src/entities"
)

type Database interface {
	StoreTransactions(transactions entities.Transactions) error
	GetAccountById(accountID int) (*entities.Account, error)
}
