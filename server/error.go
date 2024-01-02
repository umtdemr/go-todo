package server

import (
	"fmt"
	"strings"
)

type errKind int

const (
	_ errKind = iota
	notValidMethod
	invalidRequest
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

func (e ServerError) GetMessageWithInfo(actualMsg string) string {
	errMessage := strings.Builder{}
	errMessage.WriteString(actualMsg)
	if e.info != "" {
		errMessage.WriteString(fmt.Sprintf(". %s", e.info))
	}
	return errMessage.String()
}

func (e ServerError) Error() string {
	switch e.kind {
	case notValidMethod:
		return e.GetMessageWithInfo("method is not allowed")
	case invalidRequest:
		return e.GetMessageWithInfo("invalid request")
	}
	return "error with server"
}

var (
	ErrNotValidMethod = ServerError{kind: notValidMethod}
	ErrInvalidRequest = ServerError{kind: invalidRequest}
)
