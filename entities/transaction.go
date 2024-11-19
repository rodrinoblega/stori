package entities

import "time"

type Transactions []Transaction

type Transaction struct {
	ID               string
	Date             time.Time
	Amount           float64
	TransactionsType string
}
