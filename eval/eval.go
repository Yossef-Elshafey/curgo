package eval

import (
	"curgo/ast"
	"curgo/object"
	"fmt"
)

var (
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
	NULL  = &object.Null{}
)

// TODO : create Maybe<T>

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	case *ast.Boolean:
		return nativeBooleanObject(node.Value)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.UnaryExpression:
		right := Eval(node.Right)
		return evalUnaryExpression(node.Operator, right)
	case *ast.BinaryExpression:
		left := Eval(node.Left)
		right := Eval(node.Right)
		return evalBinaryExpr(node.Operator, left, right)
	case *ast.BlockStatment:
		return evalBlockStatments(node)
	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue)
		return &object.ReturnValue{Value: val}
	case *ast.IfExpression:
		return evalIfExpression(node)
	}
	return nil
}

func evalProgram(program *ast.Program) object.Object {
	var result object.Object
	for _, statement := range program.Statements {
		fmt.Printf("Currently processing stmt=%T, %+v", statement, statement)
		result = Eval(statement)
		if returnValue, ok := result.(*object.ReturnValue); ok {
			return returnValue.Value
		}
	}
	return result
}

func evalIfExpression(ie *ast.IfExpression) object.Object {
	condition := Eval(ie.Condition)
	if isTruthy(condition) {
		return Eval(ie.Consequence)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative)
	} else {
		return NULL
	}
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		return true
	}
}

func evalBinaryExpr(op string, left, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER && right.Type() == object.INTEGER:
		return evalIntegerBinaryExpr(op, left, right)
	case op == "==":
		return nativeBooleanObject(left == right)
	case op == "!=":
		return nativeBooleanObject(left != right)
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s", left.Type(), op, right.Type())
	default:
		return newError("unknown operator: %s %s %s", left.Type(), op, right.Type())
	}
}

func evalIntegerBinaryExpr(op string, left, right object.Object) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value
	switch op {
	case "+":
		return &object.Integer{Value: leftVal + rightVal}
	case "-":
		return &object.Integer{Value: leftVal - rightVal}
	case "*":
		return &object.Integer{Value: leftVal * rightVal}
	case "/":
		return &object.Integer{Value: leftVal / rightVal}
	case "<":
		return nativeBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBooleanObject(leftVal > rightVal)
	case "==":
		return nativeBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBooleanObject(leftVal != rightVal)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), op, right.Type())
	}
}

func evalUnaryExpression(op string, right object.Object) object.Object {
	switch op {
	case "!":
		return evalNotOpExpr(right)
	case "-":
		return evalMinusOpExpr(right)
	default:
		return newError("unknown operator:%s%s", op, right.Type())
	}
}

func evalMinusOpExpr(right object.Object) object.Object {
	if right.Type() != object.INTEGER {
		return newError("unknown operator: -%s", right.Type())
	}
	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

func evalNotOpExpr(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}
func evalBlockStatments(block *ast.BlockStatment) object.Object {
	var result object.Object

	for _, stmt := range block.Statements {
		result = Eval(stmt)
		if result != nil && result.Type() == object.RETURN_VALUE_OBJ {
			return result
		}
	}
	return result
}

func nativeBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}
