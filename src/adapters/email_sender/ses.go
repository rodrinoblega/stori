package email_sender

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/ses"
)

type SESEmailSender struct {
	client    *ses.SES
	fromEmail string
}

func NewSESEmailSender(client *ses.SES, fromEmail string) *SESEmailSender {
	return &SESEmailSender{client: client, fromEmail: fromEmail}
}

func (s *SESEmailSender) SendEmail(to string, subject string, body string) error {
	rawMessage := &ses.RawMessage{
		Data: []byte(fmt.Sprintf(
			"From: %s\r\nTo: %s\r\nSubject: %s\r\nContent-Type: text/html; charset=\"utf-8\"\r\n\r\n%s",
			s.fromEmail, to, subject, body,
		)),
	}

	emailInput := &ses.SendRawEmailInput{
		RawMessage: rawMessage,
	}

	_, err := s.client.SendRawEmail(emailInput)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
