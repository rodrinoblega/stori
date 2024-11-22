package uses_cases

import (
	"fmt"
	"github.com/rodrinoblega/stori/src/entities"
	"testing"
	"time"
)

func TestExecute(t *testing.T) {
	tests := []struct {
		name        string
		filePath    string
		expected    entities.Transactions
		expectedErr error
	}{
		{
			name:     "Valid file with transactions",
			filePath: "testdata/valid_transactions.csv",
			expected: entities.Transactions{
				{
					TransactionID:    "1",
					Date:             time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC),
					Amount:           100.00,
					TransactionsType: "CREDIT",
					AccountID:        "123",
				},
				{
					TransactionID:    "2",
					Date:             time.Date(2024, time.January, 2, 0, 0, 0, 0, time.UTC),
					Amount:           -50.00,
					TransactionsType: "DEBIT",
					AccountID:        "456",
				},
			},
			expectedErr: nil,
		},
		{
			name:        "Invalid CSV format",
			filePath:    "testdata/invalid_format.csv",
			expected:    nil,
			expectedErr: fmt.Errorf("invalid row format at line 2: [\"TransactionID\", \"Date\", \"Amount\"]"),
		},
		{
			name:        "Empty file",
			filePath:    "testdata/empty.csv",
			expected:    entities.Transactions{},
			expectedErr: nil,
		},
		{
			name:        "File with invalid amount",
			filePath:    "testdata/invalid_amount.csv",
			expected:    nil,
			expectedErr: fmt.Errorf("error processing row 2: invalid amount 'not_a_number': strconv.ParseFloat: parsing \"not_a_number\": invalid syntax"),
		},
		{
			name:        "File with invalid date",
			filePath:    "testdata/invalid_date.csv",
			expected:    nil,
			expectedErr: fmt.Errorf("error processing row 2: error parsing date: parsing time \"01-01-2024\" as \"01/02/2006\": cannot parse \"-01-2024\" as \"/\""),
		},
		{
			name:        "File with zero amount transaction",
			filePath:    "testdata/zero_amount.csv",
			expected:    nil,
			expectedErr: fmt.Errorf("error processing row 2: check the file, there is a transaction with amount 0"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fileReader := NewFileReaderUseCase()
			actual, err := fileReader.Execute(tt.filePath)
			if err != nil && err.Error() != tt.expectedErr.Error() {
				t.Errorf("expected error: %v, got: %v", tt.expectedErr, err)
			}
			if len(actual) != len(tt.expected) {
				t.Errorf("expected transactions length: %d, got: %d", len(tt.expected), len(actual))
			}
		})
	}
}
