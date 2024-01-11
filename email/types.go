package email

type Config struct {
	IsEmailEnabled bool
	From           string
	Username       string
	Password       string
	Host           string
	Port           string
}

type SendEmailData struct {
	To      []string
	Subject string
	Message string
}
