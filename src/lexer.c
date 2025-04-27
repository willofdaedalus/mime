#include "lexer.h"
#include <stdlib.h>
#include <stdio.h>
#include <string.h>
#include <ctype.h>

/* newLexer takes some input as text and returns
* a new lexer object which we can use to tokenize
*/
lexer_t *newLexer(const char *input) {
	size_t inputLen = 0;
	/* use calloc to allocate and initialize the in one shot */
	lexer_t *lexer = calloc(1, sizeof(lexer_t));
	if (lexer == NULL) {
		fprintf(stderr, "failed to malloc for lexer");
		return NULL;
	}

	if (input == NULL || (inputLen = strlen(input)) == 0) {
		fprintf(stderr, "empty input passed to newLexer");
		free(lexer);
		return NULL;
	}

	if ((lexer->input = strdup(input)) == NULL) {
		fprintf(stderr, "allocation for lexer input failed");
		free(lexer);
		return NULL;
	}

	readChar(lexer);
	lexer->inputLen = strlen(input);

	return lexer;
}

/* readChar basically initializes the first character and poistions
* in the lexer to keep track of characters for tokenizing
*/
void readChar(lexer_t *lexer) {
	if (lexer->readPosition >= strlen(lexer->input)) {
		lexer->ch = LEXEOF;
	} else {
		lexer->ch = lexer->input[lexer->readPosition];
	}
	lexer->position = lexer->readPosition;
	lexer->readPosition += 1;
}

/* peekChar reads the very next character in the input */
char peekChar(lexer_t *lexer) {
	if (lexer->readPosition >= lexer->inputLen) {
		return 0;
	}

	return lexer->input[lexer->readPosition];
}

token_t *nextToken(lexer_t *lexer) {
	token_t *tok = {0};

	switch (lexer->ch) {
		case '.':
			if ((tok = newToken(TokenDot, lexer->ch)) == NULL) {
				return NULL;
			}
		break;

		case '{':
			if ((tok = newToken(TokenConsOpen, lexer->ch)) == NULL) {
				return NULL;
			}
		break;

		case '}':
			if ((tok = newToken(TokenConsClose, lexer->ch)) == NULL) {
				return NULL;
			}
		break;

		case '[':
			if ((tok = newToken(TokenListOpen, lexer->ch)) == NULL) {
				return NULL;
			}
		break;

		case ']':
			if ((tok = newToken(TokenListClose, lexer->ch)) == NULL) {
				return NULL;
			}
		break;
		case '(':
			if ((tok = newToken(TokenEnumOpen, lexer->ch)) == NULL) {
				return NULL;
			}
		break;
		case ')':
			if ((tok = newToken(TokenEnumClose, lexer->ch)) == NULL) {
				return NULL;
			}
		break;
		case ':':
			if ((tok = newToken(TokenColon, lexer->ch)) == NULL) {
				return NULL;
			}
		break;
		case '\n':
			if ((tok = newToken(TokenNewline, lexer->ch)) == NULL) {
				return NULL;
			}
		break;
		case '"':
			if ((tok = newToken(TokenUnknown, lexer->ch)) == NULL) {
				return NULL;
			}

			char *str;
			if ((str = readString(lexer)) == NULL) {
				freeToken(tok);
				return NULL;
			}

			tok->type = TokenString;
			tok->literal = str;

		break;
		case '-':
			if ((tok = matchOrUnknown(lexer, '>', TokenArrow, TokenUnknown)) == NULL) {
				return NULL;
			}
		break;
		case '/':
			/* set tok to a default value then alter it if
			 * the isalpha condition is true
			 */
			if ((tok = newToken(TokenUnknown, lexer->ch)) == NULL) {
				return NULL;
			}

			if (isalpha(peekChar(lexer))) {
				char *endpoint = collectEndpointVal(lexer);
				if (endpoint == NULL) {
					freeToken(tok);
					return NULL;
				}

				free(tok->literal);
				tok->literal = endpoint;
				tok->type = TokenEndpoint;
			} 
		break;
		case 0:
			if ((tok = newToken(TokenEOF, 0)) == NULL) {
				return NULL;
			}
		break;
			
	}

	readChar(lexer);
	return tok;
}

char *readString(lexer_t *lexer) {

}

/* collectEndpointVal basically reads through the input until the scanner
 * hits an EOF, empty space or a newline character. the result is then
 * backtracked and collected as an endpoint. whether it's valid or not
 * will be checked in the parser
 */
char *collectEndpointVal(lexer_t *lexer) {
	int start = lexer->position;
	char *endpointStr = NULL;

	readChar(lexer);
	while (lexer->ch != ' ' && lexer->ch != '\n' && lexer->ch != LEXEOF) {
		readChar(lexer);
	}

	int endpointLen = lexer->position - start;
	if ((endpointStr = malloc(sizeof(char) * (endpointLen + 1))) == NULL) {
		return NULL;
	}

	strncpy(endpointStr, lexer->input + start, endpointLen);
	endpointStr[endpointLen] = '\0';

	return endpointStr;
}

/* matchOrUnknown basically checks if the next char is next and returns
 * a token with the good type else using the bad type
 */
token_t *matchOrUnknown(lexer_t *lexer, char next,  token_type good,  token_type bad) {
	token_t *tok = newToken(bad, lexer->ch);
	if (tok == NULL) {
		return NULL;
	}

	if (peekChar(lexer) == next) {
		char ch = lexer->ch;
		/* advance the lexer position so now we're pointing to the
		 * next character and the previous one is stored in c
		 */
		readChar(lexer);

		char newLiteral[3] = { ch, lexer->ch, '\0' };
		/* free the previous literal before strdup-ing a new one */
		free(tok->literal);

		if ((tok->literal = strdup(newLiteral)) == NULL) {
			return NULL;
		}

		tok->type = good;
		return tok;
	}

	return tok;
}

/* newToken simply returns a token of type and literal */
token_t *newToken(token_type type, char ch) {
	token_t *tok = calloc(1, sizeof(token_t));
	if (tok == NULL) {
		return NULL;
	}

	char lit[2] = {ch, '\0'};
	tok->type = type;
	if ((tok->literal = strdup(lit))==NULL) {
		free(tok);
		return NULL;
	}

	return tok;
}

void freeToken(token_t *tok) {
	if (tok == NULL) {
		return;
	}
	free(tok->literal);
	free(tok);
}
