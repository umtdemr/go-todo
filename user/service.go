package user

import (
	"errors"
	"github.com/alexedwards/argon2id"
	"github.com/jackc/pgx/v5"
	"regexp"
)

type Service struct {
	repository IRepository
}

func NewUserService(repo IRepository) *Service {
	return &Service{repository: repo}
}

func (service *Service) CreateUser(data *CreateUserData) error {
	// validate username, email and password
	if userNameLength := len(data.Username); userNameLength < 3 || userNameLength > 20 {
		return ErrorUsernameLength
	}

	userNameRegex := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	if !userNameRegex.MatchString(data.Username) {
		return ErrorUserNameNotValidCharacters
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

	// check if the username exists
	userWithUsername := service.repository.GetUserByUsername(data.Username)
	if userWithUsername != nil {
		return ErrorUserNameNotValid
	}

	// check if the email exists
	userWithEmail := service.repository.GetUserByEmail(data.Email)
	if userWithEmail != nil {
		return ErrorEmailNotValid
	}

	// hash password to secure
	hash, err := argon2id.CreateHash(data.Password, argon2id.DefaultParams)
	if err != nil {
		return errors.New("error while hashing the password")
	}

	data.Password = hash

	return service.repository.CreateUser(data)
}

func (service *Service) Login(data *LoginUserData) (string, error) {
	if data.Password == nil {
		return "", ErrorPasswordLength
	}

	if data.Username == nil && data.Email == nil {
		return "", ErrorLoginIdEmpty
	}

	user, userQueryErr := service.repository.GetUserWithAllParams(data)

	if userQueryErr != nil {
		if errors.Is(userQueryErr, pgx.ErrNoRows) {
			return "", ErrorUsernameOrPasswordIncorrect
		}
		return "", userQueryErr
	}

	isPassMatched, hashErr := argon2id.ComparePasswordAndHash(*data.Password, user.Password)

	if hashErr != nil {
		return "", hashErr
	}

	if !isPassMatched {
		return "", ErrorUsernameOrPasswordIncorrect
	}

	// if credentials are correct, generate a token
	tokenString, tokenErr := GenerateNewJWT(user.Username)

	if tokenErr != nil {
		return "", tokenErr
	}

	return tokenString, nil
}
