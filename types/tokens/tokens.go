package tokens

import "log"

type TokenKind int

const (
	_ TokenKind = iota
	TRANSPILE_ASSIGN
	STRING
	IDENTIFIER

	SEMI_COLON
	BACKTICK
	PLUS
	COLON
	OPEN_CURLY
	CLOSE_CURLY
	COMMA

	FETCH
	DATA
	GLOBAL
	ENDFETCH

	EOF
	NEW_LINE
)

type Token struct {
	Kind  TokenKind
	Value string
	Pos   Position
}

type Position struct {
	Line  int
	Start int
	End   int
}

func TokenKindStringify(k TokenKind) string {
	switch k {
	case TRANSPILE_ASSIGN:
		return "transpile_assign"
	case STRING:
		return "string"
	case SEMI_COLON:
		return "semi_colon"
	case COLON:
		return "colon"
	case OPEN_CURLY:
		return "open_curly"
	case CLOSE_CURLY:
		return "close_curly"
	case FETCH:
		return "fetch"
	case DATA:
		return "data"
	case GLOBAL:
		return "global"
	case IDENTIFIER:
		return "identifier"
	case EOF:
		return "eof"
	case COMMA:
		return "comma"
	case ENDFETCH:
		return "endfetch"
	case BACKTICK:
		return "backtick"
	default:
		log.Fatalf("Cannot stringfy token: %d no case match", k)
		return ""
	}
}
