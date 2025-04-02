package parser

type tokenType int

// token types for lexer
const (
	TokenEntity    tokenType = iota + 1 // entity
	TokenNumber                         // number
	TokenText                           // text
	TokenTimestamp                      // timestamp
	TokenUuid                           // uuid
	TokenRoutes                         // routes
	TokenAlter                          // alter
	TokenRef                            // ref
	TokenSelf                           // self
	TokenBool                           // bool
	TokenEndpoint                       // /employees/:id
	TokenArrow                          // ->
	TokenOpenAngle                      // <>
	TokenLBracket                       // [
	TokenRBracket                       // ]
	TokenColon                          // :
	TokenLBrace                         // {
	TokenRBrace                         // }
	TokenComma                          // ,
	TokenNewline                        // \n
	TokenDot                            // .
	TokenIdent                          // identifiers like "id", "student", "payload"
	TokenString                         // string literals (e.g., `"male"`, `"female"`)
	TokenDigits                         // 123, 45.6
	TokenGet                            // GET
	TokenPost                           // POST
	TokenPut                            // PUT
	TokenDelete                         // DELETE
	TokenEOF
	TokenUnknown
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
	"entity":    TokenEntity,
	"number":    TokenNumber,
	"text":      TokenText,
	"timestamp": TokenTimestamp,
	"uuid":      TokenUuid,
	"routes":    TokenRoutes,
	"alter":     TokenAlter,
	"ref":       TokenRef,
	"self":      TokenSelf,
	"GET":       TokenGet,
	"POST":      TokenPost,
	"DELETE":    TokenDelete,
	"PUT":       TokenPut,
}

func lookUpIdent(ident string) tokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}

	return TokenIdent
}

// tokenTypeToString converts the tokenType to its string representation
func (t tokenType) String() string {
	switch t {
	case TokenEntity:
		return "TokenENTITY"
	case TokenNumber:
		return "TokenNUMBER"
	case TokenText:
		return "TokenTEXT"
	case TokenTimestamp:
		return "TokenTIMESTAMP"
	case TokenUuid:
		return "TokenUUID"
	case TokenRoutes:
		return "TokenROUTES"
	case TokenAlter:
		return "TokenALTER"
	case TokenRef:
		return "TokenREF"
	case TokenSelf:
		return "TokenSELF"
	case TokenBool:
		return "TokenBOOL"
	case TokenEndpoint:
		return "TokenENDPOINT"
	case TokenArrow:
		return "TokenARROW"
	case TokenOpenAngle:
		return "TokenOPEN_ANGLE"
	case TokenLBracket:
		return "TokenLBRACKET"
	case TokenRBracket:
		return "TokenRBRACKET"
	case TokenColon:
		return "TokenCOLON"
	case TokenLBrace:
		return "TokenLBRACE"
	case TokenRBrace:
		return "TokenRBRACE"
	case TokenComma:
		return "TokenCOMMA"
	case TokenNewline:
		return "TokenNEWLINE"
	case TokenDot:
		return "TokenDOT"
	case TokenIdent:
		return "TokenIDENT"
	case TokenString:
		return "TokenSTRING"
	case TokenDigits:
		return "TokenDIGITS"
	case TokenEOF:
		return "TokenEOF"
	case TokenUnknown:
		return "TokenUNKNOWN"
	default:
		return "Unknown token type"
	}
}
