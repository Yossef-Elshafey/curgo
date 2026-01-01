package main

import (
	"curgo/eval"
	"curgo/lexer"
	"curgo/parser"
)

func Interpret(source string) {
	tokens := lexer.Tokenize(source)
	p := parser.Parse(tokens)
	eval.Eval(p)
}
