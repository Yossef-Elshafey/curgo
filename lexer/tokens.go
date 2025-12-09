package lexer

import "fmt"

type TokenKind int

type Token struct {
	Value string
	Type  TokenKind
}

// Token Enum
const (
	EOF TokenKind = iota
	NULL
	TRUE
	FALSE
	NUMBER
	STRING
	IDENTIFIER

	OPEN_BRACKET
	CLOSE_BRACKET
	OPEN_CURLY
	CLOSE_CURLY
	OPEN_PAREN
	CLOSE_PAREN

	ASSIGNMENT
	EQUALS
	NOT_EQUALS
	NOT

	LESS
	LESS_EQUALS
	GREATER
	GREATER_EQUALS

	OR
	AND

	DOT
	SEMI_COLON
	COLON
	QUESTION
	COMMA

	PLUS_EQUALS
	MINUS_EQUALS

	PLUS
	DASH
	SLASH
	STAR
	PERCENT

	FETCH
	DATA

	LET
	IMPORT
	FROM
	IF
	ELSE
	FOREACH
	FOR
	TYPEOF
	IN
	RETURN
)

var reserved_lu map[string]TokenKind = map[string]TokenKind{
	"let":     LET,
	"import":  IMPORT,
	"foreach": FOREACH,
	"return":  RETURN,
	"if":      IF,
	"else":    ELSE,
	"fetch":   FETCH,
	"data":    DATA,
	"true":    TRUE,
	"false":   FALSE,
}

func NewToken(kind TokenKind, value string) Token {
	return Token{
		Type:  kind,
		Value: value,
	}
}

func TokenKindString(kind TokenKind) string {
	switch kind {
	case EOF:
		return "eof"
	case RETURN:
		return "return"
	case NULL:
		return "null"
	case NUMBER:
		return "number"
	case STRING:
		return "string"
	case TRUE:
		return "true"
	case FALSE:
		return "false"
	case IDENTIFIER:
		return "identifier"
	case OPEN_BRACKET:
		return "open_bracket"
	case CLOSE_BRACKET:
		return "close_bracket"
	case OPEN_CURLY:
		return "open_curly"
	case CLOSE_CURLY:
		return "close_curly"
	case OPEN_PAREN:
		return "open_paren"
	case CLOSE_PAREN:
		return "close_paren"
	case ASSIGNMENT:
		return "assignment"
	case EQUALS:
		return "equals"
	case NOT_EQUALS:
		return "not_equals"
	case NOT:
		return "not"
	case LESS:
		return "less"
	case LESS_EQUALS:
		return "less_equals"
	case GREATER:
		return "greater"
	case GREATER_EQUALS:
		return "greater_equals"
	case OR:
		return "or"
	case AND:
		return "and"
	case DOT:
		return "dot"
	case SEMI_COLON:
		return "semi_colon"
	case COLON:
		return "colon"
	case QUESTION:
		return "question"
	case COMMA:
		return "comma"
	case PLUS_EQUALS:
		return "plus_equals"
	case MINUS_EQUALS:
		return "minus_equals"
	case PLUS:
		return "plus"
	case DASH:
		return "dash"
	case SLASH:
		return "slash"
	case STAR:
		return "star"
	case PERCENT:
		return "percent"
	case LET:
		return "let"
	case IMPORT:
		return "import"
	case FROM:
		return "from"
	case IF:
		return "if"
	case ELSE:
		return "else"
	case FOREACH:
		return "foreach"
	case FOR:
		return "for"
	case IN:
		return "in"
	default:
		return fmt.Sprintf("unknown(%d)", kind)
	}
}
