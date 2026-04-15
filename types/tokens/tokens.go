package token

type TokenKind string

const (
	ILLEGAL          =  "ILLEGAL"
	EOF              =  "EOF"
	IDENTIFIER       =  "IDENTIFIER"
	NUMBER           =  "NUMBER"
	ASSIGN           =  "="
	PLUS             =  "+"
	MINUS            =  "-"
	BANG             =  "!"
	ASTERISK         =  "*"
	SLASH            =  "/"
	LT               =  "<"
	GT               =  ">"
	EQ               =  "=="
	NOT_EQ           =  "!="
	COMMA            =  ","
	SEMICOLON        =  ";"
	COLON            =  ":"
	LPAREN           =  "("
	RPAREN           =  ")"
	LBRACE           =  "{"
	RBRACE           =  "}"
	LBRACKET         =  "["
	RBRACKET         =  "]"
	TRANSPILEASSIGN  =  "->"
	COMMENT          =  "//"
	FETCH            =  "FETCH"
	ENDFETCH         =  "ENDFETCH"
	LET              =  "LET"
	TRUE             =  "TRUE"
	FALSE            =  "FALSE"
	IF               =  "IF"
	ELSE             =  "ELSE"
	STRING           =  "STRING"
)

type Token struct {
	Value  string
	Kind   TokenKind
}

var keywords = map[string]TokenKind{
	"fetch":     FETCH,
	"endfet":    ENDFETCH,
	"let":       LET,
	"true":      TRUE,
	"false":     FALSE,
	"if":        IF,
	"else":      ELSE,
}

func LookupIdent(ident string) TokenKind {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENTIFIER
}
