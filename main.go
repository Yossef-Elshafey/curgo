package main

import (
	"log"
	"os"
	"parser/lexer"
	"parser/parser"
)

func main() {
	bytes, err := os.ReadFile("./examples/parser.txt")
	if err != nil {
		log.Fatalf("Error: cannot open file")
	}
	source := string(bytes)
	tokens := lexer.Tokenize(source)

	ast := parser.Parse(tokens)
}
