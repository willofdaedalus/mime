package parser

import (
	"reflect"
	"testing"

	"willofdaedalus/mime/internal/engine/lexer"
)

func TestParseEntity(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected *entityNode
	}{
		{
			name:  "simple entity",
			input: `entity user -> end`,
			expected: &entityNode{
				entityName: "user",
				fields:     nil,
			},
		},
		{
			name:  "single field",
			input: `entity user -> name text end`,
			expected: &entityNode{
				entityName: "user",
				fields: []field{
					{name: "name", dt: dataText},
				},
			},
		},
		{
			name: "multiple fields",
			input: `entity user ->
	name text
	age int
end`,
			expected: &entityNode{
				entityName: "user",
				fields: []field{
					{name: "name", dt: dataText},
					{name: "age", dt: dataInt},
				},
			},
		},
		{
			name: "field with enums",
			input: `entity user ->
	gender text ("male" "female")
end`,
			expected: &entityNode{
				entityName: "user",
				fields: []field{
					{
						name:  "gender",
						dt:    dataText,
						enums: []any{"male", "female"},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser(lexer.New(tt.input))
			actual := p.parseEntity()

			if !reflect.DeepEqual(actual, tt.expected) {
				t.Fatalf("for %s:\nexpected:\n%v\ngot:\n%v", tt.name, tt.expected, actual)
			}
		})
	}
}
