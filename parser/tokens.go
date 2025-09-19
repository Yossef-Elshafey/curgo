package parser

type TokenKind int

const (
	EOF TokenKind = iota
	BLOCK_START
	BLOCK_END
	CURL
	DEFINE
)

type Token struct {
	kind  TokenKind
	value string
}
