package lexer

import (
	"curgo/types/tokens"
	"log"
	"regexp"
	"strings"
)

type regexHandler func(lex *lexer, regex *regexp.Regexp, location []int)

type regexPattern struct {
	regex   *regexp.Regexp
	handler regexHandler
}

type lexer struct {
	patterns []regexPattern
	tokens   []Token
	source   string
	pos      int
	char     int
	line     int
}

func Tokenize(source string) []Token {
	lex := createLexer(source)
	for !lex.isEOF() {
		matched := false
		for _, pattern := range lex.patterns {
			loc := pattern.regex.FindStringIndex(lex.remainder())
			if loc != nil && loc[0] == 0 {
				matched = true
				pattern.handler(lex, pattern.regex, loc)
				break
			}
		}
		if !matched {
			limit := min((len(lex.remainder())), 10)
			lines := strings.ReplaceAll(string(lex.remainder()[0]) + lex.remainder()[0:limit], "","")
			log.Fatalf("Lexer:%d: Unrecognized '%s'%s\n'%s' is not supported",
				lex.line,
				string(lex.remainder()[0]),
				lines, string(lex.remainder()[0]))
		}
	}

	lex.push(NewToken(tokens.EOF, "EOF", lex.line, 0, -1))
	return lex.tokens
}

func (l *lexer) shift(n int) {
	l.pos += n
	l.char += n
}

func (l *lexer) push(token Token) {
	l.tokens = append(l.tokens, token)
}

func (l *lexer) getPositions(end int) (int, int, int) {
	return l.line, l.char, l.char + end
}

func (l *lexer) remainder() string {
	return l.source[l.pos:]
}

func (l *lexer) isEOF() bool {
	return l.pos >= len(l.source)
}

func defaultHandler(kind tokens.TokenKind, value string) regexHandler {
	return func(lex *lexer, regex *regexp.Regexp, loc []int) {
		lex.shift(len(value))
		line, start, end := lex.getPositions(loc[1])
		lex.push(NewToken(kind, value, line, start, end))
	}
}

// func numberHandler(lex *lexer, regex *regexp.Regexp) {
// 	match := regex.FindString(lex.remainder())
// 	lex.push(NewToken(NUMBER, match))
// 	lex.shift(len(match))
// }

func skipHandler(lex *lexer, regex *regexp.Regexp, loc []int) {
	skip := regex.FindString(lex.remainder())
	lex.shift(len(skip))
}

func stringHandler(lex *lexer, regex *regexp.Regexp, loc []int) {
	match := regex.FindString(lex.remainder())
	line, start, end := lex.getPositions(loc[1])
	doubleQuote := 2
	lex.push(NewToken(tokens.STRING, match[1:len(match)-1], line, start, end-doubleQuote))
	lex.shift(len(match))
}

func symbolHandler(lex *lexer, regex *regexp.Regexp, loc []int) {
	match := regex.FindString(lex.remainder())
	line, start, end := lex.getPositions(loc[1])
	if kind, exists := tokens.Reserved_Keyword[match]; exists {
		lex.push(NewToken(kind, match, line, start, end))
	} else {
		lex.push(NewToken(tokens.IDENTIFIER, match, line, start, end))
	}
	lex.shift(len(match))
}

func curlTranspileAssignment(lex *lexer, regex *regexp.Regexp, loc []int) {
	match := regex.FindString(lex.remainder())
	line, start, end := lex.getPositions(loc[1])
	lex.push(NewToken(tokens.TRANSPILE_ASSIGN, match, line, start, end))
	lex.shift(len(match))
}

func newLineHandler(lex *lexer, regex *regexp.Regexp, loc []int) {
	skipHandler(lex, regex, loc)
	lex.line += 1
	lex.char = 0
}

func multiLineStringHandler(lex *lexer, regex *regexp.Regexp, loc []int) {
	match := regex.FindString(lex.remainder())
	target := strings.ReplaceAll(match, "\n", "")
	target = strings.TrimSpace(target[1:len(target)-1])
	line, start, end := lex.getPositions(loc[1])
	asterisk := 2
	lex.push(NewToken(tokens.STRING, target, line, start, end-asterisk))
	lex.shift(len(match))
}

func createLexer(source string) *lexer {
	return &lexer{
		patterns: []regexPattern{
			{regexp.MustCompile(`[a-zA-Z_][a-zA-Z0-9_]*`), symbolHandler},
			{regexp.MustCompile(`"([^"\\]*(\\.[^"\\]*)*)"`), stringHandler},
			{regexp.MustCompile(`\/\/.*`), skipHandler},
			{regexp.MustCompile(`->`), curlTranspileAssignment},
			{regexp.MustCompile(`(?s)\x60(.*?)\x60`), multiLineStringHandler}, // x60 = `
			{regexp.MustCompile(`\n`), newLineHandler},
			{regexp.MustCompile(`\s+`), skipHandler},
			{regexp.MustCompile(`\/\/.*`), skipHandler},
			// {regexp.MustCompile(`[0-9]+(\.[0-9]+)?`), numberHandler},
			{regexp.MustCompile(`\{`), defaultHandler(tokens.OPEN_CURLY, "{")},
			{regexp.MustCompile(`\+`), defaultHandler(tokens.PLUS, "+")},
			{regexp.MustCompile(`,`), defaultHandler(tokens.COMMA, ",")},
			{regexp.MustCompile(`\}`), defaultHandler(tokens.CLOSE_CURLY, "}")},
			{regexp.MustCompile(`;`), defaultHandler(tokens.SEMI_COLON, ";")},
			{regexp.MustCompile(`:`), defaultHandler(tokens.COLON, ":")},
		},
		source: source,
		line:   1,
	}
}

// patterns: []regexPattern{
// 	{regexp.MustCompile(`[a-zA-Z_][a-zA-Z0-9_]*`), symbolHandler},
// 	{regexp.MustCompile(`"([^"\\]*(\\.[^"\\]*)*)"`), stringHandler},
// 	{regexp.MustCompile(`(?s)\x60.*?(\n*?)\x60`), multiLineStringHandler}, // x60 = `
// 	{regexp.MustCompile(`\s+`), skipHandler},

// 	{regexp.MustCompile(`[0-9]+(\.[0-9]+)?`), numberHandler},
// 	{regexp.MustCompile(`\[`), defaultHandler(OPEN_BRACKET, "[")},
// 	{regexp.MustCompile(`\]`), defaultHandler(CLOSE_BRACKET, "]")},
// 	{regexp.MustCompile(`\{`), defaultHandler(OPEN_CURLY, "{")},
// 	{regexp.MustCompile(`\}`), defaultHandler(CLOSE_CURLY, "}")},
// 	{regexp.MustCompile(`\(`), defaultHandler(OPEN_PAREN, "(")},
// 	{regexp.MustCompile(`\)`), defaultHandler(CLOSE_PAREN, ")")},
// 	{regexp.MustCompile(`==`), defaultHandler(EQUALS, "==")},
// 	{regexp.MustCompile(`!=`), defaultHandler(NOT_EQUALS, "!=")},
// 	{regexp.MustCompile(`=`), defaultHandler(ASSIGNMENT, "=")},
// 	{regexp.MustCompile(`!`), defaultHandler(NOT, "!")},
// 	{regexp.MustCompile(`<=`), defaultHandler(LESS_EQUALS, "<=")},
// 	{regexp.MustCompile(`<`), defaultHandler(LESS, "<")},
// 	{regexp.MustCompile(`>=`), defaultHandler(GREATER_EQUALS, ">=")},
// 	{regexp.MustCompile(`>`), defaultHandler(GREATER, ">")},
// 	{regexp.MustCompile(`\|\|`), defaultHandler(OR, "||")},
// 	{regexp.MustCompile(`&&`), defaultHandler(AND, "&&")},
// 	{regexp.MustCompile(`\.\.`), defaultHandler(DOT_DOT, "..")},
// 	{regexp.MustCompile(`\.`), defaultHandler(DOT, ".")},
// 	{regexp.MustCompile(`;`), defaultHandler(SEMI_COLON, ";")},
// 	{regexp.MustCompile(`:`), defaultHandler(COLON, ":")},
// 	{regexp.MustCompile(`\?\?=`), defaultHandler(NULLISH_ASSIGNMENT, "??=")},
// 	{regexp.MustCompile(`\?`), defaultHandler(QUESTION, "?")},
// 	{regexp.MustCompile(`,`), defaultHandler(COMMA, ",")},
// 	{regexp.MustCompile(`\+\+`), defaultHandler(PLUS_PLUS, "++")},
// 	{regexp.MustCompile(`--`), defaultHandler(MINUS_MINUS, "--")},
// 	{regexp.MustCompile(`\+=`), defaultHandler(PLUS_EQUALS, "+=")},
// 	{regexp.MustCompile(`-=`), defaultHandler(MINUS_EQUALS, "-=")},
// 	{regexp.MustCompile(`\+`), defaultHandler(PLUS, "+")},
// 	{regexp.MustCompile(`-`), defaultHandler(DASH, "-")},
// 	{regexp.MustCompile(`/`), defaultHandler(SLASH, "/")},
// 	{regexp.MustCompile(`\*`), defaultHandler(STAR, "*")},
// 	{regexp.MustCompile(`%`), defaultHandler(PERCENT, "%")},
// },
