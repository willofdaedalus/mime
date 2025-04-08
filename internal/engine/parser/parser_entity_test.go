package parser

import (
	"reflect"
	"testing"

	"willofdaedalus/mime/internal/engine/lexer"
)

func TestParseEntitySimple(t *testing.T) {
	input := `entity user ->
end`
	p := NewParser(lexer.New(input))

	expected := &entityNode{
		entityName: "user",
		fields:     nil,
	}
	actual := p.parseEntity()

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("expected %v got %v", expected, actual)
	}
}

func TestParseEntitySingleField(t *testing.T) {
	input := `entity user ->
	name text
end`
	p := NewParser(lexer.New(input))

	expected := &entityNode{
		entityName: "user",
		fields: []field{
			{
				name:        "name",
				dt:          dataText,
				constraints: nil,
				enums:       nil,
			},
		},
	}
	actual := p.parseEntity()

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("expected %v got %v", expected, actual)
	}
}

func TestParseEntityMultipleFields(t *testing.T) {
	input := `entity user ->
	name text
	age int
end`
	p := NewParser(lexer.New(input))

	expected := &entityNode{
		entityName: "user",
		fields: []field{
			{
				name:        "name",
				dt:          dataText,
				constraints: nil,
				enums:       nil,
			},
			{
				name:        "age",
				dt:          dataInt,
				constraints: nil,
				enums:       nil,
			},
		},
	}
	actual := p.parseEntity()

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("expected %v got %v", expected, actual)
	}
}

// there's a program crashing bug where the if the user uses the wrong opening
// the program freezes during the parsing stage and can't proceed
func TestParseEntitySimpleEnums(t *testing.T) {
	input := `entity user ->
	gender text ("male" "female")
end`
	p := NewParser(lexer.New(input))

	expected := &entityNode{
		entityName: "user",
		fields: []field{
			{
				name:        "gender",
				dt:          dataText,
				constraints: nil,
				enums:       []any{"male", "female"},
			},
		},
	}
	actual := p.parseEntity()

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("expected %v got %v", expected, actual)
	}
}
