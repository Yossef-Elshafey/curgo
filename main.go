package main

import (
	"curgo/eval"
	"curgo/lexer"
	"curgo/object"
	"curgo/parser"
	"curgo/repl"
	"fmt"
	"log"
	"os"
)

func main() {
	bytes, err := os.ReadFile("./examples/testFn.txt")
	if err != nil {
		log.Fatalf("Error: cannot open file")
	}
	source := string(bytes)
	tokens := lexer.Tokenize(source)
	p := parser.New(tokens)
	progarm := p.ParseProgram()
	env := object.NewEnvironment()
	ev := eval.Eval(progarm, env)
	fmt.Printf("Eval: %+v\n", ev)
	repl.Start(os.Stdin, os.Stdout)
}
