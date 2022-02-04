package criteria

import (
	"fmt"

	"github.com/thedustin/go-gmail-query-parser/token"
)

var ErrUnknownToken = fmt.Errorf("unknown token")
var ErrUnexpectedListEnd = fmt.Errorf("unexpected list end")
var ErrUnexpectedToken = fmt.Errorf("unexpected token")

type ValueTransformer func(field string, v interface{}) []string

type Translator struct {
	i      int
	tokens token.List

	optimize bool

	lastCriteria InnerCriteria

	valFunc ValueTransformer
}

func NewTranslator(valFunc ValueTransformer) *Translator {
	return &Translator{
		valFunc: valFunc,
	}
}

func (t *Translator) WithOptimizations(enabled bool) *Translator {
	t.optimize = enabled

	return t
}

func (t *Translator) Reset() {
	t.i = 0
	t.tokens = nil
	t.lastCriteria = nil
}

func (t *Translator) ParseTree(ts token.List) (Criteria, error) {
	t.Reset()

	t.tokens = ts

	and := NewAnd()
	and.SetParent(and)

	t.lastCriteria = and

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

	if !t.optimize {
		return and, nil
	}

	finalCrit, err := t.optimizeCriteria(and)

	if err != nil {
		return nil, err
	}

	return finalCrit, nil
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

func (t *Translator) criteriaFromToken(tok token.Token) (InnerCriteria, error) {
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

		not := NewNot(crit)
		not.SetParent(t.lastCriteria.Parent())
		t.lastCriteria = not

		return not, nil
	case token.GroupStart:
		group := NewAnd()
		group.SetParent(t.lastCriteria.Parent())

		t.lastCriteria = group

		for nextTok := t.pop(); nextTok != nil; nextTok = t.pop() {
			if nextTok.Kind() == token.GroupEnd {
				break
			}

			crit, err := t.criteriaFromToken(*nextTok)

			if err != nil {
				return nil, err
			}

			if crit != nil {
				t.lastCriteria = crit
				group.Add(crit)
			}
		}

		return group, nil
	case token.Fulltext:
		// @todo: get criteria creator for token + irgendwie den Feldnamen übergeben und nen ValueTransformer...
		match := NewMatch(FieldFulltext, tok.Value(), t.valFunc)
		match.SetParent(t.lastCriteria.Parent())
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
		match := NewMatch(tok.Value(), valTok.Value(), t.valFunc)
		match.SetParent(t.lastCriteria.Parent())
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

		or := NewOr(left, right)

		parentCrit, ok := t.lastCriteria.Parent().(GroupCriteria)

		if !ok {
			fmt.Println(left)
			return nil, ErrUnexpectedToken
		}

		parentCrit.Replace(left, or)
		t.lastCriteria = or

		return nil, nil
	}

	return nil, ErrUnknownToken
}

func (t Translator) optimizeCriteria(c InnerCriteria) (InnerCriteria, error) {
	group, ok := c.(GroupCriteria)

	if !ok {
		return c, nil
	}

	switch group.Length() {
	case 0:
		return NewNoop(), nil
	case 1:
		return t.optimizeCriteria(group.All()[0])
	}

	for _, old := range group.All() {
		new, err := t.optimizeCriteria(old)

		if err != nil {
			return nil, err
		}

		if old == new {
			continue
		}

		group.Replace(old, new)
	}

	return c, nil
}

func isSameGroup(a, b InnerCriteria) bool {
	switch a.(type) {
	case *criteriaAnd:
		_, ok := b.(*criteriaAnd)

		return ok
	case *criteriaOr:
		_, ok := b.(*criteriaOr)

		return ok
	}

	return false
}
