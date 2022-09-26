package errors

import (
	"fmt"
)

var (
	ErrInvalidArgument = fmt.Errorf("INVALID_ARG_ERROR")
	ErrNotFound        = fmt.Errorf("NOT_FOUND_ERROR")
	ErrUnexpected      = fmt.Errorf("UNEXPECTED_ERORR")
	ErrUnauthorized    = fmt.Errorf(("UNAUTHORIZED"))
	ErrAlreadyExist    = fmt.Errorf("ALREADY_EXIST_ERROR")
)

type Error struct {
	Code    int
	Message string
	Err     error
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s: %s", e.Err, e.Message)
}

func (e *Error) Unwrap() error { return e.Err }

func NewInvalidArgumentErr(msg string) *Error {
	return &Error{
		Code:    1,
		Message: msg,
		Err:     ErrInvalidArgument,
	}
}

func NewNotFoundErr(msg string) *Error {
	return &Error{
		Code:    2,
		Message: msg,
		Err:     ErrNotFound,
	}
}

func NewUnexpectedErr(msg string) *Error {
	return &Error{
		Code:    3,
		Message: msg,
		Err:     ErrUnexpected,
	}
}

func NewUnAuthorized(msg string) *Error {
	return &Error{
		Code:    4,
		Message: msg,
		Err:     ErrUnauthorized,
	}
}

func NewAlreadyExistErr(msg string) *Error {
	return &Error{
		Code:    5,
		Message: msg,
		Err:     ErrAlreadyExist,
	}
}
