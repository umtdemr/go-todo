package server

import (
	"fmt"
	"strings"
)

type errKind int

const (
	_ errKind = iota
	notValidMethod
)

type ServerError struct {
	kind errKind
	info string
}

func (e ServerError) With(info string) ServerError {
	err := e
	err.info = info
	return err
}

func (e ServerError) Error() string {
	switch e.kind {
	case notValidMethod:
		errMessage := strings.Builder{}
		errMessage.WriteString("Method is not allowed")
		if e.info != "" {
			errMessage.WriteString(fmt.Sprintf(". %s", e.info))
		}
		return errMessage.String()
	}
	return "error with server"
}

var (
	ErrNotValidMethod = ServerError{kind: notValidMethod}
)
