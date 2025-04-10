package lexer

// token types for lexer
const (
	// identifiers
	TokenEntity        TokenType = iota + 1 // entity
	TokenTypeInt                            // int
	TokenTypeFloat                          // float
	TokenTypeText                           // text
	TokenTypeTimestamp                      // timestamp
	TokenTypeUuid                           // uuid
	TokenTypeRoutes                         // routes
	// keywords
	TokenAlter    // alter
	TokenRef      // ref
	TokenSelf     // self
	TokenBool     // bool
	TokenEnd      // end
	TokenEndpoint // /employees/:id
	// symbols
	TokenArrow     // ->
	TokenEnumOpen  // (
	TokenEnumClose // )
	TokenListOpen  // [
	TokenListClose // ]
	TokenConsOpen  // {
	TokenConsClose // }
	TokenComment   // #
	TokenNewline   // \n
	TokenDot       // .
	TokenColon     // :
	// values
	TokenIdent       // identifiers like id, student, payload
	TokenString      // string literals (e.g., `"male"`, `"female"`)
	TokenDigits      // 123
	TokenDigitsFloat // 45.6
	// http
	TokenGet    // GET
	TokenPost   // POST
	TokenPut    // PUT
	TokenDelete // DELETE
	// constraints
	TokenConstraintAutoIncrement // increment
	TokenConstraintUnique        // unique
	TokenConstraintForeignKey    // fk
	TokenConstraintPrimaryKey    // primary
	TokenConstraintNotNull       // required
	TokenConstraintDefault       // default
	TokenEOF
	TokenUnknown
)

var keywords = map[string]TokenType{
	// normal keywords
	"entity":    TokenEntity,
	"float":     TokenTypeFloat,
	"int":       TokenTypeInt,
	"text":      TokenTypeText,
	"timestamp": TokenTypeTimestamp,
	"uuid":      TokenTypeUuid,
	"routes":    TokenTypeRoutes,
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
	"default":   TokenConstraintDefault,
	"fk":        TokenConstraintForeignKey,
	"primary":   TokenConstraintPrimaryKey,
	"required":  TokenConstraintNotNull,
}

var AllDataTypes = map[TokenType]struct{}{
	TokenTypeInt:       {},
	TokenTypeTimestamp: {},
	TokenTypeText:      {},
	TokenTypeFloat:     {},
	TokenTypeUuid:      {},
}

var allConstraints = map[TokenType]struct{}{
	TokenConstraintAutoIncrement: {},
	TokenConstraintUnique:        {},
	TokenConstraintForeignKey:    {},
	TokenConstraintPrimaryKey:    {},
	TokenConstraintNotNull:       {},
}

var AnnotationOpens = map[TokenType]struct{}{
	TokenEnumOpen: {},
	TokenListOpen: {},
	TokenConsOpen: {},
}

func IsValidMemberOf(tt TokenType, list map[TokenType]struct{}) bool {
	_, ok := list[tt]
	return ok
}

func lookUpIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}

	return TokenIdent
}

// String converts the tokenType to its string representation
func (t TokenType) String() string {
	switch t {
	case TokenEntity:
		return "entity"
	case TokenTypeFloat:
		return "float"
	case TokenTypeInt:
		return "int"
	case TokenTypeText:
		return "text"
	case TokenTypeTimestamp:
		return "timestamp"
	case TokenTypeUuid:
		return "uuid"
	case TokenTypeRoutes:
		return "routes"
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
