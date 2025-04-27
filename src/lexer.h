#ifndef LEXER_H
#define LEXER_H

#include <stdlib.h>

#define LEXEOF 0

enum token_type {
	TokenEntity,
	/* identifiers */
	TokenTypeInt,
	TokenTypeFloat,
	TokenTypeText,
	TokenTypeBool,
	TokenTypeUUID,
	TokenTypeRoutes,
	TokenTypeTimestamp,
	/* keywords */
	TokenAlter,
	TokenRef,
	TokenSelf,
	TokenEnd,
	TokenEndpoint,
	// symbols
	TokenArrow,
	TokenEnumOpen,
	TokenEnumClose,
	TokenListOpen,
	TokenListClose,
	TokenConsOpen,
	TokenConsClose,
	TokenComment,
	TokenNewline,
	TokenDot,
	TokenColon,
	// values
	TokenIdent,
	TokenString,
	TokenDigits,
	TokenDigitsFloat,
	// http
	TokenGet,
	TokenPost,
	TokenPut,
	TokenDelete,
	// constraints
	TokenConstraintAutoIncrement,
	TokenConstraintUnique,
	TokenConstraintForeignKey,
	TokenConstraintPrimaryKey,
	TokenConstraintNotNull,
	TokenConstraintDefault,
	TokenEOF,
	TokenUnknown,
};

typedef enum token_type token_type;
typedef struct lexer_t lexer_t;
typedef struct token_t token_t;

struct lexer_t {
	char *input;
	int position;
	int readPosition;
	size_t inputLen;
	char ch;
};

struct token_t {
	enum token_type type;
	char *literal;
	int lineNum;
};

/* lexer functions */
token_t *matchOrUnknown(lexer_t *, char,  token_type, token_type);
token_t *nextToken(lexer_t *);
token_t *newToken(token_type, char);
lexer_t *newLexer(const char *); 

char *collectEndpointVal(lexer_t *);
char *readString(lexer_t *);
char *readNumber(lexer_t *, int *);

void freeToken(token_t *);
void readChar(lexer_t *);
void skipComment(lexer_t *);
void skipWhitespace(lexer_t *);
void skipWhitespace(lexer_t *);
void skipWhitespace(lexer_t *);
void skipWhitespace(lexer_t *);

#endif
