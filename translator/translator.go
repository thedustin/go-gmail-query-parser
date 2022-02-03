package translator

import (
	"fmt"

	"github.com/thedustin/go-gmail-query-parser/criteria"
	"github.com/thedustin/go-gmail-query-parser/token"
)

var ErrUnknownToken = fmt.Errorf("unknown token")
var ErrUnexpectedListEnd = fmt.Errorf("unexpected list end")
var ErrUnexpectedToken = fmt.Errorf("unexpected token")

type Translator struct {
	i      int
	tokens token.List

	lastCriteria criteria.Criteria
}

func (t *Translator) Reset() {
	t.i = 0
	t.tokens = nil
	t.lastCriteria = nil
}

func (t *Translator) ParseTree(ts token.List) (criteria.Criteria, error) {
	t.Reset()

	t.tokens = ts

	and := criteria.NewAnd()

	for t.i < len(ts) {
		tok := t.pop()

		if tok == nil {
			break
		}

		crit, err := t.criteriaFromToken(*tok)

		if err != nil {
			return nil, err
		}

		if crit != nil {
			t.lastCriteria = crit

			fmt.Println(*tok, "->", crit)
			and.Add(crit)
		}
	}

	return and, nil
}

func (t *Translator) peek() *token.Token {
	if t.i+1 < len(t.tokens) {
		return &t.tokens[t.i+1]
	}

	return nil
}

func (t *Translator) pop() *token.Token {
	tok := t.peek()

	if tok == nil {
		return nil
	}

	t.i++

	return tok
}

func (t *Translator) criteriaFromToken(tok token.Token) (criteria.Criteria, error) {
	switch tok.Kind() {
	case token.Start, token.End, token.GroupEnd:
		return nil, nil
	case token.Negate:
		nextTok := t.pop()

		if nextTok == nil {
			return nil, ErrUnexpectedListEnd
		}

		crit, err := t.criteriaFromToken(*nextTok)

		if err != nil {
			return nil, err
		}

		not := criteria.NewNot(crit)
		t.lastCriteria = not

		return not, nil
	case token.GroupStart:
		group := criteria.NewAnd()

		for nextTok := t.pop(); nextTok != nil; nextTok = t.pop() {
			crit, err := t.criteriaFromToken(*nextTok)

			if err != nil {
				return nil, err
			}

			if crit != nil {
				t.lastCriteria = crit
				group.Add(crit)
			}
		}

		t.lastCriteria = group

		return group, nil
	case token.Fulltext:
		// @todo: get criteria creator for token + irgendwie den Feldnamen übergeben und nen ValueTransformer...
		match := criteria.NewMatch(tok.Value())
		t.lastCriteria = match

		return match, nil
	case token.Field:
		eqTok := t.pop()

		if eqTok == nil {
			return nil, ErrUnexpectedListEnd
		}

		if eqTok.Kind() != token.Equal {
			return nil, ErrUnexpectedToken
		}

		valTok := t.pop()

		if valTok == nil {
			return nil, ErrUnexpectedListEnd
		}

		if valTok.Kind() != token.FieldValue {
			return nil, ErrUnexpectedToken
		}

		// @todo: get criteria creator for token + irgendwie den Feldnamen übergeben und nen ValueTransformer...
		match := criteria.NewMatch(valTok.Value())
		t.lastCriteria = match

		return match, nil
	case token.Or:
		left := t.lastCriteria
		nextTok := t.pop()

		if nextTok == nil {
			return nil, ErrUnexpectedListEnd
		}

		right, err := t.criteriaFromToken(*nextTok)

		if err != nil {
			return nil, err
		}

		or := criteria.NewOr(left, right)
		t.lastCriteria = or

		return nil, nil
	}

	return nil, ErrUnknownToken
}
