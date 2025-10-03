package ast

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

/*
THIS IS SICK; since parsing by what should be there and what shouldn't is not the best but its the fasted to go with
there has some parsing and simple interpretation to avoid lack of flexibility and appropiate errors
TODO: pratt parsing of each tiny token
TODO: simple interpreter that is capable of the following
			bit operations
			real assignments
			syntax determination
			assignment handling for statements and exper
TODO: build an AST which in not that powerful for higher precedence but easy to move with for simple operations and assignments
*/

type Ast struct {
	Global  Global
	Closure Closure
	body    []map[string]string
}

func NewAst() Ast {
	return Ast{
		Global:  NewGlobal(),
		Closure: NewClosure(),
	}
}

func (a *Ast) ClosureDone() {
	json := a.Closure.toJson()
	a.body = append(a.body, json)
	a.Closure = NewClosure()
}

func (a *Ast) acessOrFail(assignments map[string]string, key string) string {
	value, exists := assignments[key]
	if !exists {
		fmt.Printf("Interpret: Failed to read %s command\n", key)
		os.Exit(1)
	}
	return value
}

func (a *Ast) Interpret() {
	for _, assignments := range a.body {
		command := a.acessOrFail(assignments, "curl")
		name := a.acessOrFail(assignments, "for")
		fmt.Printf("%s\n", name)
		cmd := exec.Command("bash", "-c", command)

		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		err := cmd.Run()
		if err != nil {
			log.Fatalf("Command failed with error: %v\nStderr: %s", err, stderr.String())
		}

		output := stdout.String()
		var out bytes.Buffer
		err = json.Indent(&out, []byte(output), "", "  ")
		fmt.Printf("%s\n", out.Bytes())
		fmt.Println(strings.Repeat("-", 20))
	}
}
