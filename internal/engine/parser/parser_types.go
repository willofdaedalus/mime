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
	TK_BOOL                 // bool
	TK_ARROW                // ->
	TK_OPEN_ANGLE           // <>
	TK_SLASH                // /
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
	case TK_SLASH:
		return "TK_SLASH"
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
	case TK_EOF:
		return "TK_EOF"
	case TK_UNKNOWN:
		return "TK_UNKNOWN"
	default:
		return "Unknown token type"
	}
}
