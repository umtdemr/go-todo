package email

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/smtp"
	"testing"
)

type mockEmailSender struct {
	mock.Mock
}

type mockViper struct {
	mock.Mock
}

func (m *mockEmailSender) SendMail(addr string, a smtp.Auth, from string, to []string, msg []byte) error {
	args := m.Called(addr, a, from, to, msg)
	return args.Error(0)
}

func (m *mockViper) Get(key string) interface{} {
	args := m.Called(key)
	return args.Get(0)
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

	assert.Equal(t, ErrServiceNotEnabled, err)
}

func TestSenEmail(t *testing.T) {
	mockSender := new(mockEmailSender)
	originalEmailSender := emailSender
	emailSender = mockSender
	defer func() { emailSender = originalEmailSender }()

	tests := []struct {
		name              string
		data              SendEmailData
		mockedResponseErr error
	}{
		{
			name: "should successfully send email",
			data: SendEmailData{
				To:      []string{"to"},
				Subject: "subject",
				Message: "message",
			},
			mockedResponseErr: nil,
		},
		{
			name: "should return the error if email service fails",
			data: SendEmailData{
				To:      []string{"to"},
				Subject: "subject",
				Message: "message",
			},
			mockedResponseErr: ErrServiceNotEnabled, // any error
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			config = Config{
				IsEmailEnabled: true,
				From:           "from",
				Username:       "username",
				Password:       "password",
				Host:           "host",
				Port:           "port",
			}
			mockSender.On("SendMail", mock.Anything, mock.Anything, mock.Anything, tc.data.To, mock.Anything).Return(tc.mockedResponseErr)

			err := SenEmail(tc.data)

			assert.Equal(t, tc.mockedResponseErr, err)

			mockSender.ExpectedCalls = nil
			mockSender.Calls = nil
		})
	}
}

// TestInit tests that Init sets the config correctly
func TestInit(t *testing.T) {
	type test struct {
		name           string
		envData        map[string]string
		expectedConfig Config
	}
	tests := []test{
		{
			name: "email not enabled",
			envData: map[string]string{
				"EMAIL_ENABLED": "0",
			},
			expectedConfig: Config{
				IsEmailEnabled: false,
			},
		},
		{
			name: "email enabled but missing env params",
			envData: map[string]string{
				"EMAIL_ENABLED": "1",
			},
			expectedConfig: Config{
				IsEmailEnabled: false,
			},
		},
		{
			name: "email enabled and all env params set",
			envData: map[string]string{
				"EMAIL_ENABLED":  "1",
				"EMAIL_FROM":     "from",
				"EMAIL_USERNAME": "username",
				"EMAIL_PASSWORD": "password",
				"EMAIL_HOST":     "host",
				"EMAIL_PORT":     "port",
			},
			expectedConfig: Config{
				IsEmailEnabled: true,
				From:           "from",
				Username:       "username",
				Password:       "password",
				Host:           "host",
				Port:           "port",
			},
		},
	}

	mockedViper := new(mockViper)
	originalViperGetterInstance := viperGetterInstance
	viperGetterInstance = mockedViper

	defer func() { viperGetterInstance = originalViperGetterInstance }()

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			for key, value := range tc.envData {
				mockedViper.Mock.On("Get", key).Return(value)
			}
			mockedViper.On("Get", mock.Anything).Return(nil)

			initializeConfig()
			assert.Equal(t, tc.expectedConfig, config)

			mockedViper.ExpectedCalls = nil
			mockedViper.Calls = nil
		})
	}
}
