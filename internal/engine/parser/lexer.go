// Copyright (c) 2016-2017 Thorsten Ball
// Licensed under the MIT License. See LICENSE for details.
package parser

import "unicode"

func New(input string) *lexer {
	l := &lexer{input: input}
	l.readChar()
	return l
}

func (l *lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}

func (l *lexer) NextToken() token {
	var tok token

	// this might backfire especially if we need structure
	l.skipWhitespace()

	switch l.ch {
	case ':':
		tok = newToken(TK_COLON, l.ch)
	case ',':
		tok = newToken(TK_COMMA, l.ch)
	case '{':
		tok = newToken(TK_LBRACE, l.ch)
	case '}':
		tok = newToken(TK_RBRACE, l.ch)
	case '[':
		tok = newToken(TK_LBRACKET, l.ch)
	case ']':
		tok = newToken(TK_RBRACKET, l.ch)
	case '-':
		tok = newToken(TK_DASH, l.ch)
	case '>':
		tok = newToken(TK_RANGLE, l.ch)
	case '<':
		tok = newToken(TK_LANGLE, l.ch)
	case '/':
		tok = newToken(TK_SLASH, l.ch)
	case 0:
		tok.Literal = ""
		tok.Type = TK_EOF
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = lookUpIdent(tok.Literal)
			return tok
		} else if unicode.IsDigit(rune(l.ch)) {
			tok.Type = TK_DIGITS
			tok.Literal = l.readNumber()
			return tok
		} else {
			tok = newToken(TK_UNKNOWN, l.ch)
		}
	}

	l.readChar()
	return tok
}

func (l *lexer) skipWhitespace() {
	for unicode.IsSpace(rune(l.ch)) {
		l.readChar()
	}
}

func (l *lexer) readNumber() string {
	start := l.position
	// read the integer part (digits before the decimal point)
	for unicode.IsDigit(rune(l.ch)) {
		l.readChar()
	}

	// if the next character is a decimal point, read the fractional part (digits after the decimal point)
	if l.ch == '.' {
		l.readChar() // consume the decimal point
		// ensure we have digits after the decimal point, otherwise it's not a valid float
		if unicode.IsDigit(rune(l.ch)) {
			// read the fractional part (digits after the decimal point)
			for unicode.IsDigit(rune(l.ch)) {
				l.readChar()
			}
		} else {
			// if there's no digit after the decimal, revert the read
			l.position = start
			l.readPosition = start
			l.ch = l.input[l.position]
			return l.input[start:l.position] // return as integer, no decimal part
		}
	}

	return l.input[start:l.position] // return the whole number (int or float)
}

func (l *lexer) readIdentifier() string {
	start := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[start:l.position]
}

func isLetter(ch byte) bool {
	return unicode.IsLetter(rune(ch)) || ch == '_'
}

func newToken(tokType tokenType, ch byte) token {
	return token{Type: tokType, Literal: string(ch)}
}
