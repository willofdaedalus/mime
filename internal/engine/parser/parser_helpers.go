package parser

import (
	"slices"

	"willofdaedalus/mime/internal/engine/lexer"
)

func expectTokOf(tok lexer.Token, tokType lexer.TokenType) bool {
	return tok.Type == tokType
}

func skipToTok(p *Parser, targets ...lexer.TokenType) {
	for {
		if slices.Contains(targets, p.curToken.Type) ||
			p.curToken.Type == lexer.TokenEOF {
			return
		}
		p.advanceToken()
	}
}
