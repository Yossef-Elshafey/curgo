package eval

import (
	"bytes"
	"curgo/environment"
	"curgo/transpiler"
	"curgo/types/ast"
	"curgo/types/object"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

type Evaluater struct {
	transpiler *transpiler.CurlTranspiler
}

func Eval(node ast.Node, env *environment.Env) object.Object {
	e := &Evaluater{}
	e.transpiler = transpiler.New()
	switch node := node.(type) {
	case *ast.Program: e.evalProgram(node, env)

	case *ast.FetchStmt:
		ff := &object.FetchFunction{}
		ff.Body = node.Body
		ff.Params = node.Arguments
		env.Set(node.FetchIdentifier.Value, ff)
		return nil

	case *ast.CallExpression:
		function := Eval(node.Function, env)
		if isError(function) {
			return function
		}
		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}

		return applyFunction(function, args, env)

	case *ast.LetStatement: 
		v := Eval(node.Value, env)
		if v ==  nil { panic("Let stmt value is nil") }
		env.Set(node.Identifier.Token.Value, v)
		return nil

	case *ast.BinaryExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalInfixExpression(node.Operator.Value, left, right)
	case *ast.ExpressionStatement: return Eval(node.Expression, env)
	case *ast.Identifier: return evalIdentifier(node, env)
	case *ast.StringLiteral: return &object.String{Value: node.Value}

	}
	return nil
}

func evalInfixExpression(
	operator string,
	left, right object.Object,
) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evalStringInfixExpression(operator, left, right)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalStringInfixExpression(
	operator string,
	left, right object.Object,
) object.Object {
	if operator != "+" {
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}

	leftVal := left.(*object.String).Value
	rightVal := right.(*object.String).Value
	return &object.String{Value: leftVal + rightVal}
}

func evalIntegerInfixExpression(
	operator string,
	left, right object.Object,
) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch operator {
	case "+":
		return &object.Integer{Value: leftVal + rightVal}
	case "-":
		return &object.Integer{Value: leftVal - rightVal}
	case "*":
		return &object.Integer{Value: leftVal * rightVal}
	case "/":
		return &object.Integer{Value: leftVal / rightVal}
	default:
		return newError("unknown operator: %s %s %s",
			left.Type(), operator, right.Type())
	}
}

func applyFunction(fn object.Object, args []object.Object, env *environment.Env) object.Object {
	switch fn := fn.(type) {

	case *object.FetchFunction:
		fmt.Printf("%+v\n", fn.Body)
		fmt.Printf("%+v\n", fn.Params[0])
		fmt.Printf("%+v\n", args[0])
	}
	return newError("not a function: %T", fn)
}

func evalExpressions(
	exps []ast.Expression,
	env *environment.Env,
) []object.Object {
	var result []object.Object

	for _, e := range exps {
		evaluated := Eval(e, env)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}
	return result
}

func isError(obj object.Object) bool {
	_, ok := obj.(*object.Error)
	if ok && obj != nil {
		return true
	}
	return false
}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func evalIdentifier(
	node *ast.Identifier,
	env *environment.Env,
) object.Object {
	if val, ok := env.Get(node.Value); ok {
		return val
	}

	return newError("identifier not found: " + node.Value)
}


func (e *Evaluater) evalProgram(n *ast.Program, env *environment.Env) {
	for _, stmt := range n.Statements {
		Eval(stmt, env)
	}
}

func (e *Evaluater) fail(msg string) {
	log.Fatalf("%s", msg)
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
