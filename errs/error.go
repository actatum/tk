package errs

import (
	"errors"
	"fmt"
)

// Error represents application specific errors. Application errors
// can be unwrapped by the caller to extract out the kind and message.
type Error struct {
	Kind    Kind
	Message string
}

// Error implements the error interface.
func (e *Error) Error() string {
	return fmt.Sprintf("kind=%s message=%s", e.Kind, e.Message)
}

// Kind defines what kind of error this is.
type Kind uint8

// Kinds of errors.
// New items must be added only to the end.
const (
	Internal Kind = iota
	AlreadyExists
	NotFound
)

func (k Kind) String() string {
	switch k {
	case AlreadyExists:
		return "already_exists"
	case NotFound:
		return "not_found"
	case Internal:
		return "internal"
	default:
		return "internal"
	}
}

// ErrorKind unwraps an application error and returns its kind.
// Non-application errors always return Internal.
func ErrorKind(err error) Kind {
	var e *Error
	if errors.As(err, &e) {
		return e.Kind
	}

	return Internal
}

// ErrorMessage unwarps an application error and returns its message.
// Non-application errors always return "internal error".
func ErrorMessage(err error) string {
	var e *Error
	if errors.As(err, &e) {
		return e.Message
	}

	return "internal error"
}

// Errorf is a helper function to return an Error with the given kind and formatted message.
func Errorf(kind Kind, format string, args ...any) *Error {
	return &Error{
		Kind:    kind,
		Message: fmt.Sprintf(format, args...),
	}
}
