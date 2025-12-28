package main

import (
	"curgo/errors"
	"curgo/lexer"
)

func Interpret(source string) {
	errors.SetSourceName(source)
	lexer.Tokenize(source).Parse().Eval()
}
