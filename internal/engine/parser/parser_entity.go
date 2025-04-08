package parser

import (
	"fmt"
	"strconv"

	"willofdaedalus/mime/internal/engine/lexer"
)

type entityNode struct {
	entityName string
	fields     []field
}

type field struct {
	name        string
	dt          dataType
	constraints []constraint
	enums       []any
}

func (p *Parser) parseEntity() *entityNode {
	if p.curToken.Type != lexer.TokenEntity {
		p.pushError(fmt.Sprintf("expected entity token, got %s", p.curToken.Type))
		return nil
	}
	p.advanceToken() // consume 'entity'

	if p.curToken.Type != lexer.TokenIdent {
		p.pushError(fmt.Sprintf("expected entity name, got %s", p.curToken.Type))
		fmt.Println("expected entity name")
		return nil
	}

	entity := &entityNode{
		entityName: p.curToken.Literal,
	}
	p.advanceToken() // consume entity name

	// check for arrow token
	if p.curToken.Type != lexer.TokenArrow {
		p.pushError(fmt.Sprintf("expected -> after entity name, got %s", p.curToken.Type))
		fmt.Println("expected -> after entity name")
		return nil
	}
	p.advanceToken() // consume '->'

	// parse fields until 'end' token
	for p.curToken.Type != lexer.TokenEnd && p.curToken.Type != lexer.TokenEOF {
		if p.curToken.Type == lexer.TokenNewline {
			p.advanceToken() // skip newlines
			continue
		}

		field := p.parseField()
		if field != nil {
			entity.fields = append(entity.fields, *field)
		}
	}

	if p.curToken.Type == lexer.TokenEnd {
		p.advanceToken() // consume 'end'
	}

	return entity
}

// example field
// student_id number {unique}
func (p *Parser) parseField() *field {
	// expect field name (identifier)
	if p.curToken.Type != lexer.TokenIdent {
		p.pushError(fmt.Sprintf("%s:%d; expected field name, got %s",
			p.curToken.FileName, p.curToken.LineNum, p.curToken.Type))
		return nil
	}

	f := &field{
		name: p.curToken.Literal,
	}
	p.advanceToken() // consume field name

	// parse data type
	if !lexer.IsValidMemberOf(p.curToken.Type, lexer.AllDataTypes) {
		p.pushError(fmt.Sprintf("%s:%d; expected data type, got %s",
			p.curToken.FileName, p.curToken.LineNum, p.curToken.Type))
		return nil
	}

	// map token type to data type
	t, ok := tokenToDataType[p.curToken.Type]
	if !ok {
		p.pushError(fmt.Sprintf("%s:%d; unsupported data type: %s",
			p.curToken.FileName, p.curToken.LineNum, p.curToken.Literal))
		return nil
	}
	f.dt = t
	fieldDataType := p.curToken.Type
	p.advanceToken() // consume data type

	// continue parsing annotations until newline or unexpected token
	for {
		switch p.curToken.Type {
		case lexer.TokenEnumOpen:
			enums := p.parseEnums(fieldDataType)
			if enums != nil {
				f.enums = enums
			}
		case lexer.TokenConsOpen:
			constraints := p.parseConstraints(f.dt)
			if constraints != nil {
				f.constraints = constraints
			}
		case lexer.TokenListOpen:
			if p.nextToken.Type != lexer.TokenListClose {
				p.pushError(fmt.Sprintf("%s:%d; expected ], got %s",
					p.curToken.FileName, p.curToken.LineNum, p.nextToken.Literal))
				return nil
			}
		case lexer.TokenNewline:
			p.advanceToken() // consume newline, done with field
			return f
		default:
			if _, ok := lexer.AnnotationOpens[p.curToken.Type]; !ok {
				// invalid token found where annotation was expected
				p.pushError(fmt.Sprintf("%s:%d; unexpected token %s after data type",
					p.curToken.FileName, p.curToken.LineNum, p.curToken.Type))
			}
			return f
		}
	}
}

func (p *Parser) parseEnums(fdt lexer.TokenType) []any {
	var enums []any
	p.advanceToken() // consume '('

	// make sure the data type aforehand is enumerable
	if _, ok := enumerableTypes[fdt]; !ok {
		p.pushError(fmt.Sprintf("%s:%d; %s doesn't support enums",
			p.curToken.FileName, p.curToken.LineNum, fdt.String()))
		return nil
	}

	// keeping this as only 3 types can be enumerated... for now
	expectedType := lexer.TokenString
	switch fdt {
	case lexer.TokenInt:
		expectedType = lexer.TokenDigits
	case lexer.TokenFloat:
		expectedType = lexer.TokenDigitsFloat
	}

	// parse enum values until closing parenthesis
	for p.curToken.Type != lexer.TokenEnumClose && p.curToken.Type != lexer.TokenEOF {
		// make sure the user is not adding unrelated data types
		if p.curToken.Type == expectedType {
			// convert token to appropriate value type based on data type
			var value any
			switch expectedType {
			case lexer.TokenString:
				value = p.curToken.Literal
			case lexer.TokenDigits:
				value, _ = strconv.Atoi(p.curToken.Literal)
			case lexer.TokenDigitsFloat:
				value, _ = strconv.ParseFloat(p.curToken.Literal, 64)
			}
			enums = append(enums, value)
		} else if p.curToken.Type != lexer.TokenNewline {
			p.pushError(fmt.Sprintf("%s:%d; unexpected token in enum: %s",
				p.curToken.FileName, p.curToken.LineNum, p.curToken.Type))
			return nil
		}

		p.advanceToken()
	}

	if p.curToken.Type == lexer.TokenEnumClose {
		p.advanceToken() // consume ')'
	} else {
		p.pushError(fmt.Sprintf("%s:%d; unclosed enum definition",
			p.curToken.FileName, p.curToken.LineNum))
		return nil
	}

	return enums
}

func (p *Parser) parseConstraints(fieldType dataType) []constraint {
	p.advanceToken() // consume '{'

	var constraints []constraint

	// parse constraints until closing brace
	for p.curToken.Type != lexer.TokenConsClose && p.curToken.Type != lexer.TokenEOF {
		if p.curToken.Type == lexer.TokenNewline {
			p.advanceToken() // skip newlines
			continue
		}

		// map constraint tokens to constraint types
		c, ok := tokenConstraintToConstraintType[p.curToken.Type]
		if !ok {
			p.pushError(fmt.Sprintf("%s:%d; unknown constraint %q",
				p.curToken.FileName, p.curToken.LineNum, p.curToken.Literal))
			p.advanceToken()
			continue
		}

		constraints = append(constraints, c)
		p.advanceToken() // consume constraint
	}

	if p.curToken.Type == lexer.TokenConsClose {
		p.advanceToken() // consume '}'
	} else {
		p.pushError("unclosed constraint definition")
		return nil
	}

	if verifyConstraints(p, fieldType, constraints) {
		// NOTE!
		// make checks to ensure super constraints have precedence
		// over other constraints;
		// eg primary is basically unique, autoincrement and not_null
		return constraints
	}

	return nil
}

func verifyConstraints(p *Parser, fdataType dataType, cons []constraint) bool {
	hasInvalid := false
	if len(cons) > 0 {
		for i := range cons {
			if _, ok := typeConstraintMap[fdataType]; !ok {
				hasInvalid = true
				p.pushError(fmt.Sprintf("%s doesn't support constraint %s",
					fdataType.string(), cons[i].string()))
			}
		}
	}

	return !hasInvalid
}
