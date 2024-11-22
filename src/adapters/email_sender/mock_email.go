package email_sender

type MockTestMail struct{}

func NewMockTestMail() *MockTestMail {
	return &MockTestMail{}
}

func (n *MockTestMail) SendEmail(to string, subject string, body string) error {
	return nil
}
