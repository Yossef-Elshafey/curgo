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
		t.Errorf("Evaluater should return object.Error, got= %T", eval)
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

func TestArthmeticOps(t *testing.T) {
	source := []struct{
		input string
		expected int
	}{
		{"2 - 2 + 1 * 10 / 2 - 2", 3},
		{"(2 * 2) + (2 / 2)", 5},
		{"((2 + 2) * 1)", 4},
		{"(100 + 100) / 50 - (20 + 5)", -21},
	}
	for _, s := range source {
		tokens := lexer.New(s.input)
		p := parser.New(tokens)
		program, _ := p.ParseProgram()
		env := object.NewEnvironment()
		eval := Eval(program, env)
		i, ok := eval.(*object.Integer)
		if !ok {
			t.Errorf("expect object.Integer, got=%T", eval)
		}
		if i.Value != int64(s.expected) {
			t.Errorf("expect %d, got= %d",s.expected, i.Value)
		}
	}
}
