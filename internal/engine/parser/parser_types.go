package parser

type tokenType int

// token types for lexer
const (
	TK_ENTITY     tokenType = iota + 1
	TK_NUMBER               // number
	TK_TEXT                 // text
	TK_TIMESTAMP            // timestamp
	TK_UUID                 // uuid
	TK_ROUTES               // routes
	TK_HTTP_VERB            // GET, POST, DELETE, PUT
	TK_ALTER                // alter
	TK_REF                  // ref
	TK_SELF                 // self
	TK_BOOL                 // bool
	TK_ENDPOINT             // /employees/:id
	TK_ARROW                // ->
	TK_OPEN_ANGLE           // <>
	TK_LBRACKET             // [
	TK_RBRACKET             // ]
	TK_COLON                // :
	TK_LBRACE               // {
	TK_RBRACE               // }
	TK_COMMA                // ,
	TK_NEWLINE              // \n
	TK_DOT                  // .
	TK_IDENT                // identifiers like "id", "student", "payload"
	TK_STRING               // string literals (e.g., `"male"`, `"female"`)
	TK_DIGITS               // 123, 45.6
	TK_EOF
	TK_UNKNOWN
)

// token struct
type token struct {
	Type    tokenType
	Literal string
}

// lexer struct
type lexer struct {
	input        string
	position     int  // current position in input
	readPosition int  // next position to read
	ch           byte // current character being examined
}

var keywords = map[string]tokenType{
	"entity":    TK_ENTITY,
	"number":    TK_NUMBER,
	"text":      TK_TEXT,
	"timestamp": TK_TIMESTAMP,
	"uuid":      TK_UUID,
	"routes":    TK_ROUTES,
	"alter":     TK_ALTER,
	"ref":       TK_REF,
	"self":      TK_SELF,
	"GET":       TK_HTTP_VERB,
	"POST":      TK_HTTP_VERB,
	"DELETE":    TK_HTTP_VERB,
	"PUT":       TK_HTTP_VERB,
}

func lookUpIdent(ident string) tokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}

	return TK_IDENT
}

// tokenTypeToString converts the tokenType to its string representation
func (t tokenType) String() string {
	switch t {
	case TK_ENTITY:
		return "TK_ENTITY"
	case TK_NUMBER:
		return "TK_NUMBER"
	case TK_TEXT:
		return "TK_TEXT"
	case TK_TIMESTAMP:
		return "TK_TIMESTAMP"
	case TK_UUID:
		return "TK_UUID"
	case TK_ROUTES:
		return "TK_ROUTES"
	case TK_HTTP_VERB:
		return "TK_HTTP_VERB"
	case TK_ALTER:
		return "TK_ALTER"
	case TK_REF:
		return "TK_REF"
	case TK_SELF:
		return "TK_SELF"
	case TK_BOOL:
		return "TK_BOOL"
	case TK_ENDPOINT:
		return "TK_ENDPOINT"
	case TK_ARROW:
		return "TK_ARROW"
	case TK_OPEN_ANGLE:
		return "TK_OPEN_ANGLE"
	case TK_LBRACKET:
		return "TK_LBRACKET"
	case TK_RBRACKET:
		return "TK_RBRACKET"
	case TK_COLON:
		return "TK_COLON"
	case TK_LBRACE:
		return "TK_LBRACE"
	case TK_RBRACE:
		return "TK_RBRACE"
	case TK_COMMA:
		return "TK_COMMA"
	case TK_NEWLINE:
		return "TK_NEWLINE"
	case TK_DOT:
		return "TK_DOT"
	case TK_IDENT:
		return "TK_IDENT"
	case TK_STRING:
		return "TK_STRING"
	case TK_DIGITS:
		return "TK_DIGITS"
	case TK_EOF:
		return "TK_EOF"
	case TK_UNKNOWN:
		return "TK_UNKNOWN"
	default:
		return "Unknown token type"
	}
}
