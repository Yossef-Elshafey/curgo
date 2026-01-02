package main

import (
	"curgo/eval"
	"curgo/lexer"
	"curgo/parser"
	"curgo/utils"
)

func Interpret(source string) {
	utils.SetSource(source)
	tokens := lexer.Tokenize(source)
	p := parser.Parse(tokens)
	eval.Eval(p)
}
