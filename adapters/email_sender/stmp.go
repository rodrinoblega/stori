package email_sender

import (
	"fmt"
	"net/smtp"
)

type SMTPEmailSender struct {
	host     string
	port     string
	username string
	password string
}

func NewSMTPEmailSender(host, port, username, password string) *SMTPEmailSender {
	return &SMTPEmailSender{host: host, port: port, username: username, password: password}
}

func (s *SMTPEmailSender) SendEmail(to string, subject string, body string) error {
	auth := smtp.PlainAuth("", s.username, s.password, s.host)
	msg := []byte(fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\nContent-Type: text/html; charset=\"utf-8\"\r\n\r\n%s",
		"", to, subject, body))
	address := fmt.Sprintf("%s:%s", s.host, s.port)

	err := smtp.SendMail(address, auth, s.username, []string{to}, msg)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
