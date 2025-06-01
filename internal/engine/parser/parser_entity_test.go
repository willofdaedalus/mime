package parser

import (
	"fmt"
	"strings"
	"testing"
	"time"

	l "willofdaedalus/mime/internal/engine/lexer"
	"willofdaedalus/mime/internal/engine/types"
)

func TestEntityParsing(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected *types.EntityNode
		wantErr  bool
	}{
		{
			name: "simple entity with basic fields",
			input: `entity user ->
			id uuid
			name text
			age int
		end`,
			expected: &types.EntityNode{
				Name: "user",
				Fields: []*types.Field{
					{Name: "id", Kind: types.FieldPrimitive, DataType: types.DataUUID},
					{Name: "name", Kind: types.FieldPrimitive, DataType: types.DataText},
					{Name: "age", Kind: types.FieldPrimitive, DataType: types.DataInt},
				},
			},
			wantErr: false,
		},
		{
			name: "entity with attributes",
			input: `entity user ->
			id uuid [primary required]
			name text [required unique]
			password text [hash]
			age int [default]
		end`,
			expected: &types.EntityNode{
				Name: "user",
				Fields: []*types.Field{
					{
						Name: "id", Kind: types.FieldPrimitive, DataType: types.DataUUID,
						Attributes: types.AttrPrimary | types.AttrRequired,
					},
					{
						Name: "name", Kind: types.FieldPrimitive, DataType: types.DataText,
						Attributes: types.AttrRequired | types.AttrUnique,
					},
					{
						Name: "password", Kind: types.FieldPrimitive, DataType: types.DataText,
						Attributes: types.AttrHash,
					},
					{
						Name: "age", Kind: types.FieldPrimitive, DataType: types.DataInt,
						Attributes: types.AttrDefault,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "entity with reference field",
			input: `entity note ->
			id uuid
			title text
			owner @user.id
		end`,
			expected: &types.EntityNode{
				Name: "note",
				Fields: []*types.Field{
					{Name: "id", Kind: types.FieldPrimitive, DataType: types.DataUUID},
					{Name: "title", Kind: types.FieldPrimitive, DataType: types.DataText},
					{
						Name: "owner", Kind: types.FieldReference,
						Target: &types.ReferenceTarget{Entity: "user", Field: "id"},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "entity with enum reference",
			input: `entity user ->
			id uuid
			name text
			role &user_role
		end`,
			expected: &types.EntityNode{
				Name: "user",
				Fields: []*types.Field{
					{Name: "id", Kind: types.FieldPrimitive, DataType: types.DataUUID},
					{Name: "name", Kind: types.FieldPrimitive, DataType: types.DataText},
					{
						Name: "role", Kind: types.FieldEnum, DataType: types.DataEnum,
						Target: &types.ReferenceTarget{Entity: "user_role"},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "entity with embedded entity",
			input: `entity student ->
			@person
			gpa float
			course text
		end`,
			expected: &types.EntityNode{
				Name: "student",
				Fields: []*types.Field{
					{Name: "person", Kind: types.FieldEmbedded},
					{Name: "gpa", Kind: types.FieldPrimitive, DataType: types.DataReal},
					{Name: "course", Kind: types.FieldPrimitive, DataType: types.DataText},
				},
			},
			wantErr: false,
		},
		{
			name: "complex entity with mixed field types",
			input: `entity order ->
	id uuid [primary]
	@audit_info
	customer @user.id [required]
	status &order_status
	total float [required]
	created timestamp
end`,
			expected: nil,
			wantErr:  true,
		},
		{
			name:     "entity without name should fail",
			input:    `entity ->\n\tid uuid\nend`,
			expected: nil,
			wantErr:  true,
		},
		{
			name:     "entity without arrow should fail",
			input:    `entity test\n\tid uuid\nend`,
			expected: nil,
			wantErr:  true,
		},
		{
			name:     "entity with invalid field should fail",
			input:    `entity test ->\n\t123invalid\nend`,
			expected: nil,
			wantErr:  true,
		},
		{
			name:     "entity with malformed reference should fail",
			input:    `entity test ->\n\towner @user\nend`,
			expected: nil,
			wantErr:  true,
		},
		{
			name:     "entity with malformed enum reference should fail",
			input:    `entity test ->\n\trole &\nend`,
			expected: nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := l.New(tt.input)
			parser := NewParser(lexer)

			result, err := handleEntity(parser)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			entityNode, ok := result.(*types.EntityNode)
			if !ok {
				t.Errorf("expected *types.EntityNode, got %T", result)
				return
			}

			if entityNode.Name != tt.expected.Name {
				t.Errorf("expected name %s, got %s", tt.expected.Name, entityNode.Name)
			}

			if len(entityNode.Fields) != len(tt.expected.Fields) {
				t.Errorf("expected %d fields, got %d", len(tt.expected.Fields), len(entityNode.Fields))
				return
			}

			for i, field := range entityNode.Fields {
				expected := tt.expected.Fields[i]

				if field.Name != expected.Name {
					t.Errorf("field[%d]: expected name %s, got %s", i, expected.Name, field.Name)
				}

				if field.Kind != expected.Kind {
					t.Errorf("field[%d]: expected kind %v, got %v", i, expected.Kind, field.Kind)
				}

				if field.DataType != expected.DataType {
					t.Errorf("field[%d]: expected datatype %v, got %v", i, expected.DataType, field.DataType)
				}

				if field.Attributes != expected.Attributes {
					t.Errorf("field[%d]: expected attributes %v, got %v", i, expected.Attributes, field.Attributes)
				}

				// Check reference targets
				if expected.Target != nil {
					if field.Target == nil {
						t.Errorf("field[%d]: expected target but got nil", i)
						continue
					}
					if field.Target.Entity != expected.Target.Entity {
						t.Errorf("field[%d]: expected target entity %s, got %s", i, expected.Target.Entity, field.Target.Entity)
					}
					if field.Target.Field != expected.Target.Field {
						t.Errorf("field[%d]: expected target field %s, got %s", i, expected.Target.Field, field.Target.Field)
					}
				} else if field.Target != nil {
					t.Errorf("field[%d]: expected no target but got %+v", i, field.Target)
				}
			}
		})
	}
}

func TestFieldParsing(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected *types.Field
		wantErr  bool
	}{
		{
			name:  "simple text field",
			input: "name text",
			expected: &types.Field{
				Name:     "name",
				Kind:     types.FieldPrimitive,
				DataType: types.DataText,
			},
			wantErr: false,
		},
		{
			name:  "uuid field with primary attribute",
			input: "id uuid [primary]",
			expected: &types.Field{
				Name:       "id",
				Kind:       types.FieldPrimitive,
				DataType:   types.DataUUID,
				Attributes: types.AttrPrimary,
			},
			wantErr: false,
		},
		{
			name:  "text field with multiple attributes",
			input: "username text [required unique]",
			expected: &types.Field{
				Name:       "username",
				Kind:       types.FieldPrimitive,
				DataType:   types.DataText,
				Attributes: types.AttrRequired | types.AttrUnique,
			},
			wantErr: false,
		},
		{
			name:  "reference field",
			input: "owner @user.id",
			expected: &types.Field{
				Name:   "owner",
				Kind:   types.FieldReference,
				Target: &types.ReferenceTarget{Entity: "user", Field: "id"},
			},
			wantErr: false,
		},
		{
			name:  "enum reference field",
			input: "status &order_status",
			expected: &types.Field{
				Name:     "status",
				Kind:     types.FieldEnum,
				DataType: types.DataEnum,
				Target:   &types.ReferenceTarget{Entity: "order_status"},
			},
			wantErr: false,
		},
		{
			name:  "embedded field",
			input: "@person",
			expected: &types.Field{
				Name: "person",
				Kind: types.FieldEmbedded,
			},
			wantErr: false,
		},
		{
			name:  "int field with increment",
			input: "counter int [increment]",
			expected: &types.Field{
				Name:       "counter",
				Kind:       types.FieldPrimitive,
				DataType:   types.DataInt,
				Attributes: types.AttrIncrement,
			},
			wantErr: false,
		},
		{
			name:     "field without name should fail",
			input:    "text",
			expected: nil,
			wantErr:  true,
		},
		{
			name:     "field with invalid type should fail",
			input:    "name invalidtype",
			expected: nil,
			wantErr:  true,
		},
		{
			name:     "reference without dot should fail",
			input:    "owner @user",
			expected: nil,
			wantErr:  true,
		},
		{
			name:     "reference without field should fail",
			input:    "owner @user.",
			expected: nil,
			wantErr:  true,
		},
		{
			name:     "enum reference without name should fail",
			input:    "status &",
			expected: nil,
			wantErr:  true,
		},
		{
			name:     "embedded without name should fail",
			input:    "@",
			expected: nil,
			wantErr:  true,
		},
		{
			name:     "unclosed attributes should fail",
			input:    "name text [required",
			expected: nil,
			wantErr:  true,
		},
		{
			name:     "invalid attribute should fail",
			input:    "name text [invalidattr]",
			expected: nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Add newline to simulate proper field ending
			input := tt.input + "\n"
			lexer := l.New(input)
			parser := NewParser(lexer)

			result, err := parseField(parser)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if result.Name != tt.expected.Name {
				t.Errorf("expected name %s, got %s", tt.expected.Name, result.Name)
			}

			if result.Kind != tt.expected.Kind {
				t.Errorf("expected kind %v, got %v", tt.expected.Kind, result.Kind)
			}

			if result.DataType != tt.expected.DataType {
				t.Errorf("expected datatype %v, got %v", tt.expected.DataType, result.DataType)
			}

			if result.Attributes != tt.expected.Attributes {
				t.Errorf("expected attributes %v, got %v", tt.expected.Attributes, result.Attributes)
			}

			// Check reference targets
			if tt.expected.Target != nil {
				if result.Target == nil {
					t.Errorf("expected target but got nil")
					return
				}
				if result.Target.Entity != tt.expected.Target.Entity {
					t.Errorf("expected target entity %s, got %s", tt.expected.Target.Entity, result.Target.Entity)
				}
				if result.Target.Field != tt.expected.Target.Field {
					t.Errorf("expected target field %s, got %s", tt.expected.Target.Field, result.Target.Field)
				}
			} else if result.Target != nil {
				t.Errorf("expected no target but got %+v", result.Target)
			}
		})
	}
}

func TestAttributeParsing(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected types.Attribute
		wantErr  bool
	}{
		{
			name:     "single attribute",
			input:    "[required]",
			expected: types.AttrRequired,
			wantErr:  false,
		},
		{
			name:     "multiple attributes",
			input:    "[required unique primary]",
			expected: types.AttrRequired | types.AttrUnique | types.AttrPrimary,
			wantErr:  false,
		},
		{
			name:  "all valid attributes",
			input: "[default hash unique required increment override primary hidden readonly]",
			expected: types.AttrDefault | types.AttrHash | types.AttrUnique | types.AttrRequired |
				types.AttrIncrement | types.AttrOverride | types.AttrPrimary | types.AttrHidden | types.AttrReadonly,
			wantErr: false,
		},
		{
			name:     "empty attributes should work",
			input:    "[]",
			expected: 0,
			wantErr:  false,
		},
		{
			name:     "unclosed bracket should fail",
			input:    "[required",
			expected: 0,
			wantErr:  true,
		},
		{
			name:     "invalid attribute should fail",
			input:    "[invalidattr]",
			expected: 0,
			wantErr:  true,
		},
		{
			name:     "no opening bracket should fail",
			input:    "required]",
			expected: 0,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := l.New(tt.input)
			parser := NewParser(lexer)

			result, err := parseAttributes(parser)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if result != tt.expected {
				t.Errorf("expected attributes %v, got %v", tt.expected, result)
			}
		})
	}
}

// Stress tests to try and break the parser
func TestParserStressTests(t *testing.T) {
	stressTests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name: "deeply nested references",
			input: `entity deep ->
	ref1 @level1.ref2
	ref2 @level2.ref3
	ref3 @level3.ref4
end`,
			wantErr: false,
		},
		{
			name: "entity with many fields",
			input: `entity huge ->
	field1 text
	field2 int
	field3 uuid
	field4 bool
	field5 timestamp
	field6 float
	field7 text [required]
	field8 int [unique]
	field9 uuid [primary]
	@embedded1
	@embedded2
	field12 &enum1
	field13 &enum2
	ref1 @entity1.id
	ref2 @entity2.name
end`,
			wantErr: false,
		},
		{
			name: "mixed comments everywhere",
			input: `# comment before entity
entity test -> # comment after arrow
	# comment before field
	id uuid [primary] # comment after field
	# another comment
	name text # final comment
	# comment before end
end # comment after end`,
			wantErr: false,
		},
		{
			name: "empty lines and whitespace",
			input: `entity test ->

	id uuid

	name text

end`,
			wantErr: false,
		},
		{
			name:    "malformed nested brackets",
			input:   `entity test ->\n\tid uuid [required [nested]]\nend`,
			wantErr: true,
		},
		{
			name:    "missing keywords",
			input:   `test ->\n\tid uuid\nend`,
			wantErr: true,
		},
		{
			name:    "unterminated entity",
			input:   `entity test ->\n\tid uuid`,
			wantErr: true,
		},
	}

	for _, tt := range stressTests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := l.New(tt.input)
			parser := NewParser(lexer)

			_, err := handleEntity(parser)

			if tt.wantErr && err == nil {
				t.Errorf("expected error but got none for input: %s", tt.input)
			}

			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error: %v for input: %s", err, tt.input)
			}
		})
	}
}

// Additional stress tests to add to your existing test suite

func TestParserEdgeCases(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		// Unicode and special characters
		// 		{
		// 			name: "unicode field names",
		// 			input: `entity tëst ->
		// 	naïve_field text
		// 	café_name text
		// end`,
		// 			wantErr: true, // Unicode not supported
		// 		},
		{
			name: "very long identifiers",
			input: `entity this_is_a_very_long_entity_name_that_might_break_something ->
	this_is_also_a_very_long_field_name_that_could_cause_issues text
end`,
			wantErr: false,
		},

		// Whitespace edge cases
		{
			name:    "tabs vs spaces mixing",
			input:   "entity test ->\n\tid uuid\n    name text\nend", // mixed tabs and spaces
			wantErr: false,                                           // depends on your lexer
		},
		{
			name:    "trailing whitespace after tokens",
			input:   "entity test ->   \n\tid uuid   \n\tname text   \nend   ",
			wantErr: false,
		},
		{
			name:    "windows line endings",
			input:   "entity test ->\r\n\tid uuid\r\n\tname text\r\nend",
			wantErr: false,
		},

		// Comment edge cases
		{
			name: "comments with special characters",
			input: `entity test ->
	# This is a comment with @#$%^&*()
	id uuid # Another comment with symbols !@#
	name text # Comment with unicode: café naïve
end`,
			wantErr: false,
		},
		{
			name: "very long comments",
			input: `entity test ->
	# ` + strings.Repeat("This is a very long comment ", 100) + `
	id uuid
end`,
			wantErr: false,
		},

		// Boundary conditions
		{
			name:    "single character entity name",
			input:   "entity a ->\n\tid uuid\nend",
			wantErr: false,
		},
		{
			name:    "single character field name",
			input:   "entity test ->\n\tx text\nend",
			wantErr: false,
		},
		{
			name:    "numbers in identifiers",
			input:   "entity user2 ->\n\tid2 uuid\n\tfield_3 text\nend",
			wantErr: false,
		},

		// Malformed brackets and symbols
		{
			name:    "unmatched opening bracket",
			input:   "entity test ->\n\tid uuid [\nend",
			wantErr: true,
		},
		{
			name:    "unmatched closing bracket",
			input:   "entity test ->\n\tid uuid ]\nend",
			wantErr: true,
		},
		{
			name:    "double opening brackets",
			input:   "entity test ->\n\tid uuid [[required]]\nend",
			wantErr: true,
		},
		{
			name:    "empty reference target",
			input:   "entity test ->\n\tref @\nend",
			wantErr: true,
		},
		{
			name:    "multiple dots in reference",
			input:   "entity test ->\n\tref @user..id\nend",
			wantErr: true,
		},
		{
			name:    "reference starting with dot",
			input:   "entity test ->\n\tref @.user.id\nend",
			wantErr: true,
		},

		// // Enum edge cases
		// {
		// 	name:    "duplicate enum members",
		// 	input:   "enum status ->\n\tpending\n\tpending\nend",
		// 	wantErr: true, // should probably be an error
		// },
		{
			name:    "enum member with numbers",
			input:   "enum status ->\n\tstatus_1\n\tstatus_2\nend",
			wantErr: false,
		},

		// Attribute edge cases
		{
			name:    "duplicate attributes",
			input:   "entity test ->\n\tid uuid [required required]\nend",
			wantErr: false, // might be ok, just redundant
		},
		{
			name:    "attributes with extra spaces",
			input:   "entity test ->\n\tid uuid [ required  unique ]\nend",
			wantErr: false,
		},
		{
			name:    "attributes on wrong field types",
			input:   "entity test ->\n\t@embedded [required]\nend",
			wantErr: true, // embedded fields shouldn't have attributes
		},

		// Circular references (if you want to catch these at parse time)
		{
			name: "self reference",
			input: `entity user ->
	id uuid
	parent @user.id
end`,
			wantErr: false, // this is actually valid
		},

		// Case sensitivity tests
		{
			name:    "mixed case keywords",
			input:   "Entity test ->\n\tid UUID\nEnd",
			wantErr: true, // assuming keywords are case sensitive
		},
		{
			name:    "mixed case attributes",
			input:   "entity test ->\n\tid uuid [Required]\nend",
			wantErr: true, // assuming attributes are case sensitive
		},

		// Empty input edge cases
		{
			name:    "completely empty input",
			input:   "",
			wantErr: true,
		},
		{
			name:    "only whitespace",
			input:   "   \n\t\n   ",
			wantErr: true,
		},
		{
			name:    "only comments",
			input:   "# just a comment\n# another comment",
			wantErr: true,
		},

		// Malformed arrows and keywords
		{
			name:    "wrong arrow direction",
			input:   "entity test <-\n\tid uuid\nend",
			wantErr: true,
		},
		{
			name:    "multiple arrows",
			input:   "entity test -> ->\n\tid uuid\nend",
			wantErr: true,
		},
		{
			name:    "space in arrow",
			input:   "entity test - >\n\tid uuid\nend",
			wantErr: true,
		},

		// Premature EOF
		{
			name:    "EOF after entity keyword",
			input:   "entity",
			wantErr: true,
		},
		{
			name:    "EOF after entity name",
			input:   "entity test",
			wantErr: true,
		},
		{
			name:    "EOF after arrow",
			input:   "entity test ->",
			wantErr: true,
		},
		{
			name:    "EOF in middle of field",
			input:   "entity test ->\n\tid",
			wantErr: true,
		},
		{
			name:    "EOF in middle of attributes",
			input:   "entity test ->\n\tid uuid [required",
			wantErr: true,
		},

		// Very large constructs
		{
			name: "entity with maximum reasonable fields",
			input: func() string {
				var sb strings.Builder
				sb.WriteString("entity huge ->\n")
				for i := 0; i < 1000; i++ {
					sb.WriteString(fmt.Sprintf("\tfield_%d text\n", i))
				}
				sb.WriteString("end")
				return sb.String()
			}(),
			wantErr: false,
		},
		{
			name: "enum with many members",
			input: func() string {
				var sb strings.Builder
				sb.WriteString("enum huge_enum ->\n")
				for i := 0; i < 500; i++ {
					sb.WriteString(fmt.Sprintf("\tmember_%d\n", i))
				}
				sb.WriteString("end")
				return sb.String()
			}(),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := l.New(tt.input)
			parser := NewParser(lexer)

			// Determine what type of construct we're testing based on input
			var result interface{}
			var err error
			if strings.HasPrefix(strings.TrimSpace(tt.input), "entity") {
				result, err = handleEntity(parser)
			} else if strings.HasPrefix(strings.TrimSpace(tt.input), "enum") {
				result, err = handleEnum(parser)
			} else {
				// Default to entity for ambiguous cases
				result, err = handleEntity(parser)
			}

			if tt.wantErr && err == nil {
				t.Errorf("expected error but got none for input: %q", tt.input)
			}

			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error: %v for input: %q", err, tt.input)
			}

			_ = result // avoid unused variable warning
		})
	}
}

// Test concurrent parsing (if your parser needs to be thread-safe)
func TestConcurrentParsing(t *testing.T) {
	input := `entity user ->
	id uuid [primary]
	name text [required]
	email text [unique]
end`

	const numGoroutines = 100
	errors := make(chan error, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			lexer := l.New(input)
			parser := NewParser(lexer)
			_, err := handleEntity(parser)
			errors <- err
		}()
	}

	for i := 0; i < numGoroutines; i++ {
		if err := <-errors; err != nil {
			t.Errorf("concurrent parsing failed: %v", err)
		}
	}
}

// Memory pressure test
func TestParsingMemoryPressure(t *testing.T) {
	// Create a very large but valid entity definition
	var sb strings.Builder
	sb.WriteString("entity stress_test ->\n")

	// Add many different types of fields
	for i := 0; i < 10000; i++ {
		fieldType := []string{"text", "int", "uuid", "bool", "float", "timestamp"}[i%6]
		sb.WriteString(fmt.Sprintf("\tfield_%d %s\n", i, fieldType))

		// Occasionally add references and enums
		if i%100 == 0 {
			sb.WriteString(fmt.Sprintf("\tref_%d @other_entity.id\n", i))
		}
		if i%150 == 0 {
			sb.WriteString(fmt.Sprintf("\tenum_%d &some_enum\n", i))
		}
		if i%200 == 0 {
			sb.WriteString("\t@embedded_entity\n")
		}
	}
	sb.WriteString("end")

	lexer := l.New(sb.String())
	parser := NewParser(lexer)

	start := time.Now()
	result, err := handleEntity(parser)
	duration := time.Since(start)

	if err != nil {
		t.Errorf("memory pressure test failed: %v", err)
	}

	t.Logf("Parsed large entity in %v", duration)

	if entity, ok := result.(*types.EntityNode); ok {
		t.Logf("Successfully parsed entity with %d fields", len(entity.Fields))
	}
}
