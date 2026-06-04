package eval

import (
	"curgo/types/ast"
	"curgo/types/object"
)


// TODO: attach expectArgsCount(), markAsProperty(), etc to object.Builtin
func expectArgsCount(valid bool) object.Object {
	if !valid {
		return newError("invalid number of arguments")
	}
	return nil
}

func stringInterface(left *object.String, rhsOpts ast.RightOpts) object.Object {
	switch rhsOpts.Member.Value {
	case "length":
		if rhsOpts.Callable {
			return newError("length is not callable")
		}
		return &object.Integer{Value: int64(len(left.Value))}
	case "charIndex":
		fnObj := &object.Builtin{}
		fnObj.Fn = func(args ...object.Object) object.Object {
			err := expectArgsCount(len(args) == 1)
			if isError(err) {
				return newError("%s for %s", err.Visit(), rhsOpts.Member.Value)
			}
			idx, ok := args[0].(*object.Integer)
			if !ok {
				return newError("invalid argument type for %s expect number, got=%s", rhsOpts.Member.Value, args[0].Type())
			}
			return &object.String{Value: string(left.Value[idx.Value])}
		}
		return fnObj
	default:
		return newError("object %s donest support operation %s", left.Type(), rhsOpts.Member.Value)
	}
}

func responseInterface(left *object.Response, rhsOpts ast.RightOpts) object.Object {
	switch rhsOpts.Member.Value {
	case "status":
		return &object.Integer{Value: int64(left.Res.StatusCode)}
	case "statusText":
		return &object.String{Value: left.Res.Status}
	case "cookies":
		return newError("object %s doesn't support operation %s", left.Type(), rhsOpts.Member.Value)
	case "location":
		return newError("object %s doesn't support operation %s", left.Type(), rhsOpts.Member.Value)
	case "header":
		return newError("object %s doesn't support operation %s", left.Type(), rhsOpts.Member.Value)
	case "body":
		return newError("object %s doesn't support operation %s", left.Type(), rhsOpts.Member.Value)
	default:
		return newError("object %s doesn't support operation %s", left.Type(), rhsOpts.Member.Value)
	}
}
