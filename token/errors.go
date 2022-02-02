package token

import (
	"errors"
	"fmt"
)

type ValidationError struct {
	token    kind
	next     kind
	expected []kind
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("token of type %q expected on of %q as next token but got %q", e.token, e.expected, e.next)
}

func (e ValidationError) Unwrap() error {
	return ErrValidation
}

var ErrValidation = errors.New("validation error")
