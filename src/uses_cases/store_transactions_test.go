package uses_cases

import (
	"errors"
	"github.com/rodrinoblega/stori/src/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type MockDatabase struct {
	mock.Mock
}

func (m *MockDatabase) StoreTransactions(transactions entities.Transactions) error {
	args := m.Called(transactions)
	return args.Error(0)
}

func TestStoreTransactionsUseCase_Execute_Success(t *testing.T) {
	mockDatabase := new(MockDatabase)
	useCase := NewStoreTransactionsUseCase(mockDatabase)

	transactions := entities.Transactions{
		entities.Transaction{
			TransactionID:    "1",
			Amount:           100.0,
			TransactionsType: "CREDIT",
			AccountID:        "acc123",
		},
	}

	mockDatabase.On("StoreTransactions", transactions).Return(nil)

	err := useCase.Execute(transactions)

	assert.NoError(t, err)
	mockDatabase.AssertExpectations(t)
}

func TestStoreTransactionsUseCase_Execute_Failure(t *testing.T) {
	mockDatabase := new(MockDatabase)
	useCase := NewStoreTransactionsUseCase(mockDatabase)

	transactions := entities.Transactions{
		entities.Transaction{
			TransactionID:    "1",
			Amount:           100.0,
			TransactionsType: "CREDIT",
			AccountID:        "acc123",
		},
	}

	mockDatabase.On("StoreTransactions", transactions).Return(errors.New("database error"))

	err := useCase.Execute(transactions)

	assert.Error(t, err)
	assert.Equal(t, "database error", err.Error())
	mockDatabase.AssertExpectations(t)
}
