package parser

import (
	"fmt"
	"strings"

	"willofdaedalus/mime/internal/engine/lexer"
)

type parserErrorLevel int

const (
	ParserLogError parserErrorLevel = iota
	ParserLogWarning
)

type parserError struct {
	errorLevel parserErrorLevel
	msg        string
}

type shortField struct {
	name *string
	dt   *dataType
}

// entityObject resolves the issue of payloads and responses
// by default it contains pointers to each field in the parent
// entity which the user can then override with their own
// default so long as the fields match those in the entity
type entityObject struct {
	isResponse bool
	fields     []shortField
}

type (
	dataType  int
	consType  uint8
	fieldFlag uint8
)

const (
	flagPayload  fieldFlag = 1 << 0
	flagResponse           = 1 << 1
	flagNullable           = 1 << 2
)

const consNone consType = 0
const (
	consUnique consType = 1 << iota
	consIncrement
	consPrimary
	consRequired
	consDefault
	consFK
	// consEnsure
)

const (
	dataText dataType = iota + 1
	dataInt
	dataBool
	dataReal
	dataUUID
	dataTimestamp
)

var payloadFriendly = map[dataType]struct{}{
	dataInt:  {},
	dataBool: {},
	dataReal: {},
	dataUUID: {},
	dataText: {},
}

var enumerableTypes = map[lexer.TokenType]struct{}{
	lexer.TokenTypeText:  {},
	lexer.TokenTypeInt:   {},
	lexer.TokenTypeFloat: {},
}

var constrainableTypes = map[lexer.TokenType]struct{}{
	lexer.TokenTypeInt:   {},
	lexer.TokenTypeText:  {},
	lexer.TokenTypeBool:  {},
	lexer.TokenTypeFloat: {},
}

var tokenToDataType = map[lexer.TokenType]dataType{
	lexer.TokenTypeText:      dataText,
	lexer.TokenTypeInt:       dataInt,
	lexer.TokenTypeFloat:     dataReal,
	lexer.TokenTypeTimestamp: dataTimestamp,
	lexer.TokenTypeUuid:      dataUUID,
	lexer.TokenTypeBool:      dataBool,
}

var tokenToConsType = map[lexer.TokenType]consType{
	lexer.TokenConstraintUnique:        consUnique,
	lexer.TokenConstraintAutoIncrement: consIncrement,
	lexer.TokenConstraintPrimaryKey:    consPrimary,
	lexer.TokenConstraintNotNull:       consRequired,
	lexer.TokenConstraintForeignKey:    consFK,
	lexer.TokenConstraintDefault:       consDefault,
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

func (d dataType) String() string {
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
	return fmt.Sprintf("%d not checked", d)
}

func (c consType) String() string {
	var parts []string
	if c&consUnique != 0 {
		parts = append(parts, "Unique")
	}
	if c&consIncrement != 0 {
		parts = append(parts, "Increment")
	}
	if c&consPrimary != 0 {
		parts = append(parts, "Primary")
	}
	if c&consRequired != 0 {
		parts = append(parts, "Required")
	}
	if c&consDefault != 0 {
		parts = append(parts, "Default")
	}
	if c&consFK != 0 {
		parts = append(parts, "FK")
	}
	if len(parts) == 0 {
		return "None"
	}
	return strings.Join(parts, "|")
}
