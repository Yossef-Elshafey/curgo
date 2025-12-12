package main

import (
	"curgo/lexer"
	"curgo/parser"
	"curgo/types/tokens"
	"log"
	"os"
)

func main() {
	input, error := os.ReadFile("./examples/01.txt")
	if error != nil {
		log.Fatalf("Cannot read file")
	}
	tokenCh := make(chan tokens.Token)
	go lexer.Tokenize(string(input), tokenCh)
	p := parser.NewParser()
	p.Parse(tokenCh)
}
