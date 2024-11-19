package entities

import "time"

type Transactions []Transaction

type Transaction struct {
	TransactionID    string
	Date             time.Time
	Amount           float64
	TransactionsType string
	AccountID        string
}
