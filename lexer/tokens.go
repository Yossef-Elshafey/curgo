package lexer

import (
	"curgo/types/tokens"
	"fmt"
)

var reserved_lu map[string]tokens.TokenKind = map[string]tokens.TokenKind{
	"let":     tokens.LET,
	"import":  tokens.IMPORT,
	"foreach": tokens.FOREACH,
	"return":  tokens.RETURN,
	"if":      tokens.IF,
	"else":    tokens.ELSE,
	"fetch":   tokens.FETCH,
	"data":    tokens.DATA,
	"true":    tokens.TRUE,
	"false":   tokens.FALSE,
}

func NewToken(kind tokens.TokenKind, value string) tokens.Token {
	return tokens.Token{
		Type:  kind,
		Value: value,
	}
}

func TokenKindString(kind tokens.TokenKind) string {
	switch kind {
	case tokens.EOF:
		return "eof"
	case tokens.RETURN:
		return "return"
	case tokens.NULL:
		return "null"
	case tokens.NUMBER:
		return "number"
	case tokens.STRING:
		return "string"
	case tokens.TRUE:
		return "true"
	case tokens.FALSE:
		return "false"
	case tokens.IDENTIFIER:
		return "identifier"
	case tokens.OPEN_BRACKET:
		return "open_bracket"
	case tokens.CLOSE_BRACKET:
		return "close_bracket"
	case tokens.OPEN_CURLY:
		return "open_curly"
	case tokens.CLOSE_CURLY:
		return "close_curly"
	case tokens.OPEN_PAREN:
		return "open_paren"
	case tokens.CLOSE_PAREN:
		return "close_paren"
	case tokens.ASSIGNMENT:
		return "assignment"
	case tokens.EQUALS:
		return "equals"
	case tokens.NOT_EQUALS:
		return "not_equals"
	case tokens.NOT:
		return "not"
	case tokens.LESS:
		return "less"
	case tokens.LESS_EQUALS:
		return "less_equals"
	case tokens.GREATER:
		return "greater"
	case tokens.GREATER_EQUALS:
		return "greater_equals"
	case tokens.OR:
		return "or"
	case tokens.AND:
		return "and"
	case tokens.DOT:
		return "dot"
	case tokens.SEMI_COLON:
		return "semi_colon"
	case tokens.COLON:
		return "colon"
	case tokens.QUESTION:
		return "question"
	case tokens.COMMA:
		return "comma"
	case tokens.PLUS_EQUALS:
		return "plus_equals"
	case tokens.MINUS_EQUALS:
		return "minus_equals"
	case tokens.PLUS:
		return "plus"
	case tokens.DASH:
		return "dash"
	case tokens.SLASH:
		return "slash"
	case tokens.STAR:
		return "star"
	case tokens.PERCENT:
		return "percent"
	case tokens.LET:
		return "let"
	case tokens.IMPORT:
		return "import"
	case tokens.FROM:
		return "from"
	case tokens.IF:
		return "if"
	case tokens.ELSE:
		return "else"
	case tokens.FOREACH:
		return "foreach"
	case tokens.FOR:
		return "for"
	case tokens.IN:
		return "in"
	default:
		return fmt.Sprintf("unknown(%d)", kind)
	}
}
