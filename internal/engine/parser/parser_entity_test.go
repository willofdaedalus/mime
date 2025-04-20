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
				fields: []longField{
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
  name int
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
  age int
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
				fields: []longField{
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
				fields: []longField{
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
				fields: []longField{
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
				fields: []longField{
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
				fields: []longField{
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
end`,
			expected: &entityNode{
				name:   "user",
				fields: nil,
			},
		},
		{
			name: "invalid modifier syntax",
			input: `entity user ->
  age timestamp { default:18 }
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
  id int { increment
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
				fields: []longField{
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
				fields: []longField{
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
				fields: []longField{
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

func TestParseEntityConstraints(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected *entityNode
	}{
		{
			name: "simple constraint - unique",
			input: `entity user ->
			id int {unique fk}
		end`,
			expected: &entityNode{
				name: "user",
				fields: []longField{
					{
						name: "id",
						dt:   dataInt,
						constraints: []constraint{
							{kind: consUnique},
							{kind: consFK},
						},
					},
				},
			},
		},
		{
			name: "multiple constraints on one field",
			input: `entity user ->
		  id int {unique required}
		end`,
			expected: &entityNode{
				name: "user",
				fields: []longField{
					{
						name: "id",
						dt:   dataInt,
						constraints: []constraint{
							{kind: consUnique},
							{kind: consRequired},
						},
					},
				},
			},
		},
		{
			name: "primary key constraint",
			input: `entity user ->
			  id int {primary}
			end`,
			expected: &entityNode{
				name: "user",
				fields: []longField{
					{
						name: "id",
						dt:   dataInt,
						constraints: []constraint{
							{kind: consPrimary},
						},
					},
				},
			},
		},
		{
			name: "autoincrement constraint",
			input: `entity user ->
			  id int {increment}
			end`,
			expected: &entityNode{
				name: "user",
				fields: []longField{
					{
						name: "id",
						dt:   dataInt,
						constraints: []constraint{
							{kind: consIncrement},
						},
					},
				},
			},
		},
		{
			name: "constraint with unexpected tokens",
			input: `entity user ->
			  id int {primary 123}
			end`,
			expected: nil,
		},
		{
			name: "default constraint with improper value for float",
			input: `entity user ->
			  balance float {default:"abc"}
			end`,
			expected: nil,
		},
		{
			name: "foreign key constraint",
			input: `entity post ->
			  user_id int {fk}
			end`,
			expected: &entityNode{
				name: "post",
				fields: []longField{
					{
						name: "user_id",
						dt:   dataInt,
						constraints: []constraint{
							{kind: consFK},
						},
					},
				},
			},
		},
		{
			name: "unclosed constraint",
			input: `entity user ->
			  id int {unique
			end`,
			expected: nil,
		},
		{
			name: "invalid constraint",
			input: `entity user ->
			  id int {unknown}
			end`,
			expected: nil,
		},
		{
			name: "constraint on unsupported type",
			input: `entity user ->
			  created_at timestamp {unique}
			end`,
			expected: nil,
		},
		{
			name: "constraint with missing value",
			input: `entity user ->
			  status text {default:}
			end`,
			expected: nil,
		},
		{
			name: "constraint with value on constraint that doesn't support values",
			input: `entity user ->
			  id int {unique:"yes"}
			end`,
			expected: nil,
		},
		{
			name: "type mismatch in default value",
			input: `entity user ->
			  age int {default:"not-a-int"}
			end`,
			expected: nil,
		},
		{
			name: "default constraint with value",
			input: `entity user ->
			active int {default:"1"}
		end`,
			expected: &entityNode{
				name: "user",
				fields: []longField{
					{
						name: "active",
						dt:   dataInt,
						constraints: []constraint{
							{
								kind:  consDefault,
								value: stringPtr("1"),
							},
						},
					},
				},
			},
		},
		{
			name: "default constraint for text field",
			input: `entity user ->
			  status text {default:"active"}
			end`,
			expected: &entityNode{
				name: "user",
				fields: []longField{
					{
						name: "status",
						dt:   dataText,
						constraints: []constraint{
							{
								kind:  consDefault,
								value: stringPtr("active"),
							},
						},
					},
				},
			},
		},
		{
			name: "multiple constraints with a default value",
			input: `entity user ->
			  age int {required default:"18"}
			end`,
			expected: &entityNode{
				name: "user",
				fields: []longField{
					{
						name: "age",
						dt:   dataInt,
						constraints: []constraint{
							{kind: consRequired},
							{
								kind:  consDefault,
								value: stringPtr("18"),
							},
						},
					},
				},
			},
		},
		{
			name: "constraints with enum",
			input: `entity user ->
			  role text ("admin" "user") {unique}
			end`,
			expected: &entityNode{
				name: "user",
				fields: []longField{
					{
						name:  "role",
						dt:    dataText,
						enums: []any{"admin", "user"},
						constraints: []constraint{
							{
								kind: consUnique,
							},
						},
					},
				},
			},
		},
		{
			name: "newline within constraint block",
			input: `entity user ->
			  id int {
			    primary
			  }
			end`,
			expected: nil,
		},
		{
			name: "default constraint with proper value for float",
			input: `entity user ->
			  balance float {default:"123.45"}
			end`,
			expected: &entityNode{
				name: "user",
				fields: []longField{
					{
						name: "balance",
						dt:   dataReal,
						constraints: []constraint{
							{
								kind:  consDefault,
								value: stringPtr("123.45"),
							},
						},
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

// Helper function to create string pointers for the tests
func stringPtr(s string) *string {
	return &s
}
