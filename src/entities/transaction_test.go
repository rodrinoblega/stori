package entities

import (
	"fmt"
	"testing"
	"time"
)

func TestLength(t *testing.T) {
	tests := []struct {
		name         string
		transactions Transactions
		expected     int
	}{
		{
			name:         "Non-empty Transactions",
			transactions: Transactions{{TransactionID: "1", Amount: 100, TransactionsType: "CREDIT"}},
			expected:     1,
		},
		{
			name:         "Empty Transactions",
			transactions: Transactions{},
			expected:     0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.transactions.Length()
			if actual != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, actual)
			}
		})
	}
}

func TestGetAccountID(t *testing.T) {
	tests := []struct {
		name         string
		transactions Transactions
		expectedID   string
		expectedErr  error
	}{
		{
			name:         "Get Account ID from Non-empty Transactions",
			transactions: Transactions{{AccountID: "123"}},
			expectedID:   "123",
			expectedErr:  nil,
		},
		{
			name:         "Get Account ID from Empty Transactions",
			transactions: Transactions{},
			expectedID:   "",
			expectedErr:  fmt.Errorf("no transactions available to retrieve AccountID"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualID, actualErr := tt.transactions.GetAccountID()
			if actualID != tt.expectedID || (actualErr != nil && actualErr.Error() != tt.expectedErr.Error()) {
				t.Errorf("expected %v and %v, got %v and %v", tt.expectedID, tt.expectedErr, actualID, actualErr)
			}
		})
	}
}

func TestTotalBalance(t *testing.T) {
	tests := []struct {
		name         string
		transactions Transactions
		expected     float64
	}{
		{
			name: "Total Balance of Transactions",
			transactions: Transactions{
				{Amount: 100, TransactionsType: "CREDIT"},
				{Amount: 50, TransactionsType: "DEBIT"},
			},
			expected: 150.00,
		},
		{
			name:         "Empty Transactions",
			transactions: Transactions{},
			expected:     0.00,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.transactions.TotalBalance()
			if actual != tt.expected {
				t.Errorf("expected %.2f, got %.2f", tt.expected, actual)
			}
		})
	}
}

func TestTransactionsByMonthYear(t *testing.T) {
	tests := []struct {
		name         string
		transactions Transactions
		expected     []MonthYearTransactions
	}{
		{
			name: "Transactions grouped by Month/Year",
			transactions: Transactions{
				{Date: time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC)},
				{Date: time.Date(2024, time.January, 2, 0, 0, 0, 0, time.UTC)},
				{Date: time.Date(2024, time.February, 1, 0, 0, 0, 0, time.UTC)},
			},
			expected: []MonthYearTransactions{
				{MonthYear: "January 2024", TransactionCount: 2},
				{MonthYear: "February 2024", TransactionCount: 1},
			},
		},
		{
			name:         "No Transactions",
			transactions: Transactions{},
			expected:     []MonthYearTransactions{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.transactions.TransactionsByMonthYear()
			if len(actual) != len(tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, actual)
			}
		})
	}
}

func TestAverageTransactionsAmount(t *testing.T) {
	tests := []struct {
		name         string
		transactions Transactions
		expected     AverageSummary
	}{
		{
			name: "Average Credit and Debit",
			transactions: Transactions{
				{Amount: 100, TransactionsType: "CREDIT"},
				{Amount: 50, TransactionsType: "DEBIT"},
				{Amount: 200, TransactionsType: "CREDIT"},
				{Amount: 30, TransactionsType: "DEBIT"},
			},
			expected: AverageSummary{
				AverageCreditAmount: 150.00,
				AverageDebitAmount:  40.00,
			},
		},
		{
			name:         "No Transactions",
			transactions: Transactions{},
			expected:     AverageSummary{AverageCreditAmount: 0, AverageDebitAmount: 0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.transactions.AverageTransactionsAmount()
			if actual.AverageCreditAmount != tt.expected.AverageCreditAmount || actual.AverageDebitAmount != tt.expected.AverageDebitAmount {
				t.Errorf("expected %v, got %v", tt.expected, actual)
			}
		})
	}
}
