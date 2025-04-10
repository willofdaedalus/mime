package parser

import "willofdaedalus/mime/internal/engine/lexer"

type (
	dataType int
	consType int
)

const (
	dataText dataType = iota + 1
	dataInt
	dataReal
	dataTimestamp
	dataUUID
)

const (
	consUnique consType = iota + 1
	consIncrement
	consPrimary
	consRequired
	consDefault
	// consEnsure
	consFK
)

var enumerableTypes = map[lexer.TokenType]struct{}{
	lexer.TokenTypeText:  {},
	lexer.TokenTypeInt:   {},
	lexer.TokenTypeFloat: {},
}

var constrainableTypes = map[lexer.TokenType]struct{}{
	lexer.TokenTypeText:  {},
	lexer.TokenTypeInt:   {},
	lexer.TokenTypeFloat: {},
}

var tokenToDataType = map[lexer.TokenType]dataType{
	lexer.TokenTypeText:      dataText,
	lexer.TokenTypeInt:       dataInt,
	lexer.TokenTypeFloat:     dataReal,
	lexer.TokenTypeTimestamp: dataTimestamp,
	lexer.TokenTypeUuid:      dataUUID,
}

var tokenToConsType = map[lexer.TokenType]consType{
	lexer.TokenConstraintUnique:        consUnique,
	lexer.TokenConstraintAutoIncrement: consIncrement,
	lexer.TokenConstraintPrimaryKey:    consPrimary,
	lexer.TokenConstraintNotNull:       consRequired,
	lexer.TokenConstraintForeignKey:    consFK,
}

var consWithValues = map[consType]struct{}{
	consDefault: {},
}

var typeConstraintMap = map[dataType][]consType{
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

func (c consType) string() string {
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
