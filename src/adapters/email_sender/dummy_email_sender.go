package email_sender

type DummyEmailSender struct{}

func NewDummyEmailSender() *DummyEmailSender {
	return &DummyEmailSender{}
}

func (n *DummyEmailSender) SendEmail(to string, subject string, body string) error {
	return nil
}
