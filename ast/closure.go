package ast

import (
	"curgo/lexer/types"
	"fmt"
	"regexp"
)

type pattern struct {
	regex   *regexp.Regexp
	token   int
	handler func(line string) bool
}

const (
	MUTLI_LINE_ASSIGN = iota
	SINGLE_LINE_ASSIGN
	CURL
	STATEMENT
)

type Closure struct {
	patterns []pattern
	gates    *types.Set[int]
}

// Currently :TODO: implement handlers, implement communications between closure and AST
func NewClosure() *Closure {
	return &Closure{
		patterns: []pattern{
			{regex: regexp.MustCompile(`(?s)([a-zA-Z_][a-zA-Z0-9_]*)\s*=\s*\{.*?\}\s*;`), token: SINGLE_LINE_ASSIGN, handler: handleSingleLineAssign},
			{regex: regexp.MustCompile(`(?s)([a-zA-Z_][a-zA-Z0-9_]*)\s*=\s.*?{$`), token: MUTLI_LINE_ASSIGN, handler: handleMutliLineAssign},
			{regex: regexp.MustCompile(`curl`), token: CURL, handler: handleCurl},
			{regex: regexp.MustCompile(`\$\w.*=(\s.+?)\w.\.\w.+;`), token: STATEMENT, handler: handleStatement},
		},
		gates: types.NewSet[int](),
	}
}

func handleSingleLineAssign(line string) bool {
	fmt.Println("Single line assign", line)
	return false
}

func handleMutliLineAssign(line string) bool {
	fmt.Println("multi line assign", line)
	return true
}

func handleCurl(line string) bool {
	fmt.Println("curl handler", line)
	return true
}

func handleStatement(line string) bool {
	fmt.Println("Statement", line)
	return false
}

func (c *Closure) sendLineToHandler(line string) {
	switch c.gates.GetLastValue() {
	case MUTLI_LINE_ASSIGN:
		handleMutliLineAssign(line)
	case SINGLE_LINE_ASSIGN:
		handleSingleLineAssign(line)
	case CURL:
		handleCurl(line)
	case STATEMENT:
		handleStatement(line)
	}
}

func (c *Closure) IsChild(line string) bool {
	if c.gates.HasValue() {
		c.sendLineToHandler(line)
	}

	for _, pattern := range c.patterns {
		matched := pattern.regex.MatchString(line)
		if matched {
			c.gates.Put(pattern.token)
			pattern.handler(line)
		}
	}

	return true
}
