package parser

import "willofdaedalus/mime/internal/engine/lexer"

type (
	dataType   int
	constraint int
)

const (
	dataText dataType = iota + 1
	dataInt
	dataReal
	dataTimestamp
	dataUUID
)

const (
	consUnique constraint = iota + 1
	consIncrement
	consPrimary
	consRequired
	consFK
)

var enumerableTypes = map[lexer.TokenType]struct{}{
	lexer.TokenText:  {},
	lexer.TokenInt:   {},
	lexer.TokenFloat: {},
}

var tokenToDataType = map[lexer.TokenType]dataType{
	lexer.TokenText:      dataText,
	lexer.TokenInt:       dataInt,
	lexer.TokenFloat:     dataReal,
	lexer.TokenTimestamp: dataTimestamp,
	lexer.TokenUuid:      dataUUID,
}

var tokenConstraintToConstraintType = map[lexer.TokenType]constraint{
	lexer.TokenConstraintUnique:        consUnique,
	lexer.TokenConstraintAutoIncrement: consIncrement,
	lexer.TokenConstraintPrimaryKey:    consPrimary,
	lexer.TokenConstraintNotNull:       consRequired,
	lexer.TokenConstraintForeignKey:    consFK,
}

var typeConstraintMap = map[dataType][]constraint{
	dataText: {
		consUnique, consPrimary, consFK,
	},
	dataInt: {
		consUnique, consIncrement, consPrimary, consFK,
	},
	dataReal: {
		consUnique, consIncrement, consRequired,
	},
	dataUUID: {
		consRequired,
	},
}

func (d dataType) string() string {
	switch d {
	case dataText:
		return "text"
	case dataInt:
		return "int"
	case dataReal:
		return "real"
	case dataTimestamp:
		return "timestamp"
	case dataUUID:
		return "uuid"
	}
	return "not checked!"
}

func (c constraint) string() string {
	switch c {
	case consUnique:
		return "unique"
	case consIncrement:
		return "increment"
	case consPrimary:
		return "primary"
	case consRequired:
		return "required"
	case consFK:
		return "fk"
	}
	return "not checked!"
}
