package parser

import "willofdaedalus/mime/internal/engine/lexer"

type (
	dataType   int
	constraint int
)

const (
	DataText dataType = iota + 1
	DataNumber
	DataTimestamp
	DataUUID
)

const (
	ConsUnique constraint = iota + 1
	ConsIncrement
	ConsPrimary
	ConsRequired
	ConsFK
)

type entityNode struct {
	entityName string
	fields     []field
}

type field struct {
	name          string
	dt            dataType
	constraints   []constraint
	allowedValues []any
}

func (parser *Parser) parseEntity() *entityNode {
	return nil
}

// func (p *Parser) parseFields() []*field {
// 	if p.nextToken.Type != lexer.TokenIdent {
// 		p.pushError(fmt.Sprintf("expected an %s but got %s %s on line %d",
// 			lexer.TokenIdent.String(), p.curToken.Type.String(), p.curToken.Literal, p.curToken.LineNum))
//
// 		return nil
// 	}
//
// 	return nil
// }

// email text {unique, required}
func (p *Parser) parseField() {
	var fieldName, fieldType string
	var hasEnums bool
	var cons []string

	for p.curToken.Type != lexer.TokenEOF && p.curToken.Type != lexer.TokenNewline {
		if p.curToken.Type == lexer.TokenIdent {
			fieldName = p.curToken.Literal
			p.advanceToken()
		} else {
			p.pushError("expected a field name")
			p.advanceToken()
			return
		}

		// NOTE!
		// remember types can be like this;
		// manager ref employee {required}
		// where "ref employee" is the type of manager
		if lexer.IsValidDataType(p.curToken.Type) {
			fieldType = p.curToken.Literal
			p.advanceToken()
		} else {
			p.pushError("expected a data type")
			p.advanceToken()
			return
		}

		if p.curToken.Type == lexer.TokenOpenAngle {
			hasEnums = true
		}
	}
}
