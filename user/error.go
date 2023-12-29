package user

type errKind int

const (
	_ errKind = iota
	lengthUsername
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
	}
	return "error in user"
}

var (
	ErrorUsernameLength = UserError{kind: lengthUsername, fields: []string{"username"}}
)
