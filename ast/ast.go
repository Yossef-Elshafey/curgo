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

func (a *Ast) process(i int) {
	closure := a.body[i]

	command := a.acessOrFail(closure, "curl")
	name := a.acessOrFail(closure, "for")

	fmt.Println(strings.Repeat("-", 80))
	fmt.Printf("%s Response:\n", name)

	cmd := exec.Command("bash", "-c", command)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		log.Fatalf("Command failed with error: %v\nStderr: %s", err, stderr.String())
	}

	output := stdout.Bytes()
	var out bytes.Buffer
	err = json.Indent(&out, output, "", "  ")
	if err != nil {
		fmt.Println(string(output))
	} else {
		fmt.Printf("%s\n", out.Bytes())
	}
	fmt.Println(strings.Repeat("-", 80))
	fmt.Println()
}

func (a *Ast) Interpret(block int) {
	if block > len(a.body) {
		fmt.Printf("Out of range: %d\n", block)
		os.Exit(1)
	}

	if block == -1 {
		for i, _ := range a.body {
			a.process(i)
		}
	} else {
		a.process(block - 1)
	}
}
