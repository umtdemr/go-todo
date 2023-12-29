package user

type errKind int

const (
	_ errKind = iota
	lengthUsername
	userNameNotValid
	lengthEmail
	emailNotValid
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
	case userNameNotValid:
		return "username shouldn't include special characters"
	case lengthEmail:
		return "email length should be between 6 and 255"
	case emailNotValid:
		return "email is not valid"
	}
	return "error in user"
}

var (
	ErrorUsernameLength   = UserError{kind: lengthUsername, fields: []string{"username"}}
	ErrorUserNameNotValid = UserError{kind: userNameNotValid, fields: []string{"username"}}
	ErrorEmailLength      = UserError{kind: lengthEmail, fields: []string{"email"}}
	ErrorEmailNotValid    = UserError{kind: emailNotValid, fields: []string{"email"}}
)
