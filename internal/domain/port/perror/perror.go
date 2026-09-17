package perror

import (
	"errors"
)

type Error struct {
	Type    Type
	Message string
}

func New(t Type, msg string) Error {
	return Error{
		Type:    t,
		Message: msg,
	}
}

func ParseError(err error) Error {
	if perr, ok := errors.AsType[Error](err); ok {
		return perr
	}

	return New(TypeUnspecified, "unspecified error")
}

func IsType(err error, t Type) bool {
	if perr, ok := errors.AsType[Error](err); ok {
		if perr.Type == t {
			return true
		}
	}

	return false
}

func (e Error) Error() string {
	return e.Message
}

func NotFound(msg string) Error {
	return New(TypeNotFound, msg)
}

func AlreadyExists(msg string) Error {
	return New(TypeAlreadyExists, msg)
}

func InvalidParam(msg string) Error {
	return New(TypeInvalidParam, msg)
}

func NoRowsUpdated() Error {
	return New(TypeNoRowsUpdated, "no rows updated")
}
