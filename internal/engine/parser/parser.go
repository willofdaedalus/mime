package parser

import (
	"errors"
	"fmt"

	"willofdaedalus/mime/internal/engine/lexer"
)

type keywordHandler func(parser *Parser) node

var handlers = map[lexer.TokenType]keywordHandler{
	lexer.TokenEntity: handleEntity,
	lexer.TokenEnum:   handleEnum,
}

type node interface {
	// this will serve as the node's identifier
	NodeLiteral() string
}

type Parser struct {
	lex          *lexer.Lexer
	curToken     lexer.Token
	nextToken    lexer.Token
	parserErrors []parserError
	errors       []error
	nodes        map[string]node
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

// TODO!
// name handler
// { "entity", entityHandler() }
func (p *Parser) ParseTokens() {
	for p.curToken.Type != lexer.TokenEOF {
		if handler, ok := handlers[p.curToken.Type]; ok {
			p.advanceToken()
			v := handler(p)
			if v != nil {
				p.nodes[p.curToken.Literal] = v
			}
		}
	}
}

func (p *Parser) findEntityNode(name string) (*entityNode, error) {
	if e, ok := p.nodes[name]; ok {
		return e.(*entityNode), nil
	}

	return nil, fmt.Errorf("entity of name %s doesn't exist in this context", name)
}

func (p *Parser) addError(logLevel parserErrorLevel, msg string) {
	err := parserError{
		errorLevel: logLevel,
		msg:        msg,
	}

	p.parserErrors = append(p.parserErrors, err)
}

func (p *Parser) pushError(msg string) {
	// fmt.Println(msg)
	p.errors = append(p.errors, errors.New(msg))
}
