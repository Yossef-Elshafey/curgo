package lexer

import (
	"curgo/types/tokens"
	"fmt"
	"log"
	"regexp"
)

type regexHandler func(lex *lexer, regex *regexp.Regexp, tokenCh chan<- tokens.Token)

type regexPattern struct {
	regex   *regexp.Regexp
	handler regexHandler
}

type lexer struct {
	patterns []regexPattern
	source   string
	pos      int
	line     int
}

func Tokenize(source string, tokenCh chan<- tokens.Token) {
	defer close(tokenCh)
	lex := createLexer(source)
	for !lex.isEOF() {
		matched := false
		for _, pattern := range lex.patterns {
			loc := pattern.regex.FindStringIndex(lex.remainder())
			if loc != nil && loc[0] == 0 {
				pattern.handler(lex, pattern.regex, tokenCh)
				matched = true
				break
			}
		}
		if !matched {
			log.Fatalf("Lexer: Unrecognized Token near %s", lex.remainder())
		}
	}

	tokenCh <- NewToken(tokens.EOF, "EOF")
}

func (l *lexer) shift(n int) {
	l.pos += n
}

func (l *lexer) remainder() string {
	return l.source[l.pos:]
}

func (l *lexer) stream(out chan<- tokens.Token, token tokens.Token) {
	fmt.Printf("Streaming: %+v\n", token)
	out <- token
}

func (l *lexer) isEOF() bool {
	return l.pos >= len(l.source)
}

func defaultHandler(kind tokens.TokenKind, value string) regexHandler {
	return func(lex *lexer, regex *regexp.Regexp, tokenCh chan<- tokens.Token) {
		lex.shift(len(value))
		lex.stream(tokenCh, NewToken(kind, value))
	}
}

func numberHandler(lex *lexer, regex *regexp.Regexp, tokenCh chan<- tokens.Token) {
	match := regex.FindString(lex.remainder())
	lex.stream(tokenCh, NewToken(tokens.NUMBER, match))
	lex.shift(len(match))
}

func skipHandler(lex *lexer, regex *regexp.Regexp, tokenCh chan<- tokens.Token) {
	skip := regex.FindString(lex.remainder())
	lex.shift(len(skip))
}

func stringHandler(lex *lexer, regex *regexp.Regexp, tokenCh chan<- tokens.Token) {
	match := regex.FindStringIndex(lex.remainder())
	literal := lex.remainder()[match[0]:match[1]]
	lex.stream(tokenCh, NewToken(tokens.STRING, literal[1:len(literal)-1]))
	lex.shift(len(literal))
}

func symbolHandler(lex *lexer, regex *regexp.Regexp, tokenCh chan<- tokens.Token) {
	match := regex.FindString(lex.remainder())
	if kind, exists := reserved_lu[match]; exists {
		lex.stream(tokenCh, NewToken(kind, match))
	} else {
		lex.stream(tokenCh, NewToken(tokens.IDENTIFIER, match))
	}
	lex.shift(len(match))
}

func multiLineStringHandler(lex *lexer, regex *regexp.Regexp, tokenCh chan<- tokens.Token) {
	match := regex.FindStringIndex(lex.remainder())
	literal := lex.remainder()[match[0]:match[1]]
	lex.stream(tokenCh, NewToken(tokens.STRING, literal))
	lex.shift(len(literal))
}

func createLexer(source string) *lexer {
	return &lexer{
		patterns: []regexPattern{
			{regexp.MustCompile(`[a-zA-Z_][a-zA-Z0-9_]*`), symbolHandler},
			{regexp.MustCompile(`"([^"\\]*(\\.[^"\\]*)*)"`), stringHandler},
			{regexp.MustCompile(`(?s)\x60.*?(\n*?)\x60`), multiLineStringHandler}, // x60 = `
			{regexp.MustCompile(`\s+`), skipHandler},
			{regexp.MustCompile(`\/\/.*`), skipHandler},
			{regexp.MustCompile(`[0-9]+(\.[0-9]+)?`), numberHandler},
			{regexp.MustCompile(`\[`), defaultHandler(tokens.OPEN_BRACKET, "[")},
			{regexp.MustCompile(`\]`), defaultHandler(tokens.CLOSE_BRACKET, "]")},
			{regexp.MustCompile(`\{`), defaultHandler(tokens.OPEN_CURLY, "{")},
			{regexp.MustCompile(`\}`), defaultHandler(tokens.CLOSE_CURLY, "}")},
			{regexp.MustCompile(`\(`), defaultHandler(tokens.OPEN_PAREN, "(")},
			{regexp.MustCompile(`\)`), defaultHandler(tokens.CLOSE_PAREN, ")")},
			{regexp.MustCompile(`==`), defaultHandler(tokens.EQUALS, "==")},
			{regexp.MustCompile(`!=`), defaultHandler(tokens.NOT_EQUALS, "!=")},
			{regexp.MustCompile(`=`), defaultHandler(tokens.ASSIGNMENT, "=")},
			{regexp.MustCompile(`!`), defaultHandler(tokens.NOT, "!")},
			{regexp.MustCompile(`<=`), defaultHandler(tokens.LESS_EQUALS, "<=")},
			{regexp.MustCompile(`<`), defaultHandler(tokens.LESS, "<")},
			{regexp.MustCompile(`>=`), defaultHandler(tokens.GREATER_EQUALS, ">=")},
			{regexp.MustCompile(`>`), defaultHandler(tokens.GREATER, ">")},
			{regexp.MustCompile(`\|\|`), defaultHandler(tokens.OR, "||")},
			{regexp.MustCompile(`&&`), defaultHandler(tokens.AND, "&&")},
			{regexp.MustCompile(`\.`), defaultHandler(tokens.DOT, ".")},
			{regexp.MustCompile(`;`), defaultHandler(tokens.SEMI_COLON, ";")},
			{regexp.MustCompile(`:`), defaultHandler(tokens.COLON, ":")},
			{regexp.MustCompile(`\?`), defaultHandler(tokens.QUESTION, "?")},
			{regexp.MustCompile(`,`), defaultHandler(tokens.COMMA, ",")},
			{regexp.MustCompile(`\+=`), defaultHandler(tokens.PLUS_EQUALS, "+=")},
			{regexp.MustCompile(`-=`), defaultHandler(tokens.MINUS_EQUALS, "-=")},
			{regexp.MustCompile(`\+`), defaultHandler(tokens.PLUS, "+")},
			{regexp.MustCompile(`-`), defaultHandler(tokens.DASH, "-")},
			{regexp.MustCompile(`/`), defaultHandler(tokens.SLASH, "/")},
			{regexp.MustCompile(`\*`), defaultHandler(tokens.STAR, "*")},
			{regexp.MustCompile(`%`), defaultHandler(tokens.PERCENT, "%")},
		},
		source: source,
		pos:    0,
	}
}
