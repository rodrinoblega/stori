package uses_cases

import (
	"github.com/rodrinoblega/stori/entities"
	"strconv"
)

type EmailSender interface {
	SendEmail(to string, subject string, body string) error
}

type EmailSummaryUseCase struct {
	emailSender EmailSender
	database    Database
}

func NewEmailSummaryUseCase(emailSender EmailSender, database Database) *EmailSummaryUseCase {
	return &EmailSummaryUseCase{
		emailSender: emailSender,
		database:    database,
	}
}

func (e *EmailSummaryUseCase) Execute(transactions entities.Transactions) error {
	accountID, err := transactions.GetAccountID()
	if err != nil {
		return err
	}

	num, err := strconv.Atoi(accountID)
	if err != nil {
		return err
	}

	account, err := e.database.GetAccountById(num)
	if err != nil {
		return err
	}

	body := e.formatSummary("OK") // Format the email content
	return e.emailSender.SendEmail(account.Mail, "Transactions summary", body)
}

func (e *EmailSummaryUseCase) formatSummary(summary string) string {
	// Format the summary into an email-friendly body
	return "Transaction Summary:\n" + summary
}
