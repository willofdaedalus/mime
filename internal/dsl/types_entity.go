package dsl

type fieldKind int

const (
	FieldPrimitive fieldKind = iota // `name text`
	FieldReference                  // `owner @user.id`
	FieldEmbedded                   // `@person`
	FieldEnum                       // &enum_name
)

type entityNode struct {
	name   string
	fields []*field
}

type referenceTarget struct {
	Entity string
	Field  string
}

type field struct {
	Name       string
	DataType   dataType
	Attributes attribute
	// Kind       fieldKind
	// Target     *referenceTarget
}

type EnumNode struct {
	Name    string
	Members []string
}

type (
	dataType  int
	fieldFlag uint8
)

const (
	flagPayload  fieldFlag = 1 << 0
	flagResponse           = 1 << 1
	flagNullable           = 1 << 2
)

const (
	DataText dataType = iota + 1
	DataInt
	DataBool
	DataReal
	DataUUID
	DataEnum
	DataTimestamp
	DataOther
)

var TokenToDataType = map[TokenType]dataType{
	TokenTypeText:      DataText,
	TokenTypeInt:       DataInt,
	TokenTypeFloat:     DataReal,
	TokenTypeTimestamp: DataTimestamp,
	TokenTypeUuid:      DataUUID,
	TokenTypeBool:      DataBool,
}

var tokenToAttr = map[TokenType]attribute{
	TokenAttrIncrement: AttrIncrement,
	TokenAttrUnique:    AttrUnique,
	TokenAttrDefault:   AttrDefault,
	TokenAttrHidden:    AttrHidden,
	TokenAttrCheck:     AttrCheck,
	TokenAttrHash:      AttrHash,
}

func (e entityNode) NodeLiteral() string {
	return "entity"
}

func (e EnumNode) NodeLiteral() string {
	return "entity"
}
