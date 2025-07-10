package types

import "willofdaedalus/mime/internal/dsl"

type FieldKind int

const (
	FieldPrimitive FieldKind = iota // `name text`
	FieldReference                  // `owner @user.id`
	FieldEmbedded                   // `@person`
	FieldEnum                       // &enum_name
)

type EntityNode struct {
	Name   string
	Fields []*Field
}

type ReferenceTarget struct {
	Entity string
	Field  string
}

type Field struct {
	Name       string
	Kind       FieldKind
	DataType   DataType
	Target     *ReferenceTarget
	Attributes Attribute
}

type EnumNode struct {
	Name    string
	Members []string
}

type (
	DataType  int
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
	DataText DataType = iota + 1
	DataInt
	DataBool
	DataReal
	DataUUID
	DataEnum
	DataRef
	DataTimestamp
	DataOther
)

var TokenToDataType = map[dsl.TokenType]DataType{
	dsl.TokenTypeText:      DataText,
	dsl.TokenTypeInt:       DataInt,
	dsl.TokenTypeFloat:     DataReal,
	dsl.TokenTypeTimestamp: DataTimestamp,
	dsl.TokenTypeUuid:      DataUUID,
	dsl.TokenTypeBool:      DataBool,
}

func (e EntityNode) NodeLiteral() string {
	return "entity"
}

func (e EnumNode) NodeLiteral() string {
	return "entity"
}
