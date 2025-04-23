package parser

import "willofdaedalus/mime/internal/engine/lexer"

func expectTokOf(tok lexer.Token, tokType lexer.TokenType) bool {
	return tok.Type == tokType
}
