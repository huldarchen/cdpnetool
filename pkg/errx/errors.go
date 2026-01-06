package errx

import (
	"errors"
	"fmt"
)

type Code string

type Error struct {
	Code Code
	Msg  string
	Err  error
}

func (e *Error) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s: %v", e.Code, e.Msg, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Msg)
}

func (e *Error) Unwrap() error { return e.Err }

func New(code Code, msg string) *Error { return &Error{Code: code, Msg: msg} }

func Wrap(code Code, err error, msg string) *Error { return &Error{Code: code, Msg: msg, Err: err} }

func Is(err error, code Code) bool {
	var e *Error
	if errors.As(err, &e) {
		return e.Code == code
	}
	return false
}

const (
	CodeSessionNotFound Code = "SESSION_NOT_FOUND"
)
