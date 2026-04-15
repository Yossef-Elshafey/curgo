package main

import (
	"curgo/eval"
	"curgo/lexer"
	"curgo/parser"
	"curgo/types/object"
	"curgo/utils"
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
		log.Fatalf("cant read file:%s", filename)
	}

	source := string(f)
	utils.SetSource(source)
	tokens := lexer.New(source)
	p := parser.New(tokens)
	program := p.ParseProgram()
	env := object.NewEnvironment()
	e := eval.Eval(program, env)
	if e != nil {fmt.Printf("%s\n", e.Visit())}
}

func run() {
	// file := flag.String("i", "", "source file")
	// listTranspiler := flag.Bool("ls", false, "list all transpilation options")
	// flag.Parse()

	interp("./examples/02.txt")
	// if *file != "" {
	// }

}
