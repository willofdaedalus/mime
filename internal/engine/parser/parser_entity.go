package parser

import (
	"fmt"

	l "willofdaedalus/mime/internal/engine/lexer"
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
	if p.curType() != l.TokenEntity {
		return nil, fmt.Errorf("expected entity token, got %s", p.curToken.Type)
	}
	p.advanceToken() // consume 'entity'

	if p.curType() != l.TokenIdent {
		return nil, fmt.Errorf("expected entity name, got %s", p.curToken.Type)
	}

	entity := types.EntityNode{
		Name: p.curToken.Literal,
	}
	p.advanceToken() // consume entity name

	// check for arrow token
	if p.curType() != l.TokenArrow {
		return nil, fmt.Errorf("expected -> after entity name, got %s", p.curToken.Type)
	}
	p.advanceToken() // consume '->'

	for p.curType() != l.TokenEnd {
		if p.curToken.Type == l.TokenComment {
			p.advanceToken()
			continue
		}

		// if p.curToken.Type == l.TokenEnd {
		// 	break
		// }

		f, err := parseField(p)
		if err != nil {
			return nil, err
		}

		entity.Fields = append(entity.Fields, f)
	}

	return &entity, nil
}
