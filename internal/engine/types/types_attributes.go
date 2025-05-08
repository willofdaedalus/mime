package types

type Attribute int

const (
	// runtime attributes don't affect the db migrations or schema
	AttributeHash Attribute = iota + 1
	AttributeSecret

	// db attributes
	AttributeIncrement
	AttributeUnique
	AttributePrimary
	AttributeDefault
	AttributeRequired
)
