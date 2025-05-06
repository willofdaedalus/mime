package lexer

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
	TokenAlter    // alter
	TokenRef      // ref
	TokenSelf     // self
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
	TokenAmpersand // @
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

var Keywords = map[string]TokenType{
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
	if tok, ok := Keywords[ident]; ok {
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
	case TokenEnumOpen:
		return "TOKEN_enumopen"
	case TokenEnumClose:
		return "TOKEN_enumclose"
	case TokenListOpen:
		return "TOKEN_listopen"
	case TokenListClose:
		return "TOKEN_listclose"
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
	case TokenGet:
		return "TOKEN_get"
	case TokenPost:
		return "TOKEN_post"
	case TokenPut:
		return "TOKEN_put"
	case TokenDelete:
		return "TOKEN_delete"
	case TokenConstraintAutoIncrement:
		return "TOKEN_autoincrement"
	case TokenConstraintUnique:
		return "TOKEN_unique"
	case TokenConstraintForeignKey:
		return "TOKEN_foreignkey"
	case TokenConstraintPrimaryKey:
		return "TOKEN_primarykey"
	case TokenConstraintNotNull:
		return "TOKEN_required"
	case TokenConstraintDefault:
		return "TOKEN_default"
	case TokenEOF:
		return "TOKEN_eof"
	case TokenUnknown:
		return "TOKEN_unknown"
	default:
		return "Unknown token type"
	}
}
