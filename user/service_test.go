package user

import (
	"github.com/alexedwards/argon2id"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"strings"
	"testing"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) GetUserByUsername(username string) *VisibleUser {
	args := m.Called(username)
	if args.Get(0) != nil {
		return args.Get(0).(*VisibleUser)
	}
	return nil
}

func (mock *MockRepository) GetUserByEmail(email string) *VisibleUser {
	args := mock.Called(email)
	if args.Get(0) != nil {
		return args.Get(0).(*VisibleUser)
	}
	return nil
}

func (mock *MockRepository) CreateUser(data *CreateUserData) error {
	args := mock.Called(data)
	return args.Error(0)
}

func (m *MockRepository) GetUserWithAllParams(data *LoginUserData) (*UserParams, error) {
	args := m.Called(data)
	if args.Get(0) != nil {
		return args.Get(0).(*UserParams), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockRepository) UpdateUserPassword(userId int64, newPassword string) error {
	args := m.Called(userId, newPassword)
	return args.Error(0)
}

func TestCreateUser(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewUserService(mockRepo)

	validUsername := "username"
	validEmail := "username@ddd.com"
	validPassword := "validPassword"

	tests := []struct {
		name          string
		input         *CreateUserData
		setupMock     func()
		expectedError error
	}{
		{
			name: "Username length is less than 3",
			input: &CreateUserData{
				Username: "us",
			},
			setupMock:     func() {},
			expectedError: ErrUsernameLength,
		},
		{
			name: "Username length is more than 20",
			input: &CreateUserData{
				Username: strings.Repeat("s", 21),
			},
			setupMock:     func() {},
			expectedError: ErrUsernameLength,
		},
		{
			name: "Username contains invalid characters",
			input: &CreateUserData{
				Username: "   ",
			},
			setupMock:     func() {},
			expectedError: ErrUserNameNotValidCharacters,
		},
		{
			name: "Email length is less than 6",
			input: &CreateUserData{
				Username: validUsername,
				Email:    "s@s.c",
			},
			setupMock:     func() {},
			expectedError: ErrEmailLength,
		},
		{
			name: "Email length is more than 255",
			input: &CreateUserData{
				Username: validUsername,
				Email:    strings.Repeat("s", 256),
			},
			setupMock:     func() {},
			expectedError: ErrEmailLength,
		},
		{
			name: "Email length is less than 8",
			input: &CreateUserData{
				Username: validUsername,
				Email:    "ss",
			},
			setupMock:     func() {},
			expectedError: ErrEmailLength,
		},
		{
			name: "Email is invalid",
			input: &CreateUserData{
				Username: validUsername,
				Email:    "invalidEmail",
			},
			setupMock:     func() {},
			expectedError: ErrEmailNotValid,
		},
		{
			name: "Password length is less then 8",
			input: &CreateUserData{
				Username: validUsername,
				Email:    validEmail,
				Password: "sss",
			},
			setupMock:     func() {},
			expectedError: ErrPasswordLength,
		},
		{
			name: "Password length is more than 64",
			input: &CreateUserData{
				Username: validUsername,
				Email:    validEmail,
				Password: strings.Repeat("s", 65),
			},
			setupMock:     func() {},
			expectedError: ErrPasswordLength,
		},
		{
			name: "Username exists",
			input: &CreateUserData{
				Username: validUsername,
				Email:    validEmail,
				Password: validPassword,
			},
			setupMock: func() {
				mockRepo.On("GetUserByUsername", validUsername).Return(&VisibleUser{})
			},
			expectedError: ErrUserNameNotValid,
		},
		{
			name: "Email exists",
			input: &CreateUserData{
				Username: validUsername,
				Email:    validEmail,
				Password: validPassword,
			},
			setupMock: func() {
				mockRepo.On("GetUserByUsername", validUsername).Return(nil)
				mockRepo.On("GetUserByEmail", validEmail).Return(&VisibleUser{})
			},
			expectedError: ErrEmailNotValid,
		},
		{
			name: "Valid input",
			input: &CreateUserData{
				Username: validUsername,
				Email:    validEmail,
				Password: validPassword,
			},
			setupMock: func() {
				mockRepo.On("GetUserByUsername", validUsername).Return(nil)
				mockRepo.On("GetUserByEmail", validEmail).Return(nil)
				mockRepo.On("CreateUser", mock.Anything).Return(nil)
			},
			expectedError: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMock()
			err := service.CreateUser(tc.input)
			assert.Equal(t, tc.expectedError, err)

			mockRepo.ExpectedCalls = nil
		})
	}
}

func TestLogin(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewUserService(mockRepo)

	passWord := "password"

	tests := []struct {
		name          string
		input         *LoginUserData
		setupMock     func()
		expectedError error
	}{
		{
			name:          "Password is nil",
			input:         &LoginUserData{},
			setupMock:     func() {},
			expectedError: ErrPasswordLength,
		},
		{
			name: "Email and Username are nil at the same time",
			input: &LoginUserData{
				Password: new(string),
			},
			setupMock:     func() {},
			expectedError: ErrLoginIdEmpty,
		},
		{
			name: "No user found with the given username or email",
			input: &LoginUserData{
				Password: new(string),
				Username: new(string),
			},
			setupMock: func() {
				mockRepo.On("GetUserWithAllParams", mock.Anything).Return(nil, pgx.ErrNoRows)
			},
			expectedError: ErrUsernameOrPasswordIncorrect,
		},
		{
			name: "Password is incorrect",
			input: &LoginUserData{
				Username: new(string),
				Password: &passWord,
			},
			setupMock: func() {
				hashedPass, _ := argon2id.CreateHash("inCorrect", argon2id.DefaultParams)
				mockRepo.On("GetUserWithAllParams", mock.Anything).Return(&UserParams{Password: hashedPass}, nil)
			},
			expectedError: ErrUsernameOrPasswordIncorrect,
		},
		{
			name: "Login with username is successful",
			input: &LoginUserData{
				Username: new(string),
				Password: &passWord,
			},
			setupMock: func() {
				hashedPass, _ := argon2id.CreateHash("password", argon2id.DefaultParams)
				mockRepo.On("GetUserWithAllParams", mock.Anything).Return(&UserParams{Password: hashedPass}, nil)
			},
			expectedError: nil,
		},
		{
			name: "Login with email is successful",
			input: &LoginUserData{
				Email:    new(string),
				Password: &passWord,
			},
			setupMock: func() {
				hashedPass, _ := argon2id.CreateHash("password", argon2id.DefaultParams)
				mockRepo.On("GetUserWithAllParams", mock.Anything).Return(&UserParams{Password: hashedPass}, nil)
			},
			expectedError: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMock()
			_, err := service.Login(tc.input)
			assert.Equal(t, tc.expectedError, err)

			mockRepo.ExpectedCalls = nil
		})
	}
}
