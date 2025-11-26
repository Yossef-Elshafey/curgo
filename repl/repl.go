package repl

import (
	"bufio"
	"curgo/eval"
	"curgo/lexer"
	"curgo/parser"
	"fmt"
	"io"
)

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	for {
		fmt.Printf(">> ")
		scanned := scanner.Scan()
		if !scanned {
			return
		}
		line := scanner.Text()
		l := lexer.Tokenize(line)
		p := parser.New(l)
		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}
		evaluated := eval.Eval(program)
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			// io.WriteString(out, "\n")
		}
		// io.WriteString(out, "Stringify: "+program.Stringify())
		io.WriteString(out, "\n")
	}
}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
