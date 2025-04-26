#include "lexer.h"
#include <stdlib.h>
#include <stdio.h>
#include <string.h>

/* newLexer takes some input as text and returns
* a new lexer object which we can use to tokenize
*/
struct lexer_t *newLexer(const char *input) {
	size_t inputLen = 0;
	/* use calloc to allocate and initialize the struct in one shot */
	struct lexer_t *lexer = calloc(1, sizeof(struct lexer_t));
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

	return lexer;
}

/* readChar basically initializes the first character and poistions
* in the lexer to keep track of characters for tokenizing
*/
void readChar(struct lexer_t *lexer) {
	if (lexer->readPosition >= strlen(lexer->input)) {
		lexer->ch = LEXEOF;
	} else {
		lexer->ch = lexer->input[lexer->readPosition];
	}
	lexer->position = lexer->readPosition;
	lexer->readPosition += 1;
}
