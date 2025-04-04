package lexer

// token types for lexer
const (
	TokenEntity                  tokenType = iota + 1 // entity
	TokenInt                                          // int
	TokenFloat                                        // float
	TokenText                                         // text
	TokenTimestamp                                    // timestamp
	TokenUuid                                         // uuid
	TokenRoutes                                       // routes
	TokenAlter                                        // alter
	TokenRef                                          // ref
	TokenSelf                                         // self
	TokenBool                                         // bool
	TokenEnd                                          // end
	TokenEndpoint                                     // /employees/:id
	TokenArrow                                        // ->
	TokenEnumOpen                                     // (
	TokenEnumClose                                    // )
	TokenListOpen                                     // [
	TokenListClose                                    // ]
	TokenConsOpen                                     // {
	TokenConsClose                                    // }
	TokenNewline                                      // \n
	TokenDot                                          // .
	TokenIdent                                        // identifiers like "id", "student", "payload"
	TokenString                                       // string literals (e.g., `"male"`, `"female"`)
	TokenDigits                                       // 123, 45.6
	TokenGet                                          // GET
	TokenPost                                         // POST
	TokenPut                                          // PUT
	TokenDelete                                       // DELETE
	TokenComment                                      // #
	TokenConstraintAutoIncrement                      // increment
	TokenConstraintUnique                             // unique
	TokenConstraintForeignKey                         // fk
	TokenConstraintPrimaryKey                         // primary
	TokenConstraintNotNull                            // required
	TokenEOF
	TokenUnknown
)

var keywords = map[string]tokenType{
	// normal keywords
	"entity":    TokenEntity,
	"float":     TokenFloat,
	"int":       TokenInt,
	"text":      TokenText,
	"timestamp": TokenTimestamp,
	"uuid":      TokenUuid,
	"routes":    TokenRoutes,
	"alter":     TokenAlter,
	"ref":       TokenRef,
	"self":      TokenSelf,
	"end":       TokenEnd,
	// http verbs
	"GET":    TokenGet,
	"POST":   TokenPost,
	"DELETE": TokenDelete,
	"PUT":    TokenPut,
	// constraints
	"increment": TokenConstraintAutoIncrement,
	"unique":    TokenConstraintUnique,
	"fk":        TokenConstraintForeignKey,
	"primary":   TokenConstraintPrimaryKey,
	"required":  TokenConstraintNotNull,
}

var AllDataTypes = map[tokenType]struct{}{
	TokenInt:       {},
	TokenTimestamp: {},
	TokenText:      {},
	TokenFloat:     {},
	TokenUuid:      {},
}

var allConstraints = map[tokenType]struct{}{
	TokenConstraintAutoIncrement: {},
	TokenConstraintUnique:        {},
	TokenConstraintForeignKey:    {},
	TokenConstraintPrimaryKey:    {},
	TokenConstraintNotNull:       {},
}

var AnnotationOpens = map[tokenType]struct{}{
	TokenEnumOpen: {},
	TokenListOpen: {},
	TokenConsOpen: {},
}

func IsValidMemberOf(tt tokenType, list map[tokenType]struct{}) bool {
	_, ok := list[tt]
	return ok
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
		return "entity"
	case TokenFloat:
		return "float"
	case TokenInt:
		return "int"
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
	case TokenListOpen:
		return "TokenLBRACKET"
	case TokenListClose:
		return "TokenRBRACKET"
	case TokenConsOpen:
		return "TokenLBRACE"
	case TokenConsClose:
		return "TokenRBRACE"
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
