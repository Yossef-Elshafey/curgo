package lexer

import (
	"curgo/ast"
	"regexp"
	"strings"
)

type handlers struct {
	regex *regexp.Regexp
	token TokenKind
}

type lexer struct {
	handlers []handlers
	source   []string
	pos      int
	gate     TokenKind
	Ast      ast.Ast
}

func NewLexer(source string) *lexer {
	return &lexer{
		handlers: []handlers{
			{regex: regexp.MustCompile(`^global`), token: GLOBAL},
			{regex: regexp.MustCompile(`^<>$`), token: CLOSURE_START},
			{regex: regexp.MustCompile(`end<>`), token: CLOSURE_END},
		},
		source: strings.Split(source, "\n"),
		pos:    0,
		gate:   -1,
		Ast:    ast.NewAst(),
	}
}

func (l *lexer) Tokenize() {
	for _, line := range l.source {
		l.searchForPattern(line)
	}
	l.Ast.Interpret()
}

func (l *lexer) searchForPattern(line string) {
	for _, pattern := range l.handlers {
		if matched := pattern.regex.MatchString(line); matched {
			l.gate = pattern.token
			return
		}
	}

	if l.gate != -1 {
		l.sendLineToGate(line)
	}
}

func (l *lexer) sendLineToGate(line string) {
	switch l.gate {
	case GLOBAL:
		l.gateAcknowledge(l.Ast.Global.IsChild(line))
	case CLOSURE_START:
		l.gateAcknowledge(l.Ast.Closure.IsChild(line, l.Ast.Global))
	case CLOSURE_END:
		l.gate = -1
		l.Ast.ClosureDone()
	}
}

func (l *lexer) gateAcknowledge(acknowledge bool) {
	if !acknowledge {
		l.gate = -1
	}
}
