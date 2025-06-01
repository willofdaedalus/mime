package parser

import (
	"fmt"

	"willofdaedalus/mime/internal/engine/lexer"
	"willofdaedalus/mime/internal/engine/types"
)

func handleEnum(p *Parser) (node, error) {
	if p.curType() != lexer.TokenEnum {
		return nil, fmt.Errorf("expected enum, got %s", p.curToken.Type)
	}
	p.advanceToken() // consume enum

	// enum name
	if !expectTokOf(p.curToken, lexer.TokenIdent) {
		return (node)(nil), fmt.Errorf("expected enum name, got %s", p.curToken.Type)
	}

	enumNode := &types.EnumNode{
		Members: make([]string, 0, 10),
		Name:    p.curToken.Literal,
	}
	p.advanceToken() // consume "name"

	// check for arrow token
	if !expectTokOf(p.curToken, lexer.TokenArrow) {
		return (node)(nil), fmt.Errorf("expected -> after enum name, got %s", p.curToken.Type)
	}
	p.advanceToken() // consume '->'

	for p.curToken.Type != lexer.TokenEnd {
		if p.curToken.Type == lexer.TokenComment {
			p.advanceToken() // skip comments
			continue
		}

		// unexpected end to file with no end keyword
		if p.curToken.Type == lexer.TokenEOF {
			return (node)(nil), fmt.Errorf("unexpected end to file")
		}

		if p.curToken.Type != lexer.TokenIdent {
			return (node)(nil), fmt.Errorf("expected enum member got %s", p.curToken.Literal)
		}

		v := p.curToken.Literal
		// THIS IS WORK FOR THE AST NOT THE PARSER
		// validate and make sure there are no duplicates
		// err := validateEnumValue(v)
		// if err != nil {
		// 	p.addError(ParserLogError, err.Error())
		// 	p.advanceToken()
		// 	continue
		// }
		// if slices.Contains(enumNode.Members, v) {
		// 	p.addError(ParserLogError, fmt.Sprintf("duplicate enum member %s", v))
		// 	p.advanceToken()
		// 	continue
		// }

		enumNode.Members = append(enumNode.Members, v)
		p.advanceToken()
	}

	// ANOTHER ONE THAT THE AST SHOULD HANDLE
	// if len(enumNode.Members) == 0 {
	// 	return (node)(nil), fmt.Errorf("empty enums are not allowed")
	// }

	return enumNode, nil
}
