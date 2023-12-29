package user

import "regexp"

type Service struct {
	repository IRepository
}

func NewUserService(repo IRepository) *Service {
	return &Service{repository: repo}
}

func (service *Service) CreateUser(data *CreateUserData) error {
	if userNameLength := len(data.Username); userNameLength < 3 || userNameLength > 20 {
		return ErrorUsernameLength
	}

	userNameRegex := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	if !userNameRegex.MatchString(data.Username) {
		return ErrorUserNameNotValid
	}

	if emailLength := len(data.Email); emailLength < 6 || emailLength > 255 {
		return ErrorEmailLength
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

	if !emailRegex.MatchString(data.Email) {
		return ErrorEmailNotValid
	}

	if passwordLength := len(data.Password); passwordLength < 8 || passwordLength > 64 {
		return ErrorPasswordLength
	}

	return service.repository.CreateUser(data)
}

func (service *Service) Login(data *LoginUserData) bool {
	return service.repository.Login(data)
}
