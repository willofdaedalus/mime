package parser

import (
	"fmt"

	l "willofdaedalus/mime/internal/engine/lexer"
	"willofdaedalus/mime/internal/engine/types"
)

func parseField(p *Parser) (*types.Field, error) {
	// skip over leading whitespace/comments, if necessary
	p.skipNewlines()

	// check for embed (@entity)
	if p.curToken.Type == l.TokenAtSymbol {
		p.advanceToken() // consume @

		if p.curToken.Type != l.TokenIdent {
			return nil, fmt.Errorf("expected entity name after '@', got %s", p.curToken.Literal)
		}
		entityName := p.curToken.Literal

		p.advanceToken() // consume entity name

		// after @Entity, the next token must be newline or comment
		if p.curToken.Type != l.TokenNewline && p.curToken.Type != l.TokenComment && p.curToken.Type != l.TokenEOF {
			return nil, fmt.Errorf("unexpected token after embedded entity: %s", p.curToken.Literal)
		}

		return &types.Field{
			Name: entityName,
			Kind: types.FieldEmbedded,
		}, nil
	}

	// Otherwise, treat as primitive field
	return parseFieldNormal(p)
}

func parseFieldNormal(p *Parser) (*types.Field, error) {
	var field *types.Field

	// field name
	if p.curToken.Type != l.TokenIdent {
		return nil, fmt.Errorf("expected field name, got %s", p.curToken.Literal)
	}
	field.Name = p.curToken.Literal
	p.advanceToken()

	if p.curToken.Type == l.TokenAtSymbol {
		p.advanceToken() // skip the @ symbol
		refs, err := parseReferenceTarget(p)
		if err != nil {
			p.addError(ParserLogError, err.Error())
		} else {
			field.Target = refs
		}
	} else {
		dt, ok := types.TokenToDataType[p.curToken.Type]
		if !ok {
			p.addError(ParserLogError, fmt.Sprintf("unknown data type %s", p.curToken.Literal))
		} else {
			field.DataType = dt
		}
	}

	// attributes?
	if p.curToken.Type == l.TokenEnumOpen {
		attrs, err := parseAttributes(p)
		if err != nil {
			p.addError(ParserLogError, err.Error())
		} else {
			field.Attributes = attrs
		}
	}

	// make sure line ends correctly
	if p.curToken.Type != l.TokenNewline && p.curToken.Type != l.TokenEOF && p.curToken.Type != l.TokenComment {
		return nil, fmt.Errorf("unexpected token at end of field: %s", p.curToken.Literal)
	}

	return field, nil
}

func parseReferenceTarget(p *Parser) (*types.ReferenceTarget, error) {
	// example field that satisfies this
	// owner @player.id
	var target *types.ReferenceTarget
	if p.curToken.Type != l.TokenIdent {
		return nil, fmt.Errorf("unexpected symbol %s", p.curToken.Literal)
	}
	target.Entity = p.curToken.Literal
	p.advanceToken()

	if p.curToken.Type != l.TokenDot {
		return nil, fmt.Errorf("unexpected symbol %s", p.curToken.Literal)
	}
	p.advanceToken()

	if p.curToken.Type != l.TokenIdent {
		return nil, fmt.Errorf("unexpected symbol %s", p.curToken.Literal)
	}
	target.Field = p.curToken.Literal

	return target, nil
}
