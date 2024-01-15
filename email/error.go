package email

type errKind int

const (
	_ errKind = iota
	serviceNotEnabled
)

type EmailError struct {
	kind errKind
}

func (e EmailError) Error() string {
	switch e.kind {
	case serviceNotEnabled:
		return "email service is not enabled"
	}
	return "error with email"
}

var (
	ErrServiceNotEnabled = EmailError{kind: serviceNotEnabled}
)
