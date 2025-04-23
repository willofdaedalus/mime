package parser

import (
	"fmt"

	"willofdaedalus/mime/internal/engine/lexer"
)

func (p *Parser) parseAlter() {
	if !expectTokOf(p.curToken, lexer.TokenAlter) {
		p.pushError(fmt.Sprintf("expected alter token, got %s", p.curToken.Type))
	}
	p.advanceToken() // consume "alter"

	if !expectTokOf(p.curToken, lexer.TokenRef) {
		p.pushError(fmt.Sprintf("expected ref keyword got %s", p.curToken.Type))
	}
	p.advanceToken() // consume "ref"

	// get the entity that matches the name; based on what comes after the dot, we'll
	// either modify the payload or the entity's response
	_, err := p.findEntityNode(p.curToken.Literal)
	if err != nil {
		p.pushError(err.Error())
		return
	}
}
