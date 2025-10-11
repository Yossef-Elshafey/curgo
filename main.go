package main

import (
	"curgo/lexer"
	"curgo/parser"
	"fmt"
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

	ast := parser.Parse(tokens)
	fmt.Printf("%+v\n", ast)
}
