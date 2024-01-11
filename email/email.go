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

func Init() {
	log := logger.Get()
	once.Do(func() {
		isEmailEnabled := viper.Get("EMAIL_ENABLED")
		if isEmailEnabled == "1" {
			for _, param := range EnvParams {
				if viper.Get(param) == nil {
					log.Fatal().Msg(fmt.Sprintf("environment variable %s is not set", param))
				}
			}

			log.Info().Msg("Email service is enabled")
			config = Config{
				IsEmailEnabled: true,
				From:           viper.Get("EMAIL_FROM").(string),
				Username:       viper.Get("EMAIL_USERNAME").(string),
				Password:       viper.Get("EMAIL_PASSWORD").(string),
				Host:           viper.Get("EMAIL_HOST").(string),
				Port:           viper.Get("EMAIL_PORT").(string),
			}
		} else {
			log.Info().Msg("Email service is not enabled")
			config = Config{
				IsEmailEnabled: false,
			}
		}
	})
}

func SenEmail(data SendEmailData) error {
	if !config.IsEmailEnabled {
		return fmt.Errorf("email service is not enabled")
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

	err := smtp.SendMail(addr, auth, config.From, data.To, byteMessage)
	if err != nil {
		log.Panic().Err(err).Msg("Couldn't send email")
	}

	log.Info().
		Str("to", receiverHeader).
		Str("subject", subject).
		Msg("Email sent")
	return nil
}
