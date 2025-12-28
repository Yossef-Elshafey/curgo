package lexer

import (
	"curgo/types/tokens"
)

var reserved_lookup = map[string]tokens.TokenKind{
	"global": tokens.GLOBAL,
	"fetch":  tokens.FETCH,
	"endfet": tokens.ENDFETCH,
}

func NewToken(k tokens.TokenKind, v string, l, s, e int) tokens.Token {
	return tokens.Token{Kind: k, Value: v,
		Pos: tokens.Position{Line: l, Start: s, End: e},
	}
}
