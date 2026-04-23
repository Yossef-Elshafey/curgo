package eval

import (
	"curgo/types/ast"
	"curgo/types/object"
	"fmt"
)

func Eval(node ast.Node, env *object.Env) object.Object {
	switch node := node.(type) {
	case *ast.Program: return evalProgram(node, env)

	case *ast.FetchStmt:
		ff := &object.FetchFunction{}
		ff.Body = node.Body
		ff.Params = node.Arguments
		ff.Token = &node.Token
		ff.Env = env
		env.Set(node.FetchIdentifier.Value, ff)
		return ff

	case *ast.CallExpression:
		function := Eval(node.Function, env)
		if isError(function) {
			return function
		}
		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}

		e := applyFunction(function, args)
		if isError(e) {
			return newError("Evaluator(%d): %s", node.Token.Line, e.Visit())
		}
		return e

	case *ast.LetStatement:
		v := Eval(node.Value, env)
		if isError(v) {
			return v
		}
		env.Set(node.Identifier.Token.Value, v)

	case *ast.BinaryExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		e := evalInfixExpression(node.Operator, left, right)
		if isError(e) {
			return newError("Evaluator(%d): %s", node.Token.Line, e.Visit())
		}
		return e

	case *ast.MemberAccess:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		member := node.Member.Value
		e := evalMemberAccessExpr(left, member)
		if isError(e) {
			return newError("Evaluator(%d): %s", node.Member.Token.Line, e.Visit())
		}
		return e
	case *ast.ExpressionStatement: return Eval(node.Expression, env)
	case *ast.Identifier: 
		e := evalIdentifier(node, env)
		if isError(e) {
			return newError("Evaluater(%d): %s", node.Token.Line, e.Visit())
		}
		return e
	case *ast.StringLiteral: return &object.String{Value: node.Value}
	case *ast.NumberLiteral: return &object.Integer{Value: node.Value}
	case *ast.CurgoAssignStatment: return &object.CurgoCall{Key: node.Arg.Value, Value: Eval(node.Value, env)}
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

func evalMemberAccessExpr(left object.Object, member string) object.Object {
	switch left := left.(type) {
	case *object.String:
		switch member {
		case "length":
			return &object.Integer{Value: int64(len(left.Value))}
		}
	case *object.Integer:
		switch member {
			case "plusone":
				left.Value = left.Value + 1
				return left
		}
	}
	return newError("%s doesnt support current option '%s'", left.Type(), member)
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

func applyFunction(fn object.Object, args []object.Object) object.Object {
	cr := New()
	switch fn := fn.(type) {
		case *object.FetchFunction:
			if len(args) != len(fn.Params) {
				return newError("fetch %s expects %d argument, got= %d",
					fn.Token.Value, len(fn.Params), len(args) )
			}
			extendedEnv := extendFunctionEnv(fn, args)
			for _, stmt := range fn.Body {
				cp, ok := Eval(stmt, extendedEnv).(*object.CurgoCall)
				if !ok {
					return newError("can't evluate stmt of %T", cp)
				}
				val, ok := cp.Value.(*object.String)
				if !ok { return newError("curgo request value is not string") }
				err := cr.buildRequest(cp.Key, val.Value)
				if err != nil { return newError(err.Error(), nil) }
			}
			resp, err := cr.fire()
			if err != nil { return newError(err.Error(), nil)}
			return &object.Response{Res: resp}
		case *object.Builtin:
			return fn.Fn(args...)
	}
	return newError("not a function: %T", fn)
}


func extendFunctionEnv(
	fn *object.FetchFunction,
	args []object.Object,
) *object.Env {
	env := object.NewOuterEnv(fn.Env)
	for idx, param := range fn.Params {
		env.Set(param.Value, args[idx])
	}
	return env
}

func evalExpressions(
	exps []ast.Expression,
	env *object.Env,
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
	env *object.Env,
) object.Object {
	if val, ok := env.Get(node.Value); ok {
		return val
	}

	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}

	return newError("identifier not found: %s", node.Value)
}

func evalProgram(n *ast.Program, env *object.Env) object.Object {
	var result object.Object
	for _, stmt := range n.Statements {
		result = Eval(stmt, env)

		switch result := result.(type) {
		case *object.Error:
			return result
		}
	}
	return result
}
