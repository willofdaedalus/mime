package parser

import (
	"fmt"

	"willofdaedalus/mime/internal/engine/lexer"
	"willofdaedalus/mime/internal/engine/types"
)

type entityNode struct {
	name     string
	fields   []longField
	payload  entityObject
	response entityObject
}

type constraintInfo struct {
	kind  consType // bitfield
	value *string  // only for default/other values
}

type longField struct {
	name       string
	dt         dataType
	enums      []any
	consInfo   *constraintInfo
	fieldFlags fieldFlag
}

func (e entityNode) NodeLiteral() string {
	return "entity"
}

func handleEntity(p *Parser) (node, error) {
	if !expectTokOf(p.curToken, lexer.TokenEntity) {
		// p.pushError(fmt.Sprintf("expected entity token, got %s", p.curToken.Type))
		return nil, fmt.Errorf("expected entity token, got %s", p.curToken.Type)
	}
	p.advanceToken() // consume 'entity'

	if !expectTokOf(p.curToken, lexer.TokenIdent) {
		// p.pushError(fmt.Sprintf("expected entity name, got %s", p.curToken.Type))
		// fmt.Println("expected entity name")
		return nil, fmt.Errorf("expected entity name, got %s", p.curToken.Type)
	}

	entity := types.EntityNode{
		Name: p.curToken.Literal,
	}
	p.advanceToken() // consume entity name

	// check for arrow token
	if !expectTokOf(p.curToken, lexer.TokenArrow) {
		// p.pushError(fmt.Sprintf("expected -> after entity name, got %s", p.curToken.Type))
		// fmt.Println("expected -> after entity name")
		return nil, fmt.Errorf("expected -> after entity name, got %s", p.curToken.Type)
	}
	p.advanceToken() // consume '->'

	for p.curToken.Type != lexer.TokenEnd {
		if p.curToken.Type == lexer.TokenComment {
			p.advanceToken()
		}

		f, err := parseField(p)
		if err != nil {
			return nil, err
		}

		entity.Fields = append(entity.Fields, f)
	}

	return entity, nil
}
