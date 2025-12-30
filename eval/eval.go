package eval

import (
	"bytes"
	"curgo/transpiler"
	"curgo/types/ast"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

type Evaluater struct {
	Program ast.Program
	transpiler *transpiler.CurlTranspiler
}

func (e *Evaluater) Eval() {
	e.transpiler = transpiler.New()
	for _, stmt := range e.Program.Statements {
		if n, ok := stmt.(*ast.FetchStmt); ok { e.evalFetchStmt(n) }
	}
}

func (e *Evaluater) fail(msg string) {
	log.Fatalf("%s", msg)
}

func (e *Evaluater) evalFetchStmt(n *ast.FetchStmt) {
	var cmd string
	for _, stmt := range n.Body {
		stmt, ok := stmt.(*ast.CurgoAssignStatment)
		if !ok {
			// body content is []Statements interface but now it only contains ast.CurgoAssignStatment
			e.fail(fmt.Sprintf("FetchStmt.body is not ast.CurgoAssignStatment, got=%T", stmt))
		}
		stringLiteral, ok := stmt.Value.(*ast.StringLiteral)
		if !ok {
			fmt.Printf("Eval: Cannot convert to ast.StringLiteral")
		}
		argument, value := e.transpiler.Get(stmt.Arg.Value, stringLiteral.Value)
		cmd = cmd + fmt.Sprintf("%s %s ", argument, value)
	}
	e.executeCurlCommand(n.FetchIdentifier.Value, cmd)
}

func (e *Evaluater) executeCurlCommand(title, command string) {
	// https://www.sohamkamani.com/golang/exec-shell-command/
	fmt.Printf("%s\n",strings.Repeat("-",10))
	fmt.Printf("Executing %s...\n",title)

	var stdout, stderr bytes.Buffer
	cmd := exec.Command("/bin/sh", "-c", "curl" + command)
	cmd.Stdout = &stdout
	err := cmd.Run()
	
	if err != nil {
		log.Fatalf("Command failed with %s: %s", err, stderr.String())
	}

	fmt.Printf("Response: %s\n", stdout.String())
	fmt.Printf("%s\n",strings.Repeat("-",10))
}
