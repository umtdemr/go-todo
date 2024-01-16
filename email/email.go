package email

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/umtdemr/go-todo/logger"
	"net/smtp"
	"sync"
)

var config Config
var once sync.Once

var EnvParams = []string{
	"EMAIL_HOST",
	"EMAIL_PORT",
	"EMAIL_FROM",
	"EMAIL_USERNAME",
	"EMAIL_PASSWORD",
}

// Sender is an interface for sending emails
// can be mocked for testing
type Sender interface {
	SendMail(addr string, a smtp.Auth, from string, to []string, msg []byte) error
}

type EnvGetter interface {
	Get(key string) any
}

type defaultEmailSender struct{}
type viperGetter struct{}

func (v *viperGetter) Get(key string) any {
	return viper.Get(key)
}

var viperGetterInstance EnvGetter = &viperGetter{}

// SendMail is a wrapper for smtp.SendMail for mocking purposes
func (s *defaultEmailSender) SendMail(addr string, a smtp.Auth, from string, to []string, msg []byte) error {
	return smtp.SendMail(addr, a, from, to, msg)
}

var emailSender Sender = &defaultEmailSender{}

func Init() {
	once.Do(initializeConfig)
}

func initializeConfig() {
	log := logger.Get()
	isEmailEnabled := viperGetterInstance.Get("EMAIL_ENABLED")
	if isEmailEnabled == "1" {
		for _, param := range EnvParams {
			var isParamNil bool
			if viperGetterInstance.Get(param) == nil {
				isParamNil = true
				log.Info().Msg(fmt.Sprintf("environment variable %s is not set", param))
			}
			if isParamNil {
				config = Config{
					IsEmailEnabled: false,
				}
				return
			}
		}

		log.Info().Msg("Email service is enabled")
		config = Config{
			IsEmailEnabled: true,
			From:           viperGetterInstance.Get("EMAIL_FROM").(string),
			Username:       viperGetterInstance.Get("EMAIL_USERNAME").(string),
			Password:       viperGetterInstance.Get("EMAIL_PASSWORD").(string),
			Host:           viperGetterInstance.Get("EMAIL_HOST").(string),
			Port:           viperGetterInstance.Get("EMAIL_PORT").(string),
		}
	} else {
		log.Info().Msg("Email service is not enabled")
		config = Config{
			IsEmailEnabled: false,
		}
	}
}

func Send(data SendEmailData) error {
	if !config.IsEmailEnabled {
		return ErrServiceNotEnabled
	}
	log := logger.Get()

	addr := config.Host + ":" + config.Port

	subject := data.Subject
	body := data.Message

	auth := smtp.PlainAuth("", config.Username, config.Password, config.Host)

	receiverHeader := ""
	for _, receiver := range data.To {
		receiverHeader += receiver + ", "
	}
	receiverHeader = receiverHeader[:len(receiverHeader)-2]

	header := ""
	header += fmt.Sprintf("From: %s\r\n", config.From)
	header += fmt.Sprintf("To: %s\r\n", receiverHeader)
	header += fmt.Sprintf("Subject: %s\r\n", subject)
	header += "\r\n" // Separate header from body
	message := header + body

	byteMessage := []byte(message)

	err := emailSender.SendMail(addr, auth, config.From, data.To, byteMessage)
	if err != nil {
		log.Error().Err(err).Msg("Couldn't send email")
		return err
	}

	log.Info().
		Str("to", receiverHeader).
		Str("subject", subject).
		Msg("Email sent")
	return nil
}
