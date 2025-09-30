package lexer

type TokenKind int

const (
	CLOSURE_START TokenKind = iota
	CLOSURE_END
	GLOBAL
)
