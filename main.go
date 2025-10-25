package main

import (
	"curgo/lexer"
	"curgo/parser"
	"log"
	"os"
)

func main() {
	bytes, err := os.ReadFile("./examples/parser.txt")
	if err != nil {
		log.Fatalf("Error: cannot open file")
	}

	source := string(bytes)
	tokens := lexer.Tokenize(source)
	p := parser.New(tokens)
	p.ParseProgram()
}
