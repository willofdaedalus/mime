package parser

import (
	"fmt"

	l "willofdaedalus/mime/internal/engine/lexer"
	"willofdaedalus/mime/internal/engine/types"
)

func parseField(p *Parser) (*types.Field, error) {
	// handle embedded entity: @entity
	if p.curType() == l.TokenAtSymbol {
		return parseEmbeddedField(p)
	}

	return parseRegularField(p)
}

func parseEmbeddedField(p *Parser) (*types.Field, error) {
	field := &types.Field{}
	p.advanceToken() // consume the @

	// get entity name
	if p.curToken.Type != l.TokenIdent {
		return nil, fmt.Errorf("expected entity name after '@', got %s", p.curToken.Literal)
	}

	field.Name = p.curToken.Literal
	field.Kind = types.FieldEmbedded
	p.advanceToken()

	return field, nil
}

func parseRegularField(p *Parser) (*types.Field, error) {
	field := &types.Field{}

	// get field name
	if p.curToken.Type != l.TokenIdent {
		return nil, fmt.Errorf("expected field name, got %s", p.curToken.Literal)
	}
	field.Name = p.curToken.Literal
	p.advanceToken()

	// Parse the target (could be @entity.field, &enum, or primitive type)
	return parseFieldTarget(p, field)
}

func parseFieldTarget(p *Parser, field *types.Field) (*types.Field, error) {
	switch p.curToken.Type {
	case l.TokenAtSymbol:
		// reference: @entity.field
		return parseReference(p, field)

	case l.TokenAmpersand:
		// enum: &enum_name
		return parseEnumReference(p, field)

	default:
		// Primitive type
		return parsePrimitiveType(p, field)
	}
}

func parseReference(p *Parser, field *types.Field) (*types.Field, error) {
	// consume '@'
	p.advanceToken()

	// Get entity name
	if p.curToken.Type != l.TokenIdent {
		return nil, fmt.Errorf("expected entity name after '@', got %s", p.curToken.Literal)
	}
	entityName := p.curToken.Literal
	p.advanceToken()

	// Check for dot (for entity.field)
	var fieldName string
	if p.curToken.Type == l.TokenDot {
		p.advanceToken() // consume '.'
		if p.curToken.Type != l.TokenIdent {
			return nil, fmt.Errorf("expected field name after '.', got %s", p.curToken.Literal)
		}
		fieldName = p.curToken.Literal
		p.advanceToken()
	}

	field.Kind = types.FieldReference
	field.Target = &types.ReferenceTarget{
		Entity: entityName,
		Field:  fieldName, // empty if just @entity
	}

	return field, nil
}

func parseEnumReference(p *Parser, field *types.Field) (*types.Field, error) {
	// Consume '&'
	p.advanceToken()

	if p.curToken.Type != l.TokenIdent {
		return nil, fmt.Errorf("expected enum name after '&', got %s", p.curToken.Literal)
	}

	field.Kind = types.FieldEnum
	field.DataType = types.DataEnum
	field.Target = &types.ReferenceTarget{
		Entity: p.curToken.Literal,
	}
	p.advanceToken()

	// if there are constraints
	// not sure if I'm adding these or not for enums but we'll see
	if p.curType() == l.TokenEnumOpen {
	}

	return field, nil
}

func parsePrimitiveType(p *Parser, field *types.Field) (*types.Field, error) {
	dt, ok := types.TokenToDataType[p.curToken.Type]
	if !ok {
		return nil, fmt.Errorf("unknown data type %s", p.curToken.Literal)
	}

	field.DataType = dt
	field.Kind = types.FieldPrimitive
	p.advanceToken()

	// if there are constraints
	if p.curType() == l.TokenEnumOpen {
		attr, err := parseAttributes(p)
		if err != nil {
			return nil, err
		}

		field.Attributes = attr
	}

	return field, nil
}

func parseAttributes(p *Parser) (types.Attribute, error) {
	var attribute types.Attribute

	p.advanceToken()

	for p.curType() != l.TokenEnumClose {
		if p.curType() == l.TokenEOF {
			// p.addError(ParserLogError, fmt.Sprintf("unexpected EOF while parsing attributes"))
			return 0, fmt.Errorf("unexpected EOF while parsing attributes")
		}

		attr, ok := types.LiteralToAttr[p.curToken.Literal]
		if !ok {
			return 0, fmt.Errorf("unknown attribute: %s", p.curToken.Literal)
		}
		attribute |= attr
		p.advanceToken()
	}
	p.advanceToken() // consume the closing ']'

	return attribute, nil
}
