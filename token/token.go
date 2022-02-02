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

var trailingSpaceKinds = map[kind]bool{
	Fulltext:   true,
	FieldValue: true,
	Or:         true,
	GroupEnd:   true,
}

var validationMap = map[kind][]kind{
	Start: {Field, Fulltext, Negate, GroupStart, End},

	Field: {Equal},
	Equal: {FieldValue},

	Negate: {GroupStart, Field, Fulltext},
	Or:     {Field, Fulltext, Negate, GroupStart},

	GroupStart: {Field, Fulltext, Negate},
}

type Token struct {
	kind  kind
	value string
}

func NewToken(kind kind, value string) Token {
	return Token{kind, value}
}

func (t *Token) Kind() kind {
	return t.kind
}

func (t *Token) Value() string {
	return t.value
}

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
