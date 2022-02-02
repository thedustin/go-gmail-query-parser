package parser

import "github.com/thedustin/go-gmail-query-parser/lexer"

type Parser struct {
	lexer *lexer.Lexer
}

func NewParser() *Parser {
	return &Parser{
		&lexer.Lexer{},
	}
}

func (p *Parser) Parse(query string) error {
	if err := p.lexer.Parse(query); err != nil {
		return err
	}

	tokens := p.lexer.Result()

	if err := tokens.Validate(); err != nil {
		return err
	}

	return nil
}
