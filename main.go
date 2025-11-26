package main

import (
	"curgo/eval"
	"curgo/lexer"
	"curgo/parser"
	_ "curgo/repl"
	_ "fmt"
	"log"
	"os"
)

func main() {
	bytes, err := os.ReadFile("./examples/02.txt")
	if err != nil {
		log.Fatalf("Error: cannot open file")
	}
	source := string(bytes)
	tokens := lexer.Tokenize(source)
	// fmt.Printf("%+v", tokens)
	p := parser.New(tokens)
	program := p.ParseProgram()

	eval.Eval(program)
	// repl.Start(os.Stdin, os.Stdout)

}
