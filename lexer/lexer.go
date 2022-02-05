package lexer

import (
	"fmt"
	"strings"

	"github.com/thedustin/go-gmail-query-parser/token"
)

// See https://marianogappa.github.io/software/2019/06/05/lets-build-a-sql-parser-in-go/
// See https://www.youtube.com/watch?v=HxaD_trXwRE

type Lexer struct {
	source string
	result token.List

	fields fieldMap

	lastToken token.Token

	i     int
	group int
}

func NewLexer() *Lexer {
	l := &Lexer{}

	// @todo: Copy instead of Ref
	l.fields = defaultFields

	return l
}

func (l Lexer) Result() token.List {
	return l.result
}

func (l *Lexer) Reset() {
	l.source = ""
	l.result = nil
	l.result = append(l.result, token.NewToken(token.Start, string(token.Start)))

	l.lastToken = l.result[0]

	l.i = 0
	l.group = 0
}

// Parse transforms the given query into a list of tokens
func (l *Lexer) Parse(query string) error {
	l.Reset()
	l.source = query

	for l.i < len(query) {
		nextToken, _ := l.peek()

		if err := l.processToken(nextToken); err != nil {
			return err
		}

		l.pop()
	}

	if l.group > 0 {
		l.result = nil

		return ErrGroupNotClosed
	}

	l.result = append(l.result, token.NewToken(token.End, string(token.End)))

	return nil
}

func (l *Lexer) addToken(t token.Token) {
	l.result = append(l.result, t)
	l.lastToken = t
}

func (l *Lexer) processToken(t string) error {
	switch l.lastToken.Kind() {
	case token.Field:
		if t != ":" {
			newToken := token.NewToken(token.Fulltext, l.lastToken.Value())

			l.result[len(l.result)-1] = newToken

			l.lastToken = newToken

			return l.processToken(t)
		}

		l.addToken(token.NewToken(token.Equal, t))

		return nil
	case token.Equal:
		if t[0] == '(' {
			// todo: Unescape other "(", ")"
			t = t[1 : len(t)-1]
		}

		l.addToken(token.NewToken(token.FieldValue, t))

		return nil
	default:
		if _, ok := l.fields[t]; ok {
			l.addToken(token.NewToken(token.Field, t))

			return nil
		}

		if t == "OR" {
			l.addToken(token.NewToken(token.Or, t))

			return nil
		}

		if t == "AND" {
			// AND is only allowed as syntax sugar but not really used in our AST as it is the default behaviour
			l.lastToken = token.NewToken(token.Start, t)

			return nil
		}

		if t == "" {
			l.addToken(token.NewToken(token.GroupEnd, ")"))

			l.i++
			l.group--

			return nil
		}

		if t[0] == '-' {
			l.addToken(token.NewToken(token.Negate, "-"))
			l.i++

			nextToken, _ := l.peek()

			return l.processToken(nextToken)
		}

		if t[0] == '(' {
			l.addToken(token.NewToken(token.GroupStart, t[0:1]))
			l.i++
			l.group++

			nextToken, _ := l.peek()

			return l.processToken(nextToken)
		}

		l.addToken(token.NewToken(token.Fulltext, t))

		return nil
	}

	return fmt.Errorf("unknown state step %q for token: %q", l.lastToken, t)
}

// pop returns the next token, removes all following whitespace, and increases the index.
func (l *Lexer) pop() string {
	peeked, i := l.peek()

	l.i += i
	l.popWhitespace()

	return peeked
}

// popWhitespace skips all following whitespace/space
func (l *Lexer) popWhitespace() {
	for ; l.i < len(l.source) && l.source[l.i] == ' '; l.i++ {
		// do nothing, we skip all whitespace token by incrementing everything in the for-head
	}
}

// peek returns the next token and it length, without modifing the index
func (l Lexer) peek() (string, int) {
	if l.i >= len(l.source) {
		return "", 0
	}

	if l.source[l.i] == ':' {
		return ":", 1
	}

	for field := range l.fields {
		to := min(len(l.source), l.i+len(field))
		t := strings.ToLower(l.source[l.i:to])

		if t == field && l.source[to] == ':' {
			return t, len(t)
		}
	}

	if l.source[l.i] == '(' {
		return l.lookupNext([]byte{')'})
	}

	t, i := l.lookupNext(l.valueBoundaries())

	if i != 0 { // remove the space from the value
		return t[0:(i - 1)], i - 1
	}

	return t, len(t)
}

func (l Lexer) valueBoundaries() []byte {
	if l.group > 0 {
		return []byte{' ', ')'}
	}

	return []byte{' '}
}

// lookupNext searches for the next occurence of b byte and returns all content (including the b byte), and the length of the content.
// The length will be zero if the b byte was not found.
func (l Lexer) lookupNext(chars []byte) (string, int) {
	i := l.i

	for ; i < len(l.source); i++ {
		for _, b := range chars {
			if l.source[i] == b && l.source[i-1] != token.EscapeChar {
				t := l.source[l.i:(i + 1)]

				return token.Unescape(t), len(t)
			}
		}
	}

	return token.Unescape(l.source[l.i:i]), 0
}

func min(a, b int) int {
	if a < b {
		return a
	}

	return b
}
