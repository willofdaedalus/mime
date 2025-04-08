package lexer

import "testing"

func TestNextToken(t *testing.T) {
	input := `entity user ->
	id int (1, 2)
	dob text
	gender text ("male", "female")
end

# this is a comment and shouldn't be tokenized
routes ->
	GET /users/me -> self.id
end
`

	tests := []struct {
		expectedType    TokenType
		expectedLiteral string
	}{
		{TokenEntity, "entity"},
		{TokenIdent, "user"},
		{TokenArrow, "->"},
		{TokenIdent, "id"},
		{TokenInt, "int"},
		{TokenEnumOpen, "("},
		{TokenDigits, "1"},
		{TokenDigits, "2"},
		{TokenEnumClose, ")"},
		{TokenIdent, "dob"},
		{TokenText, "text"},
		{TokenIdent, "gender"},
		{TokenText, "text"},
		{TokenEnumOpen, "("},
		{TokenString, "\"male\""},
		{TokenString, "\"female\""},
		{TokenEnumClose, ")"},
		{TokenEnd, "end"},
		{TokenComment, "#"},
		{TokenRoutes, "routes"},
		{TokenArrow, "->"},
		{TokenGet, "GET"},
		{TokenEndpoint, "/users/me"},
		{TokenArrow, "->"},
		{TokenSelf, "self"},
		{TokenDot, "."},
		{TokenIdent, "id"},
		{TokenEnd, "end"},
	}
	// routes ->
	// 	GET /users/me -> self.id

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
