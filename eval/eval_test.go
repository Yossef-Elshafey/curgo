package eval

import (
	"curgo/lexer"
	"curgo/parser"
	"curgo/types/object"
	"fmt"
	"strings"
	"testing"
)

func TestSingleObjectIfCond(t *testing.T) {
	source := `
	let x = 1 == 1;
	if x { 200 }
	`
	tokens := lexer.New(source)
	p := parser.New(tokens)
	program, _ := p.ParseProgram()
	env := object.NewEnvironment()
	e := Eval(program, env)
	obj := e.(*object.Integer)
	if obj.Value != 200 {
		t.Errorf("expect= %d, instead got %d", 200, obj.Value)
	}
}

func TestIfCondition(t *testing.T) {
	sources := []struct {
		input    string
		expected int
	}{
		{`if (1 + 2) == 3 {let x = 1; x}`, 1},
		{`
			if 1 != 1 { let x = 1; x }
			else { let x = 2; x } `, 2},
		{`
			let x = 0;
			if 1 != 1 { let x = 1; x }
			else { let x = 2; x } `, 2},

		{`
			let x = 100;
			if "hello" == "hello" { x }
			`, 100},

		{`
			let x = 100;
			if ("hello" + " " + "yossef") != "hello yossef" { x } else {200}
			`, 200},
	}
	for _, source := range sources {
		tokens := lexer.New(source.input)
		p := parser.New(tokens)
		program, _ := p.ParseProgram()
		env := object.NewEnvironment()
		e := Eval(program, env)
		obj := e.(*object.Integer)
		if obj.Value != int64(source.expected) {
			t.Errorf("expect= %d, instead got %d", source.expected, obj.Value)
		}
	}
}

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
	source := []struct {
		input    string
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
			t.Errorf("expect %d, got= %d", s.expected, i.Value)
		}
	}
}

func TestStringConcat(t *testing.T) {
	source := `let x = "hello"; let y = "world"; x + y`
	tokens := lexer.New(source)
	p := parser.New(tokens)
	program, _ := p.ParseProgram()
	env := object.NewEnvironment()
	eval := Eval(program, env)
	s, ok := eval.(*object.String)
	if !ok {
		t.Errorf("expect object.String, got=%T", eval)
	}
	if s.Value != "helloworld" {
		t.Errorf("expect 'helloworld', got= %s", s.Value)
	}
}

func TestStringLength(t *testing.T) {
	source := `let s = "test"; s.length`
	tokens := lexer.New(source)
	p := parser.New(tokens)
	program, _ := p.ParseProgram()
	env := object.NewEnvironment()
	eval := Eval(program, env)
	i, ok := eval.(*object.Integer)
	if !ok {
		t.Errorf("expect object.Integer, got=%T", eval)
	}
	if i.Value != 4 {
		t.Errorf("expect 4, got= %d", i.Value)
	}
}

func TestLetStatement(t *testing.T) {
	source := `
	let x = 10;
	let y = 20;
	x + y
	`
	tokens := lexer.New(source)
	p := parser.New(tokens)
	program, _ := p.ParseProgram()
	env := object.NewEnvironment()
	eval := Eval(program, env)
	i, ok := eval.(*object.Integer)
	if !ok {
		t.Errorf("expect object.Integer, got=%T", eval)
	}
	if i.Value != 30 {
		t.Errorf("expect 30, got= %d", i.Value)
	}
}

func TestBuiltinPrint(t *testing.T) {
	source := `print("test")`
	tokens := lexer.New(source)
	p := parser.New(tokens)
	program, _ := p.ParseProgram()
	env := object.NewEnvironment()
	eval := Eval(program, env)
	if eval != nil {
		t.Errorf("expect nil, got=%T", eval)
	}
}

func TestEvalEmptyProgram(t *testing.T) {
	source := ``
	tokens := lexer.New(source)
	p := parser.New(tokens)
	program, _ := p.ParseProgram()
	env := object.NewEnvironment()
	eval := Eval(program, env)
	if eval != nil {
		t.Errorf("expect nil, got=%T", eval)
	}
}

func TestEvalIdentifierFromEnv(t *testing.T) {
	source := `x`
	tokens := lexer.New(source)
	p := parser.New(tokens)
	program, _ := p.ParseProgram()
	env := object.NewEnvironment()
	env.Set("x", &object.String{Value: "hello"})
	eval := Eval(program, env)
	s, ok := eval.(*object.String)
	if !ok {
		t.Errorf("expect object.String, got=%T", eval)
	}
	if s.Value != "hello" {
		t.Errorf("expect 'hello', got= %s", s.Value)
	}
}

func TestEvalNumberLiteral(t *testing.T) {
	source := `42`
	tokens := lexer.New(source)
	p := parser.New(tokens)
	program, _ := p.ParseProgram()
	env := object.NewEnvironment()
	eval := Eval(program, env)
	i, ok := eval.(*object.Integer)
	if !ok {
		t.Errorf("expect object.Integer, got=%T", eval)
	}
	if i.Value != 42 {
		t.Errorf("expect 42, got= %d", i.Value)
	}
}

func TestEvalStringLiteral(t *testing.T) {
	source := `"curgo"`
	tokens := lexer.New(source)
	p := parser.New(tokens)
	program, _ := p.ParseProgram()
	env := object.NewEnvironment()
	eval := Eval(program, env)
	s, ok := eval.(*object.String)
	if !ok {
		t.Errorf("expect object.String, got=%T", eval)
	}
	if s.Value != "curgo" {
		t.Errorf("expect 'curgo', got= %s", s.Value)
	}
}

func TestExpressionStatement(t *testing.T) {
	source := `(5 + 3)`
	tokens := lexer.New(source)
	p := parser.New(tokens)
	program, _ := p.ParseProgram()
	env := object.NewEnvironment()
	eval := Eval(program, env)
	i, ok := eval.(*object.Integer)
	if !ok {
		t.Errorf("expect object.Integer, got=%T", eval)
	}
	if i.Value != 8 {
		t.Errorf("expect 8, got= %d", i.Value)
	}
}

func TestIdentifierShadowing(t *testing.T) {
	source := `
	let x = 5;
	let x = 10;
	x
	`
	tokens := lexer.New(source)
	p := parser.New(tokens)
	program, _ := p.ParseProgram()
	env := object.NewEnvironment()
	eval := Eval(program, env)
	i, ok := eval.(*object.Integer)
	if !ok {
		t.Errorf("expect object.Integer, got=%T", eval)
	}
	if i.Value != 10 {
		t.Errorf("expect 10, got= %d", i.Value)
	}
}

func TestDivisionByZero(t *testing.T) {
	source := `
	let x = 1 / 0;
	x
	`
	tokens := lexer.New(source)
	p := parser.New(tokens)
	program, _ := p.ParseProgram()
	env := object.NewEnvironment()
	eval := Eval(program, env)
	i, ok := eval.(*object.Error)
	if !ok {
		t.Errorf("expect object.error, got=%T", eval)
	}
	if !strings.Contains(i.Message, "division by zero") {
		t.Errorf("error message donest include 'division by zero'")
	}
}

func TestArrayLiteral(t *testing.T) {
	source := `
	let y = 4;
	let x = [1 + 2, y, 5];
	x
	`
	tokens := lexer.New(source)
	p := parser.New(tokens)
	program, _ := p.ParseProgram()
	env := object.NewEnvironment()
	e := Eval(program, env)
	arr, ok := e.(*object.Array)
	if !ok {
		t.Errorf("expect object.Array, got= %T", e)
	}
	if len(arr.Elements) != 3 {
		t.Errorf("expect 3 elements in array, got=%d", len(arr.Elements))
	}
	v := arr.Elements[1]
	vInt, ok := v.(*object.Integer)
	if !ok {
		t.Errorf("expect Elements[2] to be integer, got= %T", v)
	}
	if vInt.Value != 4 {
		t.Errorf("expect elements[2] to be 4, got= %d", v)
	}

	v = arr.Elements[0]
	vInt, ok = v.(*object.Integer)
	if !ok {
		t.Errorf("expect Elements[0] to be integer, got= %T", v)
	}
	if vInt.Value != 3 {
		t.Errorf("expect elements[0] to be 3, got= %d", v)
	}
}

func TestMapLiteral(t *testing.T) {
	source := `
	let cord = {x:2, y:4};
	cord["x"]
	`
	tokens := lexer.New(source)
	p := parser.New(tokens)
	program, _ := p.ParseProgram()
	env := object.NewEnvironment()
	e := Eval(program, env)
	ret, ok := e.(*object.Integer)
	if !ok {
		t.Errorf("expect return to be object.Integer, got=%T", e)
	}
	if ret.Value != 2 {
		t.Errorf("expect return to be 3, got=%d", ret.Value)
	}
}

func TestIndexing(t *testing.T) {
	source := `
	let x = [1 + 2, 5];
	x[1]
	`
	tokens := lexer.New(source)
	p := parser.New(tokens)
	program, _ := p.ParseProgram()
	env := object.NewEnvironment()
	e := Eval(program, env)
	fmt.Printf("%+v\n", e)
	ret, ok := e.(*object.Integer)
	if !ok {
		t.Errorf("expect return to be object.Integer, got=%T", e)
	}
	if ret.Value != 7 {
		t.Errorf("expect return to be 3, got=%d", ret.Value)
	}
}

func TestInvalidIndexing(t *testing.T) {
	sources := []struct{
		input string
	}{
		{ "let x = [1 + 2, 5]; x[2]" },
		{ "let x = [1 + 2, 5]; x[-1]" },
		{ `let x = {x:50, y:100}; x["x"] + x["y"] + x["z"]` },
	}
	for _, source := range sources {
		tokens := lexer.New(source.input)
		p := parser.New(tokens)
		program, _ := p.ParseProgram()
		env := object.NewEnvironment()
		e := Eval(program, env)
		ret, ok := e.(*object.Error)
		if !ok {
			t.Errorf("expect return to be object.Error, got=%T", e)
		}
		fmt.Printf("%+v\n", ret)
	}
}
