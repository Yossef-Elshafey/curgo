package lexer

type TokenKind int

const (
	OPEN_PAREN TokenKind = iota
	CLOSE_PAREN
	GLOBAL
	LOCAL
	SET
	ASSIGNMENT
)

var reserved map[string]TokenKind = map[string]TokenKind{
	"global": GLOBAL,
	"local":  LOCAL,
	"set":    SET,
}
