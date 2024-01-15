package email

import (
	"errors"
	"github.com/stretchr/testify/mock"
	"net/smtp"
	"testing"
)

type mockEmailSender struct {
	mock.Mock
}

func (m *mockEmailSender) SendMail(addr string, a smtp.Auth, from string, to []string, msg []byte) error {
	args := m.Called(addr, a, from, to, msg)
	return args.Error(0)
}

// TestSenEmail_ShowErrorIfNotEnabled tests that SendEmail returns an error if email service is not enabled
func TestSenEmail_ShowErrorIfNotEnabled(t *testing.T) {
	mockSender := new(mockEmailSender)
	originalEmailSender := emailSender
	emailSender = mockSender
	defer func() { emailSender = originalEmailSender }()

	mockSender.On("SendMail", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	emailData := SendEmailData{
		To:      []string{"to"},
		Subject: "subject",
		Message: "message",
	}

	err := SenEmail(emailData)

	if err == nil || !errors.Is(err, ErrServiceNotEnabled) {
		t.Errorf("Error while sending email: %s", err)
	}
}
