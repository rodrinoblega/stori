package entities

import (
	"fmt"
	"time"
)

type Transactions []Transaction

type Transaction struct {
	TransactionID    string
	Date             time.Time
	Amount           float64
	TransactionsType string
	AccountID        string
}

func (t Transactions) GetAccountID() (string, error) {
	if len(t) == 0 {
		return "", fmt.Errorf("no transactions available to retrieve AccountID")
	}
	return t[0].AccountID, nil
}
