package dsl

type fieldType int

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
