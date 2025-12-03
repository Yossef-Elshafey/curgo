package repl

import (
	"bufio"
	"curgo/eval"
	"curgo/lexer"
	"curgo/object"
	"curgo/parser"
	"fmt"
	"io"
	"os"
)

const (
	_ = iota
	USER_QUIT
	NEXT_INST
)

// TODO: complete repl functionalities in terms of the relationship between parsing and repl keywords
func pipe(fn func(inp string) int, inp string) {
	ret := fn(inp)
	if ret != NEXT_INST {
		os.Exit(ret)
	}
}

func clear(inp string) int {
	if inp == "clear" {
		fmt.Printf("\033c")
	}
	return NEXT_INST
}

func quit(inp string) int {
	if inp == "quit" {
		return USER_QUIT
	}
	return NEXT_INST
}

func pipeInpToReplKeywords(inp string) {
	pipe(quit, inp)
	pipe(clear, inp)
}

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()
	for {
		fmt.Printf(">> ")
		scanned := scanner.Scan()
		if !scanned {
			return
		}
		line := scanner.Text()
		pipeInpToReplKeywords(line)
		l := lexer.Tokenize(line)
		p := parser.New(l)
		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}
		evaluated := eval.Eval(program, env)
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}
}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
