#ifndef LEXER_H
#define LEXER_H

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

struct lexer_t {
	char *input;
	int position;
	int readPosition;
	char ch;
};

struct token_t {
	enum token_type type;
	char *literal;
	int lineNum;
};

/* lexer functions */
struct lexer_t *newLexer(const char *); 
void readChar(struct lexer_t *);
void nextToken(struct lexer_t *);
void skipComment(struct lexer_t *);
void skipWhitespace(struct lexer_t *);
char *readNumber(struct lexer_t *);
void skipWhitespace(struct lexer_t *);
void skipWhitespace(struct lexer_t *);
void skipWhitespace(struct lexer_t *);

#endif
