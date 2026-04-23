package eval

import (
	"curgo/lexer"
	"curgo/parser"
	"curgo/types/object"
	"testing"
)

func TestTypeMismatchInfixExpression(t *testing.T) {
	source := `
	let x = "";
	let y = 5;
	print(x + y);
	`
	tokens := lexer.New(source)
	p := parser.New(tokens)
	program, _ := p.ParseProgram()
	env := object.NewEnvironment()
	eval := Eval(program, env)
	_, ok := eval.(*object.Error)
	if !ok {
		t.Errorf("Evaluater should return object.Error, got= %t", eval)
	}
}

func TestMemberAccessError(t *testing.T) {
	source := `
	let a = "";
	a.foo
	`
	tokens := lexer.New(source)
	p := parser.New(tokens)
	program, _ := p.ParseProgram()
	env := object.NewEnvironment()
	eval := Eval(program, env)
	_, ok := eval.(*object.Error)
	if !ok {
		t.Errorf("Evaluater should return object.Error, got= %T", eval)
	}
}

func TestUnknownIdentifier(t *testing.T) {
	source := `
	copy();
	`
	tokens := lexer.New(source)
	p := parser.New(tokens)
	program, _ := p.ParseProgram()
	env := object.NewEnvironment()
	eval := Eval(program, env)
	_, ok := eval.(*object.Error)
	if !ok {
		t.Errorf("Evaluater should return object.Error, got= %T", eval)
	}
}

func TestIncorrectCallExpr(t *testing.T) {
	source := `
	fetch test(url):
		host -> url;
	endfet
	test()
	`
	tokens := lexer.New(source)
	p := parser.New(tokens)
	program, _ := p.ParseProgram()
	env := object.NewEnvironment()
	eval := Eval(program, env)
	_, ok := eval.(*object.Error)
	if !ok {
		t.Errorf("Evaluater should return object.Error, got= %T", eval)
	}
}

func TestIncorrectStringConcat(t *testing.T) {
	source := `
	let x = "";
	let y = "";
	print(x - y)
	`
	tokens := lexer.New(source)
	p := parser.New(tokens)
	program, _ := p.ParseProgram()
	env := object.NewEnvironment()
	eval := Eval(program, env)
	_, ok := eval.(*object.Error)
	if !ok {
		t.Errorf("Evaluater should return object.Error, got= %T", eval)
	}
}
