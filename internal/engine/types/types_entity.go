package types

type EntityNode struct {
	Name     string
	Fields   []LongField
	name     string
	fields   []LongField
	payload  entityObject
	response entityObject
}

type LongField struct {
	Name       string
	Dt         dataType
	Attributes Attribute
}

type DataType struct {
	basic dataType
}

type entityField struct {
	name string
	dt   dataType
}

type EnumNode struct {
	Name    string
	Members []string
}

type constraintInfo struct {
	kind  consType // bitfield
	value *string  // only for default/other values
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
	dataOther
)

func (e EntityNode) NodeLiteral() string {
	return "entity"
}

func (e EnumNode) NodeLiteral() string {
	return "entity"
}
