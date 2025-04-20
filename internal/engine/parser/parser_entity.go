package parser

import (
	"fmt"
	"strconv"

	"willofdaedalus/mime/internal/engine/lexer"
)

type entityNode struct {
	name     string
	fields   []longField
	payload  entityObject
	response entityObject
}

// entityObject resolves the issue of payloads and responses
// by default it contains pointers to each field in the parent
// entity which the user can then override with their own
// default so long as the fields match those in the entity
type entityObject struct {
	fields []shortField
}

type shortField struct {
	name *string
	dt   *dataType
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

func (p *Parser) parseEntity() *entityNode {
	if !expectTokOf(p.curToken, lexer.TokenEntity) {
		p.pushError(fmt.Sprintf("expected entity token, got %s", p.curToken.Type))
		return nil
	}
	p.advanceToken() // consume 'entity'

	if !expectTokOf(p.curToken, lexer.TokenIdent) {
		p.pushError(fmt.Sprintf("expected entity name, got %s", p.curToken.Type))
		fmt.Println("expected entity name")
		return nil
	}

	entity := &entityNode{
		name: p.curToken.Literal,
	}
	p.advanceToken() // consume entity name

	// check for arrow token
	if !expectTokOf(p.curToken, lexer.TokenArrow) {
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

	if err := entity.cleanupEntity(); err != nil {
		for _, e := range entity.cleanupEntity() {
			p.pushError(fmt.Sprintf("entity %q: %s", entity.name, e.Error()))
		}
		return nil
	}

	entity.makePayload()
	entity.makeResponse()

	return entity
}

// example field
// student_id number {unique}
func (p *Parser) parseField() *longField {
	// expect field name (identifier)
	if p.curToken.Type != lexer.TokenIdent {
		p.pushError(fmt.Sprintf("%s:%d; expected field name, got %s",
			p.curToken.FileName, p.curToken.LineNum, p.curToken.Type))
		return nil
	}

	f := &longField{
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
			cons := p.parseConstraints(fieldDataType)
			if cons == nil {
				// returning nil because the user specified a constraint and didn't finish
				return nil
			}
			f.consInfo = cons
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

	// assign the payload friendly flag to the field based on the following
	// rule; it's payload friendly (see list/map) and doesn't have any offending
	// constraints such as increment
	if _, ok := payloadFriendly[f.dt]; ok {
		if f.consInfo.kind&consIncrement == 0 {
			f.fieldFlags |= flagPayload
		}
	}

	// every field is automatically included in the response unless
	// the user overrides them with alter ref <entity>.response
	f.fieldFlags |= flagResponse

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
	case lexer.TokenTypeInt:
		expectedType = lexer.TokenDigits
	case lexer.TokenTypeFloat:
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

func (p *Parser) parseConstraints(fdt lexer.TokenType) *constraintInfo {
	result := &constraintInfo{
		kind:  consNone,
		value: nil,
	}
	// var constraints []constraint

	// Make sure the field's data type can be constrained
	if _, ok := constrainableTypes[fdt]; !ok {
		p.pushError(fmt.Sprintf("%s:%d; data type %q doesn't support constraints",
			p.curToken.FileName, p.curToken.LineNum, fdt.String()))
		return nil
	}

	p.advanceToken() // consume '{'

	// Parse constraints until closing brace or EOF
	for p.curToken.Type != lexer.TokenConsClose && p.curToken.Type != lexer.TokenEOF {
		if p.curToken.Type == lexer.TokenNewline {
			p.pushError(fmt.Sprintf("%s:%d; unexpected newline in constraint definition",
				p.curToken.FileName, p.curToken.LineNum))
			return nil
		}

		// Map constraint tokens to constraint types
		c, ok := tokenToConsType[p.curToken.Type]
		if !ok {
			p.pushError(fmt.Sprintf("%s:%d; unknown constraint %q",
				p.curToken.FileName, p.curToken.LineNum, p.curToken.Literal))
			p.advanceToken()
			return nil
		}

		// we found a duplicate constraint; fail fast
		if result.kind&c != 0 {
			p.pushError(fmt.Sprintf("%s:%d; duplicate constraint %s",
				p.curToken.FileName, p.curToken.LineNum, p.curToken.Literal))
			return nil
		}

		// assign the valid constype to the result if not a duplicate
		result.kind |= c
		p.advanceToken() // consume constraint token

		// check if constraint requires a value (e.g., default value)
		if p.curToken.Type == lexer.TokenColon {
			if _, ok := consWithValues[c]; !ok {
				p.pushError(fmt.Sprintf("%s:%d; constraint %q doesn't support values",
					p.curToken.FileName, p.curToken.LineNum, p.curToken.Literal))
				return nil
			}
			p.advanceToken() // consume the colon

			if p.curToken.Type != lexer.TokenString {
				p.pushError(fmt.Sprintf("%s:%d; expected a string for constraint value",
					p.curToken.FileName, p.curToken.LineNum))
				return nil
			}

			// verify that the value matches the field's data type
			if !verifyConstraintValue(fdt, p.curToken.Literal) {
				p.pushError(fmt.Sprintf("%s:%d; data type for constraint value doesn't match %q",
					p.curToken.FileName, p.curToken.LineNum, fdt.String()))
				return nil
			}

			// set the value for the constraint
			result.value = &p.curToken.Literal
			p.advanceToken() // consume value token
		}
	}

	// Consume the closing brace
	if p.curToken.Type == lexer.TokenConsClose {
		p.advanceToken() // consume '}'
	}

	return result
}

func verifyConstraintValue(fdt lexer.TokenType, v string) bool {
	var err error

	switch fdt {
	case lexer.TokenTypeInt:
		_, err = strconv.Atoi(v)
	case lexer.TokenTypeFloat:
		_, err = strconv.ParseFloat(v, 64)
	case lexer.TokenTypeText:
		// text by default is whatever the default value is
		return true
	}

	return err == nil
}

func (e *entityNode) cleanupEntity() []error {
	errs := make([]error, 0)
	validFields := make(map[string]struct{}, 0)

	for _, f := range e.fields {
		fieldName := f.name

		if _, ok := validFields[fieldName]; ok {
			errs = append(errs, fmt.Errorf("duplicate field name %q", fieldName))
		} else if _, ok := lexer.Keywords[fieldName]; ok {
			errs = append(errs, fmt.Errorf("field name %q is a reserved keyword", fieldName))
		} else {
			validFields[fieldName] = struct{}{}
		}

		// we don't need to verify the constraints only enums
		if len(f.enums) > 0 {
			errs = append(errs, verifyEnums(f.enums)...)
		}
	}

	if len(errs) > 0 {
		return errs
	}

	return nil
}

func verifyEnums(enums []any) []error {
	errs := make([]error, 0)
	seenEnums := make(map[any]struct{}, 0)

	for _, enum := range enums {
		if _, ok := seenEnums[enum]; ok {
			errs = append(errs, fmt.Errorf("duplicate enum value %v", enum))
			continue
		}
		seenEnums[enum] = struct{}{}
	}

	return errs
}

// helper for both payload & response
func (e *entityNode) buildObject(flag fieldFlag) entityObject {
	var obj entityObject
	for _, f := range e.fields {
		if f.fieldFlags&flag == 0 {
			continue
		}
		obj.fields = append(obj.fields, shortField{
			name: &f.name,
			dt:   &f.dt,
		})
	}
	return obj
}

func (e *entityNode) makePayload() {
	if obj := e.buildObject(flagPayload); len(obj.fields) > 0 {
		e.payload = obj
	}
}

func (e *entityNode) makeResponse() {
	if obj := e.buildObject(flagResponse); len(obj.fields) > 0 {
		e.response = obj
	}
}
