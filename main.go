package main

import (
	"curgo/eval"
	"curgo/lexer"
	"curgo/parser"
	"curgo/types/object"
	"fmt"
	"log"
	"os"
)

func main() {
	// if len(os.Args) < 2 {
	// 	log.Fatalf("repl is not yet implemented, pass any source file or use --help")
	// }

	run()
}


func interp(filename string) {
	f, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("cannot read file %s: %s\n", filename, err)
	}
	tokens := lexer.New(string(f))
	p := parser.New(tokens)
	program, err := p.ParseProgram()

	if err != nil {
		log.Fatal(err)
	}

	env := object.NewEnvironment()
	e, ok := eval.Eval(program, env).(*object.Error)
	if ok {
		fmt.Printf("Error != nil: %s\n", e.Visit())
	}
}

func run() {
	// file := flag.String("i", "", "source file")
	// listTranspiler := flag.Bool("ls", false, "list all transpilation options")
	// flag.Parse()

	interp("./examples/02.txt")
	// if *file != "" {
	// }

}
