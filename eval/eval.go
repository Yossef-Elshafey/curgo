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
	transpiler *transpiler.CurlTranspiler
}

func Eval(node ast.Node) {
	e := &Evaluater{}
	e.transpiler = transpiler.New()
	switch n := node.(type) {
	case *ast.Program:
		e.evalProgram(n)
	case *ast.FetchStmt: 
		e.evalFetchStmt(n)
	}
}

func (e *Evaluater) evalProgram(n *ast.Program) {
	for _, stmt := range n.Statements {
		Eval(stmt)
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
	e.executeCurlCommand(n.FetchIdentifier.Value, cmd) // TODO:
}

func (e *Evaluater) executeCurlCommand(title, command string) {
	// https://www.sohamkamani.com/golang/exec-shell-command/
	var stdout bytes.Buffer
	cmd := exec.Command("/bin/sh", "-c", "curl"+command)
	cmd.Stdout = &stdout

	if err := cmd.Run(); err != nil {
		fmt.Printf("Stdout: %s\n", stdout.String())
		fmt.Printf("Command failed with %s\n", err)
	}

	fmt.Printf("Response: %s\n", stdout.String())
	fmt.Printf("%s\n", strings.Repeat("-", 10))
}
