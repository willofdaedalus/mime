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
			name: "trailing garbage after end",
			input: `entity user ->
  name text
end something`,
			expected: &entityNode{
				name: "user",
				fields: []field{
					{
						name:        "name",
						dt:          dataText,
						constraints: nil,
						enums:       nil,
					},
				},
			},
		},
		{
			// this test is fine for the entityParser
			// but the parser level will reject this for
			// duplicate field names
			name: "duplicated field names",
			input: `entity user ->
  name text
  name number
end`,
			expected: nil,
		},
		{
			name: "empty enum list",
			input: `entity user ->
  status text ()
end`,
			expected: nil,
		},
		{
			name: "stray token between fields",
			input: `entity user ->
  name text
  @
  age number
end`,
			expected: nil,
		},
		{
			name: "multiple entities in one input",
			input: `entity user ->
  name text
end
entity post ->
  title text
end`,
			expected: &entityNode{ // optional: only if you're parsing one at a time
				name: "user",
				fields: []field{
					{name: "name", dt: dataText},
				},
			},
		},
		{
			name: "missing end keyword",
			input: `entity user ->
  name text`,
			expected: nil,
		},
		{
			name: "invalid data type",
			input: `entity user ->
  name string
end`,
			expected: nil,
		},
		{
			name: "field with missing type",
			input: `entity user ->
  name
end`,
			expected: nil,
		},
		{
			name: "field with missing name",
			input: `entity user ->
  text
end`,
			expected: nil,
		},
		{
			name:     "unterminated entity declaration",
			input:    `entity user ->`,
			expected: nil,
		},
		{
			name: "invalid enum with nested parens",
			input: `entity user ->
  gender text (("male" "female"))
end`,
			expected: nil,
		},
		{
			name: "enum without quotes",
			input: `entity user ->
  gender text (male female)
end`,
			expected: nil,
		},
		{
			name: "enum with mixed types",
			input: `entity user ->
  status text ("active" true)
end`,
			expected: nil,
		},
		{
			name:  "entity with trailing newline",
			input: "entity user ->\nend\n",
			expected: &entityNode{
				name:   "user",
				fields: nil,
			},
		},
		{
			name: "entity with one field and trailing comment",
			input: `entity user ->
  name text # user's full name
end`,
			expected: &entityNode{
				name: "user",
				fields: []field{
					{name: "name", dt: dataText},
				},
			},
		},
		{
			name: "field names with underscores",
			input: `entity user ->
  full_name text
  date_of_birth timestamp
end`,
			expected: &entityNode{
				name: "user",
				fields: []field{
					{name: "full_name", dt: dataText},
					{name: "date_of_birth", dt: dataTimestamp},
				},
			},
		},
		{
			name: "entity with fields and interspersed comments",
			input: `entity user ->
  # name of the user
  name text
  # age is optional
  age int
end`,
			expected: &entityNode{
				name: "user",
				fields: []field{
					{name: "name", dt: dataText},
					{name: "age", dt: dataInt},
				},
			},
		},
		{
			name: "enum with one value",
			input: `entity user ->
  status text ("active")
end`,
			expected: &entityNode{
				name: "user",
				fields: []field{
					{
						name:  "status",
						dt:    dataText,
						enums: []any{"active"},
					},
				},
			},
		},
		{
			name: "empty field list with comments",
			input: `entity user ->
  # this is a comment
  # this is a comment
  # this is a comment
  # this is a comment
  # this is a comment
  # this is a comment
end`,
			expected: &entityNode{
				name:   "user",
				fields: nil,
			},
		},
		{
			name: "invalid modifier syntax",
			input: `entity user ->
  age number { default:18 ensure:"age > 0" }
end`,
			expected: nil, // unquoted default value should fail
		},
		{
			name: "invalid enum type",
			input: `entity user ->
  status text ("active" true)
end`,
			expected: nil, // boolean is not a valid enum
		},
		{
			name: "unterminated modifier block",
			input: `entity user ->
  id number { increment
end`,
			expected: nil,
		},
		{
			name:  "simple entity",
			input: `entity user -> end`,
			expected: &entityNode{
				name:   "user",
				fields: nil,
			},
		},
		{
			name:  "single field",
			input: `entity user -> name text end`,
			expected: &entityNode{
				name: "user",
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
				name: "user",
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
				name: "user",
				fields: []field{
					{
						name:  "gender",
						dt:    dataText,
						enums: []any{"male", "female"},
					},
				},
			},
		},
		{
			name: "field with enums",
			input: `entity user ->
	# this test will return nil because of the 7
	gender text ("male" "female" 7)
end`,
			expected: nil,
		},
		{
			name: "unclosed enums",
			input: `entity user ->
	gender text ("male" "female"
end`,
			expected: nil,
		},
		{
			name: "unclosed string",
			input: `entity user ->
	gender text ("male" "female
end`,
			expected: nil,
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
