package dsl

import "fmt"

type node interface {
	// this will serve as the node's identifier
	NodeLiteral() string
}

type keywordHandler func(parser *Parser) (node, error)

var handlers = map[TokenType]keywordHandler{
	TokenEntity: handleEntity,
}

func NewParser(l *Lexer) *Parser {
	p := &Parser{
		lexer: l,
		nodes: make(map[string]node),
	}
	// first call assigns the next token to nextToken
	// and the subsequent one assigns curToken to nextToken
	p.advanceToken()
	p.advanceToken()

	return p
}

func (p *Parser) advanceToken() {
	p.curToken = p.nextToken
	p.nextToken = p.lexer.NextToken()
}

func (p *Parser) ParseTokens() error {
	var curLiteral string

	for p.curToken.Type != TokenEOF {
		curLiteral = p.curToken.Literal
		fmt.Printf("cur token is %s\n", curLiteral)

		if handler, ok := handlers[p.curToken.Type]; ok {
			v, err := handler(p)
			if err != nil {
				return err
			}

			p.nodes[curLiteral] = v
		}

		p.advanceToken()
	}

	return nil
}
