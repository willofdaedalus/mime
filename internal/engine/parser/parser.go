package parser

import (
	"errors"

	"willofdaedalus/mime/internal/engine/lexer"
)

type Parser struct {
	lex       *lexer.Lexer
	curToken  lexer.Token
	nextToken lexer.Token
	errs      []error
}

func NewParser(l *lexer.Lexer) *Parser {
	p := &Parser{
		lex: l,
	}
	// first call assigns the next token to nextToken
	// and the subsequent one assigns curToken to nextToken
	p.advanceToken()
	p.advanceToken()

	return p
}

func (p *Parser) advanceToken() {
	p.curToken = p.nextToken
	p.nextToken = p.lex.NextToken()
}

func (p *Parser) ParseTokens() {
	for p.curToken.Type != lexer.TokenEOF {
		if p.curToken.Type == lexer.TokenEntity {
			en := p.parseEntity()
		}
	}
}

func (p *Parser) pushError(msg string) {
	p.errs = append(p.errs, errors.New(msg))
}
