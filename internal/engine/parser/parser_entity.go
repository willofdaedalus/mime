package parser

import (
	"fmt"
	"strconv"

	"willofdaedalus/mime/internal/engine/lexer"
)

type entityNode struct {
	name   string
	fields []field
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
		name: p.curToken.Literal,
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
		if p.curToken.Type == lexer.TokenNewline || p.curToken.Type == lexer.TokenComment {
			p.advanceToken() // skip newlines and comments
			continue
		}

		field := p.parseField()
		if field == nil {
			return nil
		}
		entity.fields = append(entity.fields, *field)
	}

	if p.curToken.Type == lexer.TokenEnd {
		p.advanceToken() // consume 'end'
	} else {
		p.pushError(fmt.Sprintf("%s:%d; expected end keyword at end of entity definition",
			p.curToken.FileName, p.curToken.LineNum))
		return nil
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
	for p.curToken.Type != lexer.TokenNewline {
		switch p.curToken.Type {
		case lexer.TokenEnumOpen:
			enums := p.parseEnums(fieldDataType)
			// enums being nil means we got an error
			if enums == nil {
				// if one thing is null the whole entity is dead
				return nil
			}
			f.enums = enums
		case lexer.TokenConsOpen:
			// returning nil because the user specified a constraint and didn't finish
			constraints := p.parseConstraints(fieldDataType)
			if constraints == nil {
				return nil
			}
			f.constraints = constraints
		case lexer.TokenListOpen:
			if p.nextToken.Type != lexer.TokenListClose {
				p.pushError(fmt.Sprintf("%s:%d; expected ], got %s",
					p.curToken.FileName, p.curToken.LineNum, p.nextToken.Literal))
				return nil
			}
		default:
			if _, ok := lexer.AnnotationOpens[p.curToken.Type]; !ok {
				// invalid token found where annotation was expected
				p.pushError(fmt.Sprintf("%s:%d; unexpected token %s after data type",
					p.curToken.FileName, p.curToken.LineNum, p.curToken.Type))
			}
			return f
		}
	}

	return f
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
			p.advanceToken()
			continue
		}

		// on mismatched data type return nil immediately; don't waste
		// resources processing what we won't return
		errMsg := fmt.Sprintf("%s:%d; unexpected type in enum: %s",
			p.curToken.FileName, p.curToken.LineNum, p.curToken.Type.String())
		if p.curToken.Type == lexer.TokenNewline {
			errMsg = fmt.Sprintf("%s:%d; unclosed enum definition: expected )",
				p.curToken.FileName, p.curToken.LineNum)
		}
		p.pushError(errMsg)
		return nil
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

// for now we only support constraints without values so something like
// default wouldn't work
func (p *Parser) parseConstraints(fdt lexer.TokenType) []constraint {
	p.advanceToken() // consume '{'

	var constraints []constraint

	// parse constraints until closing brace
	// NOTE; we should probably make sure to check against newline token
	for p.curToken.Type != lexer.TokenConsClose && p.curToken.Type != lexer.TokenEOF {
		// any newline while parsing some constraints or enums or list is considered
		// an error and should be treated as such
		if p.curToken.Type == lexer.TokenNewline {
			// p.advanceToken() // skip newlines
			// continue
			p.pushError("unclosed constraint definition")
			return nil
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
	}
	// else {
	// 	p.pushError("unclosed constraint definition")
	// 	return nil
	// }

	// if verifyConstraints(p, fdt, constraints) {
	// 	// NOTE!
	// 	// make checks to ensure super constraints have precedence
	// 	// over other constraints;
	// 	// eg primary is basically unique, autoincrement and not_null
	// 	return constraints
	// }

	return constraints
}

// func verifyConstraints(p *Parser, fdt lexer.TokenType, cons []constraint) bool {
// 	hasInvalid := false
// 	if len(cons) > 0 {
// 		for i := range cons {
// 			if _, ok := typeConstraintMap[fdt]; !ok {
// 				hasInvalid = true
// 				p.pushError(fmt.Sprintf("%s doesn't support constraint %s",
// 					fdt.string(), cons[i].string()))
// 			}
// 		}
// 	}
//
// 	return !hasInvalid
// }
