package main

import (
	"curgo/eval"
	"curgo/lexer"
	"curgo/parser"
	"curgo/types/object"
	"curgo/utils"
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
	tokens := lexer.Tokenize(source)
	p := parser.Parse(tokens)
	env := object.NewEnvironment()
	eval.Eval(p,env)
}

func run() {
	// file := flag.String("i", "", "source file")
	// listTranspiler := flag.Bool("ls", false, "list all transpilation options")
	// flag.Parse()

	interp("./examples/02.txt")
	// if *file != "" {
	// }

}
