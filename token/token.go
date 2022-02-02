package token

import "strings"

type kind string

const (
	Start kind = "^"
	End   kind = "$"

	Field      kind = "%field%"
	Equal      kind = ":"
	FieldValue kind = "%value%"

	Fulltext kind = "%fulltext%"

	GroupStart kind = "("
	GroupEnd   kind = ")"

	Negate kind = "-"

	Or kind = "OR"
)

// trailingSpaceKinds is a map of tokens which needs a whitespace on printing.
var trailingSpaceKinds = map[kind]bool{
	Fulltext:   true,
	FieldValue: true,
	Or:         true,
	GroupEnd:   true,
}

// validationMap manages the allowed tokens for specific tokens. Not listed tokens support all token as descendant.
var validationMap = map[kind][]kind{
	Start: {Field, Fulltext, Negate, GroupStart, End},

	Field: {Equal},
	Equal: {FieldValue},

	Negate: {GroupStart, Field, Fulltext},
	Or:     {Field, Fulltext, Negate, GroupStart},

	GroupStart: {Field, Fulltext, Negate},
}

// Token is a simple struct containing information about the token type (kind), and the string value of this token.
type Token struct {
	kind  kind
	value string
}

// NewToken creates a new token.
func NewToken(kind kind, value string) Token {
	return Token{kind, value}
}

func (t *Token) Kind() kind {
	return t.kind
}

func (t *Token) Value() string {
	return t.value
}

// queryString transforms the token into its query representation value (should look like the input, but with pretty print)
func (t *Token) queryString() string {
	v := t.value

	if t.kind == FieldValue && strings.Contains(v, " ") {
		v = "(" + t.value + ")"
	}

	if !t.isTrailingSpaceToken() {
		return v
	}

	return v + " "
}

func (t *Token) isTrailingSpaceToken() bool {
	_, ok := trailingSpaceKinds[t.kind]

	return ok
}

// isValid checks whether the token is valid in relation if its descendant
func (t Token) isValid(next *Token) error {
	if t.kind == End && next == nil {
		return nil
	}

	kinds, ok := validationMap[t.kind]

	if !ok {
		// no requirements mean we are fine
		return nil
	}

	if !kindInList(next.kind, kinds) {
		return ValidationError{
			token:    t.kind,
			next:     next.kind,
			expected: kinds,
		}
	}

	return nil
}
