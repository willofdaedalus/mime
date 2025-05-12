package parser

import "fmt"

// for each field in an entity, all attributes (usually after the data type)
// are collected into a struct. this struct then performs checks on each of
// the collected based on the data type of the field. if any attribute violates
// the rules of the data type i.e they're not compatitible with the data
// the struct immediately refuses

type (
	AttributeSet  int
	FieldDataType int
)

const (
	AttrDefault AttributeSet = 1 << iota
	AttrHash
	AttrUnique
	AttrRequired
	AttrIncrement
	AttrOverride
	AttrPrimary
	AttrHidden
	AttrReadonly
)

const (
	TypeText FieldDataType = iota
	TypeInt
	TypeFloat
	TypeUUID
	TypeTimestamp
	TypeBool
	TypeEnum
	TypeEntity // foreign key or nested
	TypeArray  // array of anything
)

var allowedAttrsByType = map[FieldDataType]AttributeSet{
	TypeText:      AttrDefault | AttrRequired | AttrUnique | AttrHash | AttrHidden | AttrReadonly,
	TypeInt:       AttrDefault | AttrRequired | AttrUnique | AttrIncrement | AttrHidden | AttrReadonly,
	TypeFloat:     AttrDefault | AttrRequired | AttrUnique | AttrHidden | AttrReadonly,
	TypeUUID:      AttrDefault | AttrRequired | AttrUnique | AttrHidden | AttrReadonly,
	TypeTimestamp: AttrDefault | AttrRequired | AttrHidden | AttrReadonly,
	TypeBool:      AttrDefault | AttrRequired | AttrHidden | AttrReadonly,
	TypeEnum:      AttrDefault | AttrRequired | AttrHidden | AttrReadonly,
	TypeEntity:    AttrDefault | AttrRequired | AttrOverride | AttrHidden | AttrReadonly,
	TypeArray:     AttrDefault | AttrRequired | AttrOverride | AttrHidden | AttrReadonly,
}

type FieldAttributes struct {
	Name       string
	DataType   FieldDataType
	Attributes AttributeSet
}

func (f FieldAttributes) Validate() error {
	allowed, ok := allowedAttrsByType[f.DataType]
	if !ok {
		return fmt.Errorf("unknown data type for field %s", f.Name)
	}
	if invalid := f.Attributes &^ allowed; invalid != 0 {
		return fmt.Errorf("field %s has invalid attributes: %b", f.Name, invalid)
	}
	return nil
}

func (f FieldAttributes) validateAttributeConflicts() error {
	attrs := f.Attributes
	dt := f.DataType

	if attrs&AttrIncrement != 0 && dt != TypeInt {
		return fmt.Errorf("increment only valid for int type")
	}
	if attrs&AttrHash != 0 && dt != TypeText {
		return fmt.Errorf("hash only valid for text type")
	}
	if attrs&(AttrIncrement|AttrReadonly) == (AttrIncrement | AttrReadonly) {
		return fmt.Errorf("field %s cannot be increment and readonly", f.Name)
	}
	if attrs&(AttrPrimary|AttrUnique) == (AttrPrimary | AttrUnique) {
		// optional warning — not fatal
		fmt.Printf("warning: field %s has both primary and unique (redundant)\n", f.Name)
	}
	if attrs&AttrOverride != 0 && dt != TypeEntity && dt != TypeArray {
		fmt.Printf("warning: field %s uses override but isn't a nested/array field\n", f.Name)
	}
	return nil
}
