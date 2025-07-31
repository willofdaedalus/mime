package dsl

import "fmt"

func handleEntity(p *Parser) (node, error) {
	var entity entityNode
	p.advanceToken() // skip entity

	if p.curToken.Type != TokenIdent {
		return nil, fmt.Errorf("expected entity name")
	}
	entity.name = p.curToken.Literal
	p.advanceToken() // skip name

	if p.curToken.Type != TokenArrow {
		return nil, fmt.Errorf("expected arrow for entity")
	}
	p.advanceToken() // skip arrow

	for p.curToken.Type != TokenEnd {
		if p.curToken.Type == TokenEOF {
			return nil, fmt.Errorf("unexpected end of file")
		}

		f, err := parseField(p)
		if err != nil {
			return nil, err
		}

		entity.fields = append(entity.fields, f)
	}

	return entity, nil
}

func parseField(p *Parser) (*field, error) {
	field := &field{}

	// field name
	if p.curToken.Type != TokenIdent {
		return nil, fmt.Errorf("unexpected token %s\n", p.curToken.Literal)
	}
	field.Name = p.curToken.Literal
	p.advanceToken()

	// data type
	dt, ok := TokenToDataType[p.curToken.Type]
	if !ok {
		return nil, fmt.Errorf("unknown data type %s\n", p.curToken.Literal)
	}
	field.DataType = dt
	p.advanceToken()

	if p.curToken.Type == TokenAttrOpen {
		p.advanceToken() // skip [
		parseAttributes(p)
	}

	return field, nil
}

func parseAttributes(p *Parser) {
	for p.curToken.Type != TokenAttrClose {
	}
}
