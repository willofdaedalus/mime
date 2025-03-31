// Copyright (c) 2016-2017 Thorsten Ball
// Licensed under the MIT License. See LICENSE for details.
package parser

import (
	"unicode"
)

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
		tok = l.matchOrUnknown('>', TK_ARROW, TK_UNKNOWN)
	case '<':
		tok = l.matchOrUnknown('>', TK_OPEN_ANGLE, TK_UNKNOWN)
	case '/':
		tok = newToken(TK_SLASH, l.ch)
	case '"':
		tok.Type = TK_STRING
		tok.Literal = l.readString()
		if tok.Literal == "UNKNOWN" {
			tok.Type = TK_UNKNOWN
		}
		return tok
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

func (l *lexer) matchOrUnknown(expected byte, multiType, singleType tokenType) token {
	if l.peekChar() == expected {
		ch := l.ch
		l.readChar()
		return token{Type: multiType, Literal: string(ch) + string(l.ch)}
	}
	return newToken(singleType, l.ch)
}

func (l *lexer) readNumber() string {
	start := l.position

	// read integer part
	for unicode.IsDigit(rune(l.ch)) {
		l.readChar()
	}

	// handle decimal point
	if l.ch == '.' {
		l.readChar() // consume the decimal point

		// if there's at least one digit after the decimal, read the fractional part
		if unicode.IsDigit(rune(l.ch)) {
			for unicode.IsDigit(rune(l.ch)) {
				l.readChar()
			}
			return l.input[start:l.position] // return full float number
		}
	}

	return l.input[start:l.position] // return integer
}

func (l *lexer) readString() string {
	start := l.position // Start position (including the opening quote)
	l.readChar()        // Consume opening quote

	for l.ch != '"' && l.ch != '\n' && l.ch != 0 {
		// Handle escape sequences (\" or \\ or \n, etc.)
		if l.ch == '\\' {
			l.readChar() // Skip past the backslash to include the escaped char
		}
		l.readChar()
	}

	// If we hit a newline or EOF before finding a closing quote, it's unknown
	if l.ch != '"' {
		return "UNKNOWN"
	}

	l.readChar() // Consume the closing quote
	return l.input[start:l.position]
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

func (l *lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}

	return l.input[l.readPosition]
}

func newToken(tokType tokenType, ch byte) token {
	return token{Type: tokType, Literal: string(ch)}
}
