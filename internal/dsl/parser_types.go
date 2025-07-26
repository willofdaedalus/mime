package dsl

type fieldType int

type Parser struct {
	lexer          *Lexer
	curToken       Token
	nextToken      Token
	errors         []error
	nodes          map[string]node
	invalidParsing bool
	// parserErrors   []parserError
}

const (
	typeNumber fieldType = iota + 1
	typeText
	typeBool
	typeFloat
)

type EntityNode struct {
	EntityName string
}

type entityField struct {
	name      string
	fieldType fieldType
}
