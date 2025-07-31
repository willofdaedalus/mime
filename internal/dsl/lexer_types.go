package dsl

type TokenType int

// lexer struct
type Lexer struct {
	input        string
	position     int  // current position in input
	readPosition int  // next position to read
	ch           byte // current character being examined
}

// token struct
type Token struct {
	FileName string
	Type     TokenType
	Literal  string
	LineNum  int
}

// token types for lexer
const (
	// identifiers
	TokenEntity TokenType = iota + 1 // entity

	TokenTypeInt       // int
	TokenTypeFloat     // float
	TokenTypeText      // text
	TokenTypeBool      // bool
	TokenTypeTimestamp // timestamp
	TokenTypeUuid      // uuid
	TokenTypeRoutes    // routes
	// keywords
	TokenApp      // app
	TokenAlter    // alter
	TokenEnum     // enum
	TokenRef      // ref
	TokenSelf     // self
	TokenEnd      // end
	TokenEndpoint // /employees/:id
	// symbols
	TokenArrow     // ->
	TokenAttrOpen  // [
	TokenAttrClose // ]
	// TokenListOpen  // [
	// TokenListClose // ]
	TokenConsOpen  // {
	TokenConsClose // }
	TokenComment   // #
	TokenAtSymbol  // @
	TokenAmpersand // &
	TokenStar      // * for using a keyword as a field or enum member name
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
	TokenAttrIncrement // increment
	TokenAttrUnique    // unique
	TokenAttrNotNull   // required
	TokenAttrDefault   // default
	TokenAttrHidden    // hash
	TokenAttrCheck
	TokenAttrHash

	TokenEOF
	TokenUnknown
)

var keywords = map[string]TokenType{
	// normal keywords
	"entity":    TokenEntity,
	"enum":      TokenEnum,
	"float":     TokenTypeFloat,
	"int":       TokenTypeInt,
	"text":      TokenTypeText,
	"timestamp": TokenTypeTimestamp,
	"uuid":      TokenTypeUuid,
	"bool":      TokenTypeBool,
	"routes":    TokenTypeRoutes,
	"alter":     TokenAlter,
	"ref":       TokenRef,
	"self":      TokenSelf,
	"end":       TokenEnd,
	"app":       TokenApp,
	// http verbs
	"GET":    TokenGet,
	"POST":   TokenPost,
	"DELETE": TokenDelete,
	"PUT":    TokenPut,
	// constraints
	"increment": TokenAttrIncrement,
	"hash":      TokenAttrHash,
	"unique":    TokenAttrUnique,
	"default":   TokenAttrDefault,
	"required":  TokenAttrNotNull,
	"hidden":    TokenAttrHidden,
	"check":     TokenAttrCheck,
}

var AllDataTypes = map[TokenType]struct{}{
	TokenTypeInt:       {},
	TokenTypeTimestamp: {},
	TokenTypeText:      {},
	TokenTypeFloat:     {},
	TokenTypeUuid:      {},
}

var allAttrs = map[TokenType]struct{}{
	TokenAttrIncrement: {},
	TokenAttrUnique:    {},
	TokenAttrNotNull:   {},
}

// var AnnotationOpens = map[TokenType]struct{}{
// 	TokenAttrOpen: {},
// 	// TokenListOpen: {},
// 	TokenConsOpen: {},
// }

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
		return "TOKEN_entity"
	case TokenTypeInt:
		return "TOKEN_int"
	case TokenTypeFloat:
		return "TOKEN_float"
	case TokenTypeText:
		return "TOKEN_text"
	case TokenTypeTimestamp:
		return "TOKEN_timestamp"
	case TokenTypeUuid:
		return "TOKEN_uuid"
	case TokenTypeRoutes:
		return "TOKEN_routes"
	case TokenAlter:
		return "TOKEN_alter"
	case TokenRef:
		return "TOKEN_ref"
	case TokenSelf:
		return "TOKEN_self"
	case TokenTypeBool:
		return "TOKEN_bool"
	case TokenEnd:
		return "TOKEN_end"
	case TokenEndpoint:
		return "TOKEN_endpoint"
	case TokenArrow:
		return "TOKEN_arrow"
	case TokenAttrOpen:
		return "TOKEN_enumopen"
	case TokenAttrClose:
		return "TOKEN_enumclose"
	// case TokenListOpen:
	// 	return "TOKEN_listopen"
	// case TokenListClose:
	// 	return "TOKEN_listclose"
	case TokenConsOpen:
		return "TOKEN_consopen"
	case TokenConsClose:
		return "TOKEN_consclose"
	case TokenComment:
		return "TOKEN_comment"
	case TokenNewline:
		return "TOKEN_newline"
	case TokenDot:
		return "TOKEN_dot"
	case TokenColon:
		return "TOKEN_colon"
	case TokenIdent:
		return "TOKEN_ident"
	case TokenString:
		return "TOKEN_string"
	case TokenDigits:
		return "TOKEN_digits"
	case TokenDigitsFloat:
		return "TOKEN_digitsfloat"
	case TokenAtSymbol:
		return "TOKEN_atsymbol"
	case TokenGet:
		return "TOKEN_get"
	case TokenPost:
		return "TOKEN_post"
	case TokenPut:
		return "TOKEN_put"
	case TokenDelete:
		return "TOKEN_delete"
	case TokenAttrIncrement:
		return "TOKEN_autoincrement"
	case TokenAttrUnique:
		return "TOKEN_unique"
	case TokenAttrNotNull:
		return "TOKEN_required"
	case TokenAttrDefault:
		return "TOKEN_default"
	case TokenEOF:
		return "TOKEN_eof"
	case TokenUnknown:
		return "TOKEN_unknown"
	default:
		return "Unknown token type"
	}
}
