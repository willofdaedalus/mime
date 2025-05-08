// Copyright (c) 2016-2017 Thorsten Ball
// Licensed under the MIT License. See LICENSE for details.
package lexer

import (
	"fmt"
	"unicode"
)

type TokenType int

// lexer struct
type Lexer struct {
	input        string
	position     int  // current position in input
	readPosition int  // next position to read
	ch           byte // current character being examined
}

// token struct
type Token struct {
	FileName string
	Type     TokenType
	Literal  string
	LineNum  int
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}

func (l Lexer) RenderTokens() {
	tok := l.NextToken()

	for tok.Type != TokenEOF {
		fmt.Printf("%v\n", tok)
		tok = l.NextToken()
	}
}

func (l *Lexer) NextToken() Token {
	var tok Token

	// this might backfire especially if we need structure
	l.skipWhitespace()

	switch l.ch {
	case '.':
		tok = newToken(TokenDot, l.ch)
	case '{':
		tok = newToken(TokenConsOpen, l.ch)
	case '}':
		tok = newToken(TokenConsClose, l.ch)
	case '#':
		tok = newToken(TokenComment, l.ch)
		l.skipComment()
	case '[':
		tok = newToken(TokenListOpen, l.ch)
	case ']':
		tok = newToken(TokenListClose, l.ch)
	case '(':
		tok = newToken(TokenEnumOpen, l.ch)
	case ')':
		tok = newToken(TokenEnumClose, l.ch)
	case ':':
		tok = newToken(TokenColon, l.ch)
	case '@':
		tok = newToken(TokenAmpersand, l.ch)
	case '*':
		tok = newToken(TokenStar, l.ch)
	case '-':
		tok = l.matchOrUnknown('>', TokenArrow, TokenUnknown)
	case '/':
		if unicode.IsLetter(rune(l.peekChar())) {
			tok.Literal = l.collectEndpointStr()
			tok.Type = TokenEndpoint
			return tok
		}
		tok = newToken(TokenUnknown, l.ch)
	case '\n':
		tok = newToken(TokenNewline, l.ch)
	case '"':
		tok.Type = TokenString
		tok.Literal = l.readString()
		if tok.Literal == "UNKNOWN" {
			tok.Type = TokenUnknown
		}
		return tok
	case 0:
		tok.Literal = ""
		tok.Type = TokenEOF
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = lookUpIdent(tok.Literal)
			return tok
		} else if unicode.IsDigit(rune(l.ch)) {
			tok.Type = TokenDigits
			v, f := l.readNumber()
			if f {
				tok.Type = TokenDigitsFloat
			}
			tok.Literal = v
			return tok
		} else {
			tok = newToken(TokenUnknown, l.ch)
		}
	}

	l.readChar()
	return tok
}

func (l *Lexer) skipComment() {
	for l.ch != '\n' {
		l.readChar()
	}
}

func (l *Lexer) skipWhitespace() {
	for unicode.IsSpace(rune(l.ch)) {
		l.readChar()
	}
}

func (l *Lexer) matchOrUnknown(expected byte, multiType, singleType TokenType) Token {
	if l.peekChar() == expected {
		ch := l.ch
		l.readChar()
		return Token{Type: multiType, Literal: string(ch) + string(l.ch)}
	}
	return newToken(singleType, l.ch)
}

func (l *Lexer) readNumber() (string, bool) {
	isFloat := false
	start := l.position

	// read integer part
	for unicode.IsDigit(rune(l.ch)) {
		l.readChar()
	}

	// handle decimal point
	if l.ch == '.' {
		// set the float flag
		isFloat = true
		l.readChar() // consume the decimal point

		// if there's at least one digit after the decimal, read the fractional part
		if unicode.IsDigit(rune(l.ch)) {
			for unicode.IsDigit(rune(l.ch)) {
				l.readChar()
			}
			return l.input[start:l.position], isFloat // return full float number
		}
	}

	return l.input[start:l.position], isFloat // return integer
}

func (l *Lexer) collectEndpointStr() string {
	start := l.position
	l.readChar()

	for l.ch != ' ' && l.ch != '\n' {
		l.readChar()
	}
	return l.input[start:l.position]
}

func (l *Lexer) readString() string {
	start := l.position // start position (including the opening quote)
	l.readChar()        // consume opening quote

	for l.ch != '"' && l.ch != '\n' && l.ch != 0 {
		// handle escape sequences (\" or \\ or \n, etc.)
		if l.ch == '\\' {
			l.readChar() // skip past the backslash to include the escaped char
		}
		l.readChar()
	}

	// if we hit a newline or eof before finding a closing quote, it's unknown
	if l.ch != '"' {
		return "UNKNOWN"
	}

	l.readChar() // consume the closing quote
	return l.input[start+1 : l.position-1]
}

func (l *Lexer) readIdentifier() string {
	start := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[start:l.position]
}

func isLetter(ch byte) bool {
	return unicode.IsLetter(rune(ch)) || ch == '_'
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}

	return l.input[l.readPosition]
}

func newToken(tokType TokenType, ch byte) Token {
	return Token{Type: tokType, Literal: string(ch)}
}
