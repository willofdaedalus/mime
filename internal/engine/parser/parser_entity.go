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
	// if p.curToken.Type != lexer.TokenEntity {
	// 	p.pushError(fmt.Sprintf("expected entity token, got %s", p.curToken.Type))
	// 	return nil
	// }
	// p.advanceToken() // consume 'entity'
	//
	if p.curToken.Type != lexer.TokenIdent {
		p.pushError(fmt.Sprintf("expected entity name, got %s", p.curToken.Type))
		return nil
	}

	entity := &entityNode{
		entityName: p.curToken.Literal,
	}
	p.advanceToken() // consume entity name

	// check for arrow token
	if p.curToken.Type != lexer.TokenArrow {
		p.pushError(fmt.Sprintf("expected -> after entity name, got %s", p.curToken.Type))
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
		p.pushError(fmt.Sprintf("expected field name, got %s", p.curToken.Type))
		return nil
	}

	f := &field{
		name: p.curToken.Literal,
	}
	p.advanceToken() // consume field name

	// parse data type
	if !lexer.IsValidMemberOf(p.curToken.Type, lexer.AllDataTypes) {
		p.pushError(fmt.Sprintf("expected data type, got %s", p.curToken.Type))
		return nil
	}

	// map token type to data type
	switch p.curToken.Type {
	case lexer.TokenText:
		f.dt = dataText
	case lexer.TokenInt:
		f.dt = dataInt
	case lexer.TokenFloat:
		f.dt = dataReal
	case lexer.TokenTimestamp:
		f.dt = dataTimestamp
	case lexer.TokenUuid:
		f.dt = dataUUID
	default:
		p.pushError(fmt.Sprintf("unsupported data type: %s", p.curToken.Literal))
		return nil
	}
	p.advanceToken() // consume data type

	// check for optional enums
	if p.curToken.Type == lexer.TokenEnumOpen {
		enums := p.parseEnums(f.dt)
		if enums != nil {
			f.enums = enums
		}
	}

	// check for optional constraints
	if p.curToken.Type == lexer.TokenConsOpen {
		constraints := p.parseConstraints(f.dt)
		if constraints != nil {
			// no need to push any errors; I handled that in the parseConstraints func
			f.constraints = constraints
		}
	}

	return f
}

func (p *Parser) parseEnums(dataTypeInt dataType) []any {
	p.advanceToken() // consume '('

	var enums []any
	expectedType := lexer.TokenString // default

	// determine expected token type based on data type
	switch dataTypeInt {
	case dataText:
		expectedType = lexer.TokenString
	case dataInt:
	case dataReal:
		expectedType = lexer.TokenDigits
	}

	// make sure the data type aforehand is enumerable
	if _, ok := enumerableTypes[expectedType]; !ok {
		p.pushError(fmt.Sprintf("%s doesn't support enums", expectedType.String()))
		return nil
	}

	// parse enum values until closing parenthesis
	for p.curToken.Type != lexer.TokenEnumClose && p.curToken.Type != lexer.TokenEOF {
		if p.curToken.Type == expectedType {
			// convert token to appropriate value type based on data type
			var value any
			switch expectedType {
			case lexer.TokenString:
				value = p.curToken.Literal
			case lexer.TokenDigits:
				value, _ = strconv.Atoi(p.curToken.Literal)
			}
			enums = append(enums, value)
		} else if p.curToken.Type != lexer.TokenNewline {
			p.pushError(fmt.Sprintf("unexpected token in enum: %s", p.curToken.Type))
			return nil
		}

		p.advanceToken()
	}

	if p.curToken.Type == lexer.TokenEnumClose {
		p.advanceToken() // consume ')'
	} else {
		p.pushError("unclosed enum definition")
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
		var c constraint
		switch p.curToken.Type {
		case lexer.TokenConstraintUnique:
			c = consUnique
		case lexer.TokenConstraintAutoIncrement:
			c = consIncrement
		case lexer.TokenConstraintPrimaryKey:
			c = consPrimary
		case lexer.TokenConstraintNotNull:
			c = consRequired
		case lexer.TokenConstraintForeignKey:
			c = consFK
		default:
			p.pushError(fmt.Sprintf("unknown constraint: %s", p.curToken.Literal))
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
