package types

import (
	"fmt"
	"strings"
)

// for each field in an entity, all attributes (usually after the data type)
// are collected into a struct. this struct then performs checks on each of
// the collected based on the data type of the field. if any attribute violates
// the rules of the data type i.e they're not compatitible with the data
// the struct immediately refuses

type Attribute int

const (
	AttrDefault Attribute = 1 << iota
	AttrHash
	AttrUnique
	AttrRequired
	AttrIncrement
	AttrOverride
	AttrPrimary
	AttrHidden
	AttrReadonly
)

var allowedAttrsByType = map[DataType]Attribute{
	DataText:      AttrDefault | AttrRequired | AttrUnique | AttrHash | AttrHidden | AttrReadonly,
	DataInt:       AttrDefault | AttrRequired | AttrUnique | AttrIncrement | AttrHidden | AttrReadonly | AttrPrimary,
	DataReal:      AttrDefault | AttrRequired | AttrUnique | AttrHidden | AttrReadonly,
	DataUUID:      AttrDefault | AttrRequired | AttrUnique | AttrHidden | AttrReadonly | AttrPrimary,
	DataTimestamp: AttrDefault | AttrRequired | AttrHidden | AttrReadonly,
	DataBool:      AttrDefault | AttrRequired | AttrHidden | AttrReadonly,
	DataEnum:      AttrDefault | AttrRequired | AttrHidden | AttrReadonly,
}

// helper function to convert string to attribute
func StringToAttribute(s string) (Attribute, error) {
	switch strings.ToLower(s) {
	case "default":
		return AttrDefault, nil
	case "hash":
		return AttrHash, nil
	case "unique":
		return AttrUnique, nil
	case "required":
		return AttrRequired, nil
	case "increment", "auto_increment":
		return AttrIncrement, nil
	case "override":
		return AttrOverride, nil
	case "primary":
		return AttrPrimary, nil
	case "hidden":
		return AttrHidden, nil
	case "readonly":
		return AttrReadonly, nil
	default:
		return 0, fmt.Errorf("unknown attribute: %s", s)
	}
}

// Validation helpers

// ValidateFieldAttributes checks if the given attributes are valid for the field's data type
func ValidateFieldAttributes(field *Field) error {
	if field.Kind == FieldEmbedded || field.Kind == FieldReference {
		// For embedded and reference fields, only certain attributes make sense
		allowedForRefs := AttrRequired | AttrHidden | AttrReadonly | AttrOverride
		if field.Attributes & ^allowedForRefs != 0 {
			return fmt.Errorf("invalid attributes for %s field '%s'",
				fieldKindToString(field.Kind), field.Name)
		}
		return nil
	}

	// for primitive fields, check against the allowed attributes map
	allowed, exists := allowedAttrsByType[field.DataType]
	if !exists {
		return fmt.Errorf("no attribute validation rules for data type %v", field.DataType)
	}

	// check if any disallowed attributes are set
	if field.Attributes & ^allowed != 0 {
		return fmt.Errorf("invalid attributes for %s field '%s' with type %s",
			fieldKindToString(field.Kind), field.Name, dataTypeToString(field.DataType))
	}

	// additional validation rules
	if err := validateAttributeCombinations(field); err != nil {
		return err
	}

	return nil
}

// ValidateAttributeCombinations checks for conflicting attribute combinations
func validateAttributeCombinations(field *Field) error {
	attrs := field.Attributes

	// primary key implies unique and required
	if attrs&AttrPrimary != 0 {
		if attrs&AttrUnique == 0 {
			return fmt.Errorf("primary key field '%s' must also be unique", field.Name)
		}
		if attrs&AttrRequired == 0 {
			return fmt.Errorf("primary key field '%s' must also be required", field.Name)
		}
	}

	// auto-increment typically implies unique and required (for int fields)
	if attrs&AttrIncrement != 0 {
		if field.DataType != DataInt {
			return fmt.Errorf("auto-increment attribute only valid for int fields, got %s",
				dataTypeToString(field.DataType))
		}
		if attrs&AttrRequired == 0 {
			return fmt.Errorf("auto-increment field '%s' should be required", field.Name)
		}
	}

	// hash attribute validation
	if attrs&AttrHash != 0 {
		if field.DataType != DataText {
			return fmt.Errorf("hash attribute only valid for text fields, got %s",
				dataTypeToString(field.DataType))
		}
	}

	// readonly and default don't make sense together typically
	if attrs&AttrReadonly != 0 && attrs&AttrDefault != 0 {
		return fmt.Errorf("readonly and default attributes conflict for field '%s'", field.Name)
	}

	return nil
}

// helper function to convert field kind to string for error messages
func fieldKindToString(kind FieldKind) string {
	switch kind {
	case FieldPrimitive:
		return "primitive"
	case FieldEmbedded:
		return "embedded"
	case FieldReference:
		return "reference"
	default:
		return "unknown"
	}
}

// helper function to convert data type to string for error messages
func dataTypeToString(dt DataType) string {
	switch dt {
	case DataText:
		return "text"
	case DataInt:
		return "int"
	case DataReal:
		return "real/float"
	case DataUUID:
		return "uuid"
	case DataTimestamp:
		return "timestamp"
	case DataBool:
		return "bool"
	case DataEnum:
		return "enum"
	default:
		return "unknown"
	}
}

// batch validation function for multiple fields
func ValidateFields(fields []*Field) []error {
	var errors []error

	for _, field := range fields {
		if err := ValidateFieldAttributes(field); err != nil {
			errors = append(errors, err)
		}
	}

	return errors
}
