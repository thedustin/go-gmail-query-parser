package parser

import (
	"github.com/thedustin/go-gmail-query-parser/criteria"
	"github.com/thedustin/go-gmail-query-parser/lexer"
)

type Parser struct {
	lexer      *lexer.Lexer
	translator *criteria.Translator

	flags Flag
}

type Flag int

const (
	FlagOptimize Flag = 1 << iota

	FlagDefault = FlagOptimize
)

func NewParser(valFunc criteria.ValueTransformer, flags Flag) *Parser {
	return &Parser{
		lexer:      &lexer.Lexer{},
		translator: criteria.NewTranslator(valFunc).WithOptimizations(flags&FlagOptimize > 0),
		flags:      flags,
	}
}

func (p *Parser) Parse(query string) (criteria.Criteria, error) {
	if err := p.lexer.Parse(query); err != nil {
		return nil, err
	}

	tokens := p.lexer.Result()

	if err := tokens.Validate(); err != nil {
		return nil, err
	}

	crit, err := p.translator.ParseTree(tokens)

	if err != nil {
		return nil, err
	}

	return crit, nil
}
