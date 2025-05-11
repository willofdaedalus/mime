package parser

import (
	"fmt"
	"slices"

	"willofdaedalus/mime/internal/engine/lexer"
	"willofdaedalus/mime/internal/engine/types"
)

func handleEnum(p *Parser) node {
	defer p.resetContext()

	if !expectTokOf(p.curToken, lexer.TokenEnum) {
		p.addError(ParserLogError,
			fmt.Sprintf("expected enum, got %s", p.curToken.Type))
	}
	p.advanceToken() // consume enum

	if !expectTokOf(p.curToken, lexer.TokenIdent) {
		p.addError(ParserLogError,
			fmt.Sprintf("expected enum name, got %s", p.curToken.Type))
		fmt.Println("expected num name")
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
	}
	p.advanceToken() // consume '->'

	for p.curToken.Type != lexer.TokenEnd {
		if p.curToken.Type == lexer.TokenComment {
			p.advanceToken() // skip comments
			continue
		}

		// unexpected end to file with no end keyword
		if p.curToken.Type == lexer.TokenEOF {
			return (*types.EnumNode)(nil)
		}

		if p.curToken.Type != lexer.TokenIdent {
			p.addError(ParserLogError,
				fmt.Sprintf("expected enum member got %s", p.curToken.Type))
			skipToTok(p, lexer.TokenIdent)
			// p.advanceToken()
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

	if len(enumNode.Members) == 0 {
		// this won't trigger a p.invalidParsing but will generate a warning
		p.addError(ParserLogWarning, fmt.Sprintf("enum %s is declared and might not works as expected",
			enumNode.Name))
	}

	if p.invalidParsing {
		// this passes the test instead of the usual nil
		return (*types.EnumNode)(nil)
	}

	return enumNode
}

func validateEnumValue(s string) error {
	// check that it doesn't conflict with any keywords
	if _, ok := lexer.Keywords[s]; ok {
		fmt.Println("starts with keyword")
		return fmt.Errorf("%s is a reserved keyword", s)
	}
	return nil
}
