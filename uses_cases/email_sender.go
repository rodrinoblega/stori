package uses_cases

import (
	"bytes"
	"fmt"
	"github.com/rodrinoblega/stori/entities"
	"html/template"
	"strconv"
)

type EmailSender interface {
	SendEmail(to string, subject string, body string) error
}

type EmailSummaryUseCase struct {
	emailSender EmailSender
	database    Database
}

type EmailData struct {
	TotalBalance          float64
	MonthYearTransactions []entities.MonthYearTransactions
	AverageSummary        entities.AverageSummary
}

func NewEmailSummaryUseCase(emailSender EmailSender, database Database) *EmailSummaryUseCase {
	return &EmailSummaryUseCase{
		emailSender: emailSender,
		database:    database,
	}
}

func (e *EmailSummaryUseCase) Execute(transactions entities.Transactions) error {
	account, err := e.getAccount(transactions)
	if err != nil {
		return err
	}

	emailData, err := e.buildEmailData(transactions)
	if err != nil {
		return err
	}

	htmlContent, err := e.generateEmailHTML(emailData)
	if err != nil {
		return err
	}

	return e.emailSender.SendEmail(account.Mail, "Transactions summary", htmlContent)
}

func (e *EmailSummaryUseCase) getAccount(transactions entities.Transactions) (*entities.Account, error) {
	accountID, err := transactions.GetAccountID()
	if err != nil {
		return nil, fmt.Errorf("failed to get account ID: %w", err)
	}

	num, err := strconv.Atoi(accountID)
	if err != nil {
		return nil, fmt.Errorf("invalid account ID format: %w", err)
	}

	account, err := e.database.GetAccountById(num)
	if err != nil {
		return nil, fmt.Errorf("account not found: %w", err)
	}

	return account, nil
}

func (e *EmailSummaryUseCase) buildEmailData(transactions entities.Transactions) (EmailData, error) {
	return EmailData{
		TotalBalance:          transactions.TotalBalance(),
		MonthYearTransactions: transactions.TransactionsByMonthYear(),
		AverageSummary:        transactions.AverageTransactionsAmount(),
	}, nil
}

func (e *EmailSummaryUseCase) generateEmailHTML(emailData EmailData) (string, error) {
	tmpl, err := template.ParseFiles("static/source.html")
	if err != nil {
		return "", fmt.Errorf("failed to parse email template: %w", err)
	}

	var htmlContent bytes.Buffer
	err = tmpl.Execute(&htmlContent, emailData)
	if err != nil {
		return "", fmt.Errorf("failed to generate email content: %w", err)
	}

	return htmlContent.String(), nil
}
