package parser

import (
	"fmt"
	"willofdaedalus/mime/internal/engine/types"
)

// for each field in an entity, all attributes (usually after the data type)
// are collected into a struct. this struct then performs checks on each of
// the collected based on the data type of the field. if any attribute violates
// the rules of the data type i.e they're not compatitible with the data
// the struct immediately refuses

type AttributeSet int

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

var allowedAttrsByType = map[types.DataType]AttributeSet{
	types.DataText:      AttrDefault | AttrRequired | AttrUnique | AttrHash | AttrHidden | AttrReadonly,
	types.DataInt:       AttrDefault | AttrRequired | AttrUnique | AttrIncrement | AttrHidden | AttrReadonly,
	types.DataReal:      AttrDefault | AttrRequired | AttrUnique | AttrHidden | AttrReadonly,
	types.DataUUID:      AttrDefault | AttrRequired | AttrUnique | AttrHidden | AttrReadonly,
	types.DataTimestamp: AttrDefault | AttrRequired | AttrHidden | AttrReadonly,
	types.DataBool:      AttrDefault | AttrRequired | AttrHidden | AttrReadonly,
	types.DataEnum:      AttrDefault | AttrRequired | AttrHidden | AttrReadonly,
	// types.DataEntity:    AttrDefault | AttrRequired | AttrOverride | AttrHidden | AttrReadonly,
	// types.DataArray:     AttrDefault | AttrRequired | AttrOverride | AttrHidden | AttrReadonly,
}

type FieldAttributes struct {
	Name       string
	DataType   types.DataType
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

	if attrs&AttrIncrement != 0 && dt != types.DataInt {
		return fmt.Errorf("increment only valid for int type")
	}
	if attrs&AttrHash != 0 && dt != types.DataText {
		return fmt.Errorf("hash only valid for text type")
	}
	if attrs&(AttrIncrement|AttrReadonly) == (AttrIncrement | AttrReadonly) {
		return fmt.Errorf("field %s cannot be increment and readonly", f.Name)
	}
	if attrs&(AttrPrimary|AttrUnique) == (AttrPrimary | AttrUnique) {
		// optional warning — not fatal
		fmt.Printf("warning: field %s has both primary and unique (redundant)\n", f.Name)
	}
	if attrs&AttrOverride != 0 && dt != types.DataEntity && dt != TypeArray {
		fmt.Printf("warning: field %s uses override but isn't a nested/array field\n", f.Name)
	}
	return nil
}
