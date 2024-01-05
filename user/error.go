package user

type errKind int

const (
	_ errKind = iota
	lengthUsername
	userNameNotValidCharacters
	userNameNotValid
	lengthEmail
	emailNotValid
	lengthPassword
	loginIdEmpty // in case username and email is empty
	jwtNotValid
	usernameOrPasswordWrong
	userNotFound
)

type UserError struct {
	kind   errKind
	value  int
	err    error
	fields []string
}

func (e UserError) Error() string {
	switch e.kind {
	case lengthUsername:
		return "username length should be between 3 and 30"
	case userNameNotValidCharacters:
		return "username shouldn't include special characters"
	case userNameNotValid:
		return "username is not valid"
	case lengthEmail:
		return "email length should be between 6 and 255"
	case emailNotValid:
		return "email is not valid"
	case lengthPassword:
		return "password length should be between 8 and 64"
	case loginIdEmpty:
		return "username or email need to be sent"
	case jwtNotValid:
		return "token is invalid"
	case usernameOrPasswordWrong:
		return "username or password is incorrect"
	case userNotFound:
		return "user not found"
	}
	return "error in user"
}

var (
	ErrUsernameLength              = UserError{kind: lengthUsername, fields: []string{"username"}}
	ErrUserNameNotValidCharacters  = UserError{kind: userNameNotValidCharacters, fields: []string{"username"}}
	ErrUserNameNotValid            = UserError{kind: userNameNotValid, fields: []string{"username"}}
	ErrEmailLength                 = UserError{kind: lengthEmail, fields: []string{"email"}}
	ErrEmailNotValid               = UserError{kind: emailNotValid, fields: []string{"email"}}
	ErrPasswordLength              = UserError{kind: lengthPassword, fields: []string{"password"}}
	ErrLoginIdEmpty                = UserError{kind: loginIdEmpty, fields: []string{"username", "email"}}
	ErrTokenNotValid               = UserError{kind: jwtNotValid}
	ErrUsernameOrPasswordIncorrect = UserError{kind: usernameOrPasswordWrong, fields: []string{"username", "password"}}
	ErrUserNotFound                = UserError{kind: userNotFound}
)
