package parser

import "testing"

func TestNextToken(t *testing.T) {
	input := `entity user ->
	id: number <> [1, 2]
	dob: text
	gender: text <> ["male", "female"]
`

	tests := []struct {
		expectedType    tokenType
		expectedLiteral string
	}{
		{TK_ENTITY, "entity"},
		{TK_IDENT, "user"},
		{TK_DASH, "-"},
		{TK_RANGLE, ">"},
		{TK_IDENT, "id"},
		{TK_COLON, ":"},
		{TK_NUMBER, "number"},
		{TK_LANGLE, "<"},
		{TK_RANGLE, ">"},
		{TK_LBRACKET, "["},
		{TK_DIGITS, "1"},
		{TK_COMMA, ","},
		{TK_DIGITS, "2"},
		{TK_RBRACKET, "]"},
		{TK_IDENT, "dob"},
		{TK_COLON, ":"},
		{TK_TEXT, "text"},
		{TK_IDENT, "gender"},
		{TK_COLON, ":"},
		{TK_TEXT, "text"},
		{TK_LANGLE, "<"},
		{TK_RANGLE, ">"},
		{TK_LBRACKET, "["},
		{TK_STRING, "\"male\""},
		{TK_COMMA, ","},
		{TK_STRING, "\"female\""},
		{TK_RBRACKET, "]"},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokenType wrong. expected=%q, got=%q",
				i, tt.expectedType.String(), tok.Type.String())
		}

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokenType wrong. expected=%q, got=%q",
				i, tt.expectedType.String(), tok.Type.String())
		}
	}
}
