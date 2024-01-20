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
		return ErrUsernameLength
	}

	userNameRegex := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	if !userNameRegex.MatchString(data.Username) {
		return ErrUserNameNotValidCharacters
	}

	if emailLength := len(data.Email); emailLength < 6 || emailLength > 255 {
		return ErrEmailLength
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

	if !emailRegex.MatchString(data.Email) {
		return ErrEmailNotValid
	}

	if passwordLength := len(data.Password); passwordLength < 8 || passwordLength > 64 {
		return ErrPasswordLength
	}

	// check if the username exists
	userWithUsername := service.repository.GetUserByUsername(data.Username)
	if userWithUsername != nil {
		return ErrUserNameNotValid
	}

	// check if the email exists
	userWithEmail := service.repository.GetUserByEmail(data.Email)
	if userWithEmail != nil {
		return ErrEmailNotValid
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
		return "", ErrPasswordLength
	}

	if data.Username == nil && data.Email == nil {
		return "", ErrLoginIdEmpty
	}

	user, userQueryErr := service.repository.GetUserWithAllParams(data)

	if userQueryErr != nil {
		if errors.Is(userQueryErr, pgx.ErrNoRows) {
			return "", ErrUsernameOrPasswordIncorrect
		}
		return "", userQueryErr
	}

	isPassMatched, hashErr := argon2id.ComparePasswordAndHash(*data.Password, user.Password)

	if hashErr != nil {
		return "", hashErr
	}

	if !isPassMatched {
		return "", ErrUsernameOrPasswordIncorrect
	}

	// if credentials are correct, generate a token
	tokenString, tokenErr := GenerateNewJWT(user.Username)

	if tokenErr != nil {
		return "", tokenErr
	}

	return tokenString, nil
}

func (service *Service) GenerateResetPasswordToken(data *ResetPasswordRequest) (string, error) {
	if data.Email == "" {
		return "", ErrEmailNotValid
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

	if !emailRegex.MatchString(data.Email) {
		return "", ErrEmailNotValid
	}

	user := service.repository.GetUserByEmail(data.Email)

	// we don't need to inform the user if there is not an user with the given email
	// see: https://cheatsheetseries.owasp.org/cheatsheets/Authentication_Cheat_Sheet.html
	if user == nil {
		return "", nil
	}

	tokenString, err := GenerateResetPasswordToken(user.Email)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (service *Service) ApplyNewPasswordWithToken(data *NewPasswordRequest) error {
	if data.Token == "" {
		return ErrTokenNotValid
	}

	if passwordLength := len(data.Password); passwordLength < 8 || passwordLength > 64 {
		return ErrPasswordLength
	}

	email, tokenErr := ValidateResetPasswordToken(data.Token)
	if tokenErr != nil {
		return tokenErr
	}

	user := service.repository.GetUserByEmail(email)

	if user == nil {
		return ErrUserNotFound
	}

	hashedPassword, err := argon2id.CreateHash(data.Password, argon2id.DefaultParams)
	if err != nil {
		return errors.New("error while hashing the password")
	}

	updateUserErr := service.repository.UpdateUserPassword(user.Id, hashedPassword)

	return updateUserErr
}
