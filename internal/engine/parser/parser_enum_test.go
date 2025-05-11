package parser

import (
	"reflect"
	"testing"

	"willofdaedalus/mime/internal/engine/lexer"
	"willofdaedalus/mime/internal/engine/types"
)

func TestEnumHandler(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected *types.EnumNode
	}{
		{
			name: "simple enum with valid members",
			input: `enum role ->
			admin
			user
			Admin
			end`,
			expected: &types.EnumNode{
				Name:    "role",
				Members: []string{"admin", "user", "Admin"},
			},
		},
		{
			name: "duplicate enum member",
			input: `enum role ->
			admin
			user
			admin
			end`,
			expected: nil,
		},
		{
			name: "empty enum list",
			input: `enum role ->
			end`,
			expected: &types.EnumNode{
				Name:    "role",
				Members: []string{},
			},
		},
		{
			name: "member starts with digit",
			input: `enum role ->
			1admin
			end`,
			expected: nil,
		},
		{
			name: "member starts with underscore",
			input: `enum role ->
			# this is a comment
			_admin123
			end`,
			expected: &types.EnumNode{
				Name:    "role",
				Members: []string{"_admin123"},
			},
		},
		{
			name: "member starts with symbol",
			input: `enum role ->
			#admin
			end`,
			expected: &types.EnumNode{
				Name:    "role",
				Members: []string{},
			},
		},
		{
			name: "enum member is a keyword",
			input: `enum role ->
			enum
			end`,
			expected: nil,
		},
		{
			name: "member starts with underscore",
			input: `enum role ->
			_admin
			end`,
			expected: &types.EnumNode{
				Name:    "role",
				Members: []string{"_admin"},
			},
		},
		{
			name: "unexpected token between members",
			input: `enum role ->
			admin
			@
			user
			end`,
			expected: nil,
		},
		{
			name: "mixed valid and invalid members",
			input: `enum role ->
			admin
			1user
			user
			admin
			end`,
			expected: nil,
		},
		{
			name: "early EOF before end",
			input: `enum role ->
			admin
			user`,
			expected: nil,
		},
		{
			name: "valid enum followed by garbage",
			input: `enum role ->
			admin
			user
			end garbage`,
			expected: &types.EnumNode{
				Name:    "role",
				Members: []string{"admin", "user"},
			},
		},

		// --- Additional Tests ---

		{
			name: "enum with only whitespace",
			input: `enum role ->
			
			
			end`,
			expected: &types.EnumNode{
				Name:    "role",
				Members: []string{},
			},
		},
		{
			name: "enum with trailing comment after end",
			input: `enum role ->
			admin
			user
			end # no more roles`,
			expected: &types.EnumNode{
				Name:    "role",
				Members: []string{"admin", "user"},
			},
		},
		{
			name: "enum with inline comment on member line",
			input: `enum role ->
			admin # highest privileges
			user
			end`,
			expected: &types.EnumNode{
				Name:    "role",
				Members: []string{"admin", "user"},
			},
		},
		{
			name:  "enum with windows line endings",
			input: "enum role ->\r\nadmin\r\nuser\r\nend",
			expected: &types.EnumNode{
				Name:    "role",
				Members: []string{"admin", "user"},
			},
		},
		{
			name: "enum name with invalid characters",
			input: `enum ro$le ->
			admin
			end`,
			expected: nil,
		},
		{
			name: "enum with case-insensitive duplicates",
			input: `enum role ->
			admin
			Admin
			end`,
			expected: &types.EnumNode{
				Name:    "role",
				Members: []string{"admin", "Admin"},
			},
		},
		{
			name: "enum missing name",
			input: `enum ->
			admin
			end`,
			expected: nil,
		},
		{
			name: "enum preceded by garbage",
			input: `random garbage
			enum role ->
			admin
			end`,
			expected: nil,
		},
		{
			name: "enum ends with capital END",
			input: `enum role ->
			admin
			END`,
			expected: nil,
		},
		{
			name: "enum with quoted string as member",
			input: `enum role ->
			"admin"
			end`,
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser(lexer.New(tt.input))
			actual := handleEnum(p)

			if !reflect.DeepEqual(actual, tt.expected) {
				t.Fatalf("for test %s:\nexpected:\n%#v\ngot:\n%#v", tt.name, tt.expected, actual)
			}
		})
	}
}
