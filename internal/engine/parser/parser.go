package parser

import (
	"errors"
	"fmt"

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
		switch p.curToken.Type {
		case lexer.TokenEntity:
			p.advanceToken()
			entity := p.parseEntity()
			if entity != nil {
				// Store the entity or process it further
			}
		// case lexer.TokenAlter:
		// 	// Handle alter statements
		// 	p.parseAlter()
		// case lexer.TokenRoutes:
		// 	// Handle routes statements
		// 	p.parseRoutes()
		// case lexer.TokenNewline, lexer.TokenComment:
		// 	// Skip newlines and comments
		// 	p.advanceToken()
		default:
			p.pushError(fmt.Sprintf("unexpected token: %s", p.curToken.Type))
			p.advanceToken()
		}
	}
}

func (p *Parser) pushError(msg string) {
	p.errs = append(p.errs, errors.New(msg))
}
