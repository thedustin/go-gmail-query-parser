package parser

import (
	"fmt"

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

	FlagErrorOnUnknownField // @todo: implement

	FlagDefault = FlagOptimize
)

func NewParser(valFunc criteria.ValueTransformer, flags Flag) *Parser {
	return &Parser{
		lexer:      lexer.NewLexer(),
		translator: criteria.NewTranslator(valFunc).WithOptimizations(flags&FlagOptimize > 0),
		flags:      flags,
	}
}

func (p *Parser) Parse(query string) (criteria.Criteria, error) {
	if err := p.lexer.Parse(query); err != nil {
		return nil, err
	}

	tokens := p.lexer.Result()

	fmt.Printf("%#v\n", FlagErrorOnUnknownField)

	if err := tokens.Validate(); err != nil {
		return nil, err
	}

	crit, err := p.translator.ParseTree(tokens)

	if err != nil {
		return nil, err
	}

	return crit, nil
}

func (p *Parser) AddField(name string, constructor criteria.CriteriaMatchConstructor) error {
	if err := p.lexer.AddField(name); err != nil {
		return err
	}

	p.translator.SetMatchConstructor(name, constructor)

	return nil
}

func (p *Parser) SetField(name string, constructor criteria.CriteriaMatchConstructor) {
	p.lexer.SetField(name, true)
	p.translator.SetMatchConstructor(name, constructor)
}

func (p *Parser) RemoveField(name string) error {
	if err := p.lexer.RemoveField(name); err != nil {
		return err
	}

	p.translator.SetMatchConstructor(name, nil)

	return nil
}
