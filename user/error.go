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
	}
	return "error in user"
}

var (
	ErrorUsernameLength             = UserError{kind: lengthUsername, fields: []string{"username"}}
	ErrorUserNameNotValidCharacters = UserError{kind: userNameNotValidCharacters, fields: []string{"username"}}
	ErrorUserNameNotValid           = UserError{kind: userNameNotValid, fields: []string{"username"}}
	ErrorEmailLength                = UserError{kind: lengthEmail, fields: []string{"email"}}
	ErrorEmailNotValid              = UserError{kind: emailNotValid, fields: []string{"email"}}
	ErrorPasswordLength             = UserError{kind: lengthPassword, fields: []string{"password"}}
	ErrorLoginIdEmpty               = UserError{kind: loginIdEmpty, fields: []string{"username", "email"}}
)
