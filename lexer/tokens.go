package lexer

import (
	"curgo/types/tokens"
)

func NewToken(k tokens.TokenKind, v string, l, s, e int) Token {
	return Token{Kind: k, Value: v,
		Pos: Position{Line: l, Start: s, End: e},
	}
}

type Token struct {
	Kind  tokens.TokenKind
	Value string
	Pos   Position
}

type Position struct {
	Line  int
	Start int
	End   int
}

