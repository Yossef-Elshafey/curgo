//	for _, assignment := range e.Program.Statements[0].(*ast.FetchStmt).Body {
//		a := assignment.(*ast.CurgoAssignStatment)
//		fmt.Printf("Assignment: %+v\nPure Value ---%s---\n", a.Value, a.Value.(*ast.StringLiteral).Value)
//	}
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
	fmt.Printf("Eval: %+v\n", e.Program)
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
		argument, ok := e.transpiler.Get(stmt.Arg.Value)
		if !ok {
			e.fail(fmt.Sprintf("Eval: Argument not found %s", stmt.Arg.Value))
		}
		value, ok := stmt.Value.(*ast.StringLiteral)
		if !ok {
			fmt.Printf("Eval: Cannot convert to ast.StringLiteral")
		}
		cmd = cmd + fmt.Sprintf("%s %s ", argument, value.Value)
	}
	e.executeCurlCommand(n.FetchIdentifier.Value, cmd)
}

func (e *Evaluater) executeCurlCommand(title, command string) {
	cmd := exec.Command("curl", command)
	fmt.Printf("Command: %s",cmd.String())
	cmd.Stdin = strings.NewReader("Foo")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		e.fail(err.Error())
	}
	fmt.Printf("%s: %q\n",title, out.String())
}
