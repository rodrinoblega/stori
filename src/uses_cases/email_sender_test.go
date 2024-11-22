package uses_cases

import (
	"errors"
	"github.com/rodrinoblega/stori/src/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

type MockEmailSender struct {
	mock.Mock
}

func (m *MockEmailSender) SendEmail(to string, subject string, body string) error {
	args := m.Called(to, subject, body)
	return args.Error(0)
}

func (m *MockDatabase) GetAccountById(id int) (*entities.Account, error) {
	args := m.Called(id)
	account := args.Get(0)
	if account == nil {
		return nil, args.Error(1)
	}
	return account.(*entities.Account), args.Error(1)
}

func TestEmailSummaryUseCase_Execute_Success(t *testing.T) {
	mockEmailSender := new(MockEmailSender)
	mockDatabase := new(MockDatabase)

	account := &entities.Account{
		Mail: "test@example.com",
	}
	mockDatabase.On("GetAccountById", 123).Return(account, nil)

	mockEmailSender.On("SendEmail", "test@example.com", "Transactions summary", mock.Anything).Return(nil)

	useCase := NewEmailSummaryUseCase(mockEmailSender, mockDatabase, "testdata/source.html")

	transactions := entities.Transactions{
		{
			TransactionID:    "1",
			Date:             time.Now(),
			Amount:           100.0,
			TransactionsType: "CREDIT",
			AccountID:        "123",
		},
	}

	err := useCase.Execute(transactions)

	assert.NoError(t, err)

	mockEmailSender.AssertExpectations(t)
	mockDatabase.AssertExpectations(t)
}

func TestEmailSummaryUseCase_Execute_GetAccountError(t *testing.T) {
	mockEmailSender := new(MockEmailSender)
	mockDatabase := new(MockDatabase)

	mockDatabase.On("GetAccountById", 123).Return(nil, errors.New("account not found"))

	useCase := NewEmailSummaryUseCase(mockEmailSender, mockDatabase, "testdata/source.html")

	transactions := entities.Transactions{
		{
			TransactionID:    "1",
			Date:             time.Now(),
			Amount:           100.0,
			TransactionsType: "CREDIT",
			AccountID:        "123",
		},
	}

	err := useCase.Execute(transactions)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "account not found")

	mockEmailSender.AssertExpectations(t)
	mockDatabase.AssertExpectations(t)
}

func TestEmailSummaryUseCase_Execute_SendEmailError(t *testing.T) {
	mockEmailSender := new(MockEmailSender)
	mockDatabase := new(MockDatabase)

	account := &entities.Account{
		Mail: "test@example.com",
	}
	mockDatabase.On("GetAccountById", 123).Return(account, nil)

	mockEmailSender.On("SendEmail", "test@example.com", "Transactions summary", mock.Anything).Return(errors.New("email sending failed"))

	useCase := NewEmailSummaryUseCase(mockEmailSender, mockDatabase, "testdata/source.html")

	transactions := entities.Transactions{
		{
			TransactionID:    "1",
			Date:             time.Now(),
			Amount:           100.0,
			TransactionsType: "CREDIT",
			AccountID:        "123",
		},
	}

	err := useCase.Execute(transactions)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "email sending failed")

	mockEmailSender.AssertExpectations(t)
	mockDatabase.AssertExpectations(t)
}
