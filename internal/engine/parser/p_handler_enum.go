package parser

import (
	"fmt"
	"unicode"

	"willofdaedalus/mime/internal/engine/lexer"
	"willofdaedalus/mime/internal/engine/types"
)

func handleEnum(p *Parser) node {
	var hasStar bool
	enumNode := &types.EnumNode{
		Members: make([]string, 0),
	}

	for p.curToken.Type != lexer.TokenEnd && p.curToken.Type != lexer.TokenEOF {
		if p.curToken.Type == lexer.TokenStar {
			hasStar = true
			p.advanceToken() // consume the star
		}

		if p.curToken.Type == lexer.TokenIdent {
			v := p.curToken.Literal
			if err := validateEnumValue(v, hasStar); err == nil {
				enumNode.Members = append(enumNode.Members, p.curToken.Literal)
			}
			hasStar = false
			p.advanceToken()
		}
	}

	return enumNode
}

func validateEnumValue(s string, star bool) error {
	// check that it's not starting with a number
	if unicode.IsDigit(rune(s[0])) {
		return fmt.Errorf("%s starts with a number", s)
	}
	// next check that it doesn't start with symbol
	// if it does make sure to it's not
	if unicode.IsSymbol(rune(s[0])) && s[0] != '_' {
		return fmt.Errorf("%s starts with a symbol", s)
	}
	// next check that it doesn't conflict with any keywords
	if _, ok := lexer.Keywords[s]; !ok {
		// only trigger the error if the star is not used
		if !star {
			return fmt.Errorf("%s is a reserved keyword. use * before a keyword to use it as a field name", s)
		}
	}
	return nil
}
