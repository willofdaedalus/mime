package parser

import (
	"testing"

	l "willofdaedalus/mime/internal/engine/lexer"
	"willofdaedalus/mime/internal/engine/types"
)

func TestEnumParsing(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected *types.EnumNode
		wantErr  bool
	}{
		{
			name: "simple enum",
			input: `enum user_role ->
	admin
	user
end`,
			expected: &types.EnumNode{
				Name:    "user_role",
				Members: []string{"admin", "user"},
			},
			wantErr: false,
		},
		{
			name: "enum with many members",
			input: `enum status ->
	pending
	approved
	rejected
	cancelled
	processing
end`,
			expected: &types.EnumNode{
				Name:    "status",
				Members: []string{"pending", "approved", "rejected", "cancelled", "processing"},
			},
			wantErr: false,
		},
		{
			name: "enum with comments",
			input: `enum priority ->
	// high priority items
	high
	medium // default priority
	low
end`,
			expected: &types.EnumNode{
				Name:    "priority",
				Members: []string{"high", "medium", "low"},
			},
			wantErr: false,
		},
		{
			name:     "empty enum should fail",
			input:    `enum empty_enum ->\nend`,
			expected: nil,
			wantErr:  true,
		},
		{
			name:     "enum without arrow should fail",
			input:    `enum bad_enum\n\tadmin\nend`,
			expected: nil,
			wantErr:  true,
		},
		{
			name:     "enum without end keyword should fail",
			input:    `enum no_end ->\n\tadmin\n\tuser`,
			expected: nil,
			wantErr:  true,
		},
		{
			name:     "enum without name should fail",
			input:    `enum ->\n\tadmin\nend`,
			expected: nil,
			wantErr:  true,
		},
		{
			name:     "enum with invalid member tokens should fail",
			input:    `enum test_enum ->\n\t123invalid\nend`,
			expected: nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := l.New(tt.input)
			parser := NewParser(lexer)

			result, err := handleEnum(parser)

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

			enumNode, ok := result.(*types.EnumNode)
			if !ok {
				t.Errorf("expected *types.EnumNode, got %T", result)
				return
			}

			if enumNode.Name != tt.expected.Name {
				t.Errorf("expected name %s, got %s", tt.expected.Name, enumNode.Name)
			}

			if len(enumNode.Members) != len(tt.expected.Members) {
				t.Errorf("expected %d members, got %d", len(tt.expected.Members), len(enumNode.Members))
				return
			}

			for i, member := range enumNode.Members {
				if member != tt.expected.Members[i] {
					t.Errorf("expected member[%d] = %s, got %s", i, tt.expected.Members[i], member)
				}
			}
		})
	}
}
