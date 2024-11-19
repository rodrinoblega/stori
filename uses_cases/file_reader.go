package uses_cases

import (
	"encoding/csv"
	"fmt"
	"github.com/rodrinoblega/stori/entities"
	"os"
	"strconv"
	"time"
)

type FileReaderUseCase struct{}

func NewFileReaderUseCase() *FileReaderUseCase {
	return &FileReaderUseCase{}
}

func (f *FileReaderUseCase) Execute(filePath string) (entities.Transactions, error) {
	file, err := openFile(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	rows, err := readCSV(file)
	if err != nil {
		return nil, err
	}

	transactions, err := processRows(rows)
	if err != nil {
		return nil, err
	}

	return transactions, nil
}

func openFile(filePath string) (*os.File, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %v", filePath, err)
	}
	return file, nil
}

func readCSV(file *os.File) ([][]string, error) {
	reader := csv.NewReader(file)
	rows, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV: %v", err)
	}
	return rows, nil
}

func processRows(rows [][]string) (entities.Transactions, error) {
	var transactions entities.Transactions

	for i, row := range rows {
		if i == 0 { // Skip header row
			continue
		}

		if len(row) != 3 {
			return nil, fmt.Errorf("invalid row format at line %d: %v", i+1, row)
		}

		transaction, err := processLine(row)
		if err != nil {
			return nil, fmt.Errorf("error processing row %d: %v", i+1, err)
		}

		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

func processLine(row []string) (entities.Transaction, error) {
	amount, err := strconv.ParseFloat(row[2], 64)
	if err != nil {
		return entities.Transaction{}, fmt.Errorf("invalid amount '%s': %v", row[2], err)
	}

	transationType, err := resolveTransactionType(amount)
	if err != nil {
		return entities.Transaction{}, err
	}

	date, err := parseDate(row[1], err)
	if err != nil {
		return entities.Transaction{}, err
	}

	return entities.Transaction{
		ID:               row[0],
		Date:             date,
		Amount:           amount,
		TransactionsType: transationType,
	}, nil
}

func parseDate(rowInfo string, err error) (time.Time, error) {
	layout := "01/02/2006"

	date, err := time.Parse(layout, rowInfo)
	if err != nil {
		return time.Time{}, fmt.Errorf("error parsing date: %v", err)
	}
	return date, nil
}

func resolveTransactionType(amount float64) (string, error) {
	switch {
	case amount > 0:
		return "CREDIT", nil
	case amount < 0:
		return "DEBIT", nil
	default:
		return "", fmt.Errorf("check the file, there is a transaction with amount 0")
	}
}
