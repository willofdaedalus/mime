package parser

import (
	"fmt"
	"slices"
	"unicode"

	"willofdaedalus/mime/internal/engine/lexer"
	"willofdaedalus/mime/internal/engine/types"
)

func handleEnum(p *Parser) node {
	defer p.resetContext()

	if !expectTokOf(p.curToken, lexer.TokenEnum) {
		p.addError(ParserLogError,
			fmt.Sprintf("expected enum, got %s", p.curToken.Type))
		return nil
	}
	p.advanceToken() // consume entity

	if !expectTokOf(p.curToken, lexer.TokenIdent) {
		p.addError(ParserLogError,
			fmt.Sprintf("expected enum name, got %s", p.curToken.Type))
		fmt.Println("expected num name")
		return nil
	}

	enumNode := &types.EnumNode{
		Members: make([]string, 0, 10),
		Name:    p.curToken.Literal,
	}
	p.advanceToken() // consume "name"

	// check for arrow token
	if !expectTokOf(p.curToken, lexer.TokenArrow) {
		p.addError(ParserLogError,
			fmt.Sprintf("expected -> after enum name, got %s", p.curToken.Type))
		fmt.Println("expected -> after enum name")
		return nil
	}
	p.advanceToken() // consume '->'

	for p.curToken.Type != lexer.TokenEnd && p.curToken.Type != lexer.TokenEOF {
		if p.curToken.Type == lexer.TokenNewline || p.curToken.Type == lexer.TokenComment {
			p.advanceToken() // skip newlines and comments
			continue
		}

		if p.curToken.Type != lexer.TokenIdent {
			p.addError(ParserLogError,
				fmt.Sprintf("expected enum member got %s", p.curToken.Type))
			p.advanceToken()
			continue
		}

		v := p.curToken.Literal
		// validate and make sure there are no duplicates
		err := validateEnumValue(v)
		if err != nil {
			p.addError(ParserLogError, err.Error())
			p.advanceToken()
			continue
		}
		if slices.Contains(enumNode.Members, v) {
			p.addError(ParserLogError, fmt.Sprintf("duplicate enum member %s", v))
			p.advanceToken()
			continue
		}

		enumNode.Members = append(enumNode.Members, v)
		p.advanceToken()
	}

	if p.invalidParsing || len(enumNode.Members) == 0 {
		// this passes the test instead of the usual nil
		return (*types.EnumNode)(nil)
	}

	return enumNode
}

func validateEnumValue(s string) error {
	// check that it's not starting with a number
	if unicode.IsDigit(rune(s[0])) {
		fmt.Println("starts with number")
		return fmt.Errorf("%s starts with a number", s)
	}
	// next check that it doesn't start with symbol
	// if it does make sure to it's not
	if unicode.IsSymbol(rune(s[0])) && s[0] != '_' {
		fmt.Println("starts with symbol")
		return fmt.Errorf("%s starts with a symbol", s)
	}
	// next check that it doesn't conflict with any keywords
	if _, ok := lexer.Keywords[s]; ok {
		fmt.Println("starts with keyword")
		return fmt.Errorf("%s is a reserved keyword", s)
	}
	return nil
}
