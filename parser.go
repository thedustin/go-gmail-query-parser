package parser

import (
	"github.com/thedustin/go-gmail-query-parser/criteria"
	"github.com/thedustin/go-gmail-query-parser/lexer"
	"github.com/thedustin/go-gmail-query-parser/translator"
)

type Parser struct {
	lexer      *lexer.Lexer
	translator *translator.Translator
}

func NewParser() *Parser {
	return &Parser{
		lexer:      &lexer.Lexer{},
		translator: &translator.Translator{},
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
