package lexer

import (
	"errors"
	"fmt"
)

var ErrSyntaxError = errors.New("syntax error")

var ErrGroupNotClosed = fmt.Errorf("%w: group was not closed", ErrSyntaxError)
