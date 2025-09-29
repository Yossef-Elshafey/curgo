package lexer

import (
	"curgo/ast"
	"curgo/lexer/types"
	"fmt"
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
	gates    *types.Set[TokenKind]
	Ast      *ast.Ast
}

func NewLexer(source string) *lexer {
	return &lexer{
		handlers: []handlers{
			{regex: regexp.MustCompile(`^global`), token: GLOBAL},
			{regex: regexp.MustCompile(`^\w.+\{(.+\s?)$`), token: OPEN_PAREN},
			{regex: regexp.MustCompile(`^local`), token: LOCAL},
			{regex: regexp.MustCompile(`^set`), token: SET},
			{regex: regexp.MustCompile(`\}`), token: CLOSE_PAREN},
		},
		source: strings.Split(source, "\n"),
		pos:    0,
		gates:  types.NewSet[TokenKind](),
		Ast:    ast.NewAst(),
	}
}

func (l *lexer) Tokenize() {
	for _, line := range l.source {
		fmt.Printf("Currently processing line %s\n", line)
		l.searchForPattern(line)
	}
}

func (l *lexer) searchForPattern(line string) {
	// check for higher precedence patterns if founded it puts the token
	// at l.gates, if a gate was open (l.gates.HasValue)
	// it passes the next iteration line to whatever gate matched with the last pushed value
	for _, pattern := range l.handlers {
		if matched := pattern.regex.FindStringIndex(line); matched != nil {
			fmt.Printf("Line:%s has a match\n", line)
			l.gates.Put(pattern.token)
			return
		}
	}

	if l.gates.HasValue() {
		l.sendLineToGate(line)
	}
}

func (l *lexer) sendLineToGate(line string) {
	switch l.gates.GetLastValue() {
	case GLOBAL:
		l.gateAcknowledge(l.Ast.Global.IsChild(line), GLOBAL)
	case OPEN_PAREN:
		l.gateAcknowledge(openClosureHandler(line), OPEN_PAREN)
	case LOCAL:
		l.gateAcknowledge(localHandler(line), LOCAL)
	}
}

func (l *lexer) gateAcknowledge(acknowledge bool, token TokenKind) {
	if !acknowledge {
		l.gates.Delete(token)
	}
}

func openClosureHandler(line string) bool {
	if line == "" {
		return true
	}
	fmt.Printf("openClosure recived %v\n", line)
	return false
}

func localHandler(line string) bool {
	if line == "" {
		return true
	}
	fmt.Printf("localhandeler recived %s\n", line)
	return true
}
