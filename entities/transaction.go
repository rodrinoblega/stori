package entities

import (
	"fmt"
	"math"
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

type MonthYearTransactions struct {
	MonthYear        string
	TransactionCount int
}

type AverageSummary struct {
	AverageCreditAmount float64
	AverageDebitAmount  float64
}

func (t Transactions) GetAccountID() (string, error) {
	if len(t) == 0 {
		return "", fmt.Errorf("no transactions available to retrieve AccountID")
	}
	return t[0].AccountID, nil
}

func (t Transactions) TotalBalance() float64 {
	var total float64
	for _, transaction := range t {
		total += transaction.Amount
	}

	return math.Round(total*100) / 100
}

func (t Transactions) TransactionsByMonthYear() []MonthYearTransactions {
	monthlyCounts := make(map[string]int)

	for _, transaction := range t {
		monthYear := transaction.Date.Format("2006-01")
		monthlyCounts[monthYear]++
	}

	var results []MonthYearTransactions
	for monthYear, count := range monthlyCounts {
		parsedMonthYear, _ := time.Parse("2006-01", monthYear)
		formattedMonthYear := parsedMonthYear.Format("January 2006")

		results = append(results, MonthYearTransactions{
			MonthYear:        formattedMonthYear,
			TransactionCount: count,
		})
	}

	return results
}

func (t Transactions) AverageTransactionsAmount() AverageSummary {
	var totalCreditAmount, totalDebitAmount float64
	var totalCreditCount, totalDebitCount int

	for _, transaction := range t {
		if transaction.TransactionsType == "CREDIT" {
			totalCreditAmount += transaction.Amount
			totalCreditCount++
		} else if transaction.TransactionsType == "DEBIT" {
			totalDebitAmount += transaction.Amount
			totalDebitCount++
		}
	}

	averageCreditAmount := calculateAverage(totalCreditAmount, totalCreditCount)
	averageDebitAmount := calculateAverage(totalDebitAmount, totalDebitCount)

	return AverageSummary{
		AverageCreditAmount: averageCreditAmount,
		AverageDebitAmount:  averageDebitAmount,
	}
}

func calculateAverage(total float64, count int) float64 {
	if count == 0 {
		return 0
	}
	return total / float64(count)
}
