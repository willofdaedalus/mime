package parser

import (
	"fmt"

	l "willofdaedalus/mime/internal/engine/lexer"
	"willofdaedalus/mime/internal/engine/types"
)

func parseField(p *Parser) (*types.Field, error) {
	var field *types.Field

	// check for embed (@entity)
	if p.curToken.Type == l.TokenAtSymbol {
		p.advanceToken() // consume @

		if p.curToken.Type != l.TokenIdent {
			// return nil, fmt.Errorf("expected entity name after '@', got %s", p.curToken.Literal)
			p.addError(ParserLogError,
				fmt.Sprintf("expected entity name after '@', got %s", p.curToken.Literal))
		}
		field.Name = p.curToken.Literal

		p.advanceToken() // consume entity name

		// after @entity, the next token must be newline or comment
		if p.curToken.Type != l.TokenNewline && p.curToken.Type != l.TokenComment && p.curToken.Type != l.TokenEOF {
			p.addError(ParserLogError,
				fmt.Sprintf("expected entity name after '@', got %s", p.curToken.Literal))
			// return nil, fmt.Errorf("unexpected token after embedded entity: %s", p.curToken.Literal)
		}

		return &types.Field{
			Kind: types.FieldEmbedded,
		}, nil
	}

	// otherwise, treat as primitive field
	return parseFieldNormal(p)
}

func parseFieldNormal(p *Parser) (*types.Field, error) {
	var field *types.Field

	// field name
	if p.curToken.Type != l.TokenIdent {
		p.addError(ParserLogError,
			fmt.Sprintf("expected field name, got %s", p.curToken.Literal))
		// return nil, fmt.Errorf("expected field name, got %s", p.curToken.Literal)
	}
	field.Name = p.curToken.Literal
	p.advanceToken()

	// check for reference (@entity.field) or data type
	if p.curToken.Type == l.TokenAtSymbol {
		p.advanceToken() // skip the @ symbol
		refs, err := parseReferenceTarget(p)
		if err != nil {
			p.addError(ParserLogError, err.Error())
			// return nil, err // Return error instead of just adding to log
		}

		field.Target = refs
		field.Kind = types.FieldReference
		p.advanceToken() // consume the field name from reference
	} else if p.curToken.Type == l.TokenAmpersand {
		// handle enum reference (&enum_name)
		p.advanceToken() // consume &
		if p.curToken.Type != l.TokenIdent {
			p.addError(ParserLogError,
				fmt.Sprintf("expected enum name after '&', got %s", p.curToken.Literal))
			// return nil, fmt.Errorf("expected enum name after '&', got %s", p.curToken.Literal)
		}
		// for enum references, you might want to store this differently
		field.DataType = types.DataEnum
		p.advanceToken()
	} else {
		// regular data type
		dt, ok := types.TokenToDataType[p.curToken.Type]
		if !ok {
			p.addError(ParserLogError,
				fmt.Sprintf("unknown data type %s", p.curToken.Literal))
			// return nil, fmt.Errorf("unknown data type %s", p.curToken.Literal)
		}
		field.DataType = dt
		field.Kind = types.FieldPrimitive
		p.advanceToken()
	}

	// parse attributes if present
	if p.curToken.Type == l.TokenEnumOpen {
		attrs, err := parseAttributes(p)
		if err != nil {
			return nil, err
		}

		// convert slice to bitmask
		var attrMask types.Attribute
		for _, attr := range attrs {
			attrMask |= attr
		}
		field.Attributes = attrMask
	}

	// make sure line ends correctly
	if p.curToken.Type != l.TokenNewline && p.curToken.Type != l.TokenEOF && p.curToken.Type != l.TokenComment {
		p.addError(ParserLogError, fmt.Sprintf("unexpected token at end of field: %s", p.curToken.Literal))
		// return nil, fmt.Errorf("unexpected token at end of field: %s", p.curToken.Literal)
	}

	return field, nil
}

func parseAttributes(p *Parser) ([]types.Attribute, error) {
	var attributes []types.Attribute

	// consume opening bracket
	if p.curToken.Type != l.TokenEnumOpen {
		p.addError(ParserLogError, fmt.Sprintf("expected '[' to start attributes"))
		// return nil, fmt.Errorf("expected '[' to start attributes")
	}
	p.advanceToken()

	for p.curToken.Type != l.TokenEnumClose {
		if p.curToken.Type == l.TokenEOF {
			p.addError(ParserLogError, fmt.Sprintf("unexpected EOF while parsing attributes"))
			// return nil, fmt.Errorf("unexpected EOF while parsing attributes")
		}

		if p.curToken.Type == l.TokenIdent {
			attr, err := types.StringToAttribute(p.curToken.Literal)
			if err != nil {
				return nil, fmt.Errorf("unknown attribute: %s", p.curToken.Literal)
			}
			attributes = append(attributes, attr)
			p.advanceToken()
		}

		// // Handle comma separation or whitespace
		// if p.curToken.Type == l.TokenComma {
		// 	p.advanceToken()
		// } else if p.curToken.Type != l.TokenEnumClose {
		// 	// Skip whitespace or other separators if needed
		// 	p.advanceToken()
		// }
	}

	// Consume closing bracket
	if p.curToken.Type == l.TokenEnumClose {
		p.advanceToken()
	}

	return attributes, nil
}

func parseReferenceTarget(p *Parser) (*types.ReferenceTarget, error) {
	// example field that satisfies this
	// owner @user.id
	target := &types.ReferenceTarget{}

	if p.curToken.Type != l.TokenIdent {
		return nil, fmt.Errorf("expected entity name, got %s", p.curToken.Literal)
	}
	target.Entity = p.curToken.Literal
	p.advanceToken()

	if p.curToken.Type != l.TokenDot {
		return nil, fmt.Errorf("expected '.', got %s", p.curToken.Literal)
	}
	p.advanceToken()

	if p.curToken.Type != l.TokenIdent {
		return nil, fmt.Errorf("expected field name, got %s", p.curToken.Literal)
	}
	target.Field = p.curToken.Literal
	// don't advance here - let the caller handle it

	return target, nil
}
