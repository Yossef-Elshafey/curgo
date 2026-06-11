package eval

import (
	"curgo/types/ast"
	"curgo/types/object"
	"fmt"
	"io"
)

func expectArgsCount(valid bool) object.Object {
	if !valid {
		return newError("invalid number of arguments")
	}
	return nil
}

func markAsProperty(opts ast.RightOpts) object.Object {
	if opts.Callable {
		return newError("%s is not callable", opts.Member.Value)
	}
	return nil
}

func stringInterface(left *object.String, rhsOpts ast.RightOpts) object.Object {
	switch rhsOpts.Member.Value {
	case "length":
		err := markAsProperty(rhsOpts)
		if isError(err) {
			return err
		}
		return &object.Integer{Value: int64(len(left.Value))}
	default:
		return newError("object %s donest support operation %s", left.Type(), rhsOpts.Member.Value)
	}
}


func expectContextInterface(left *object.ExpectContext, rhsOpts ast.RightOpts) object.Object {
	switch rhsOpts.Member.Value {
	case "toBe":
		fnObj := &object.Builtin{}
		fnObj.Fn = func(args ...object.Object) object.Object {
			err := expectArgsCount(len(args) == 1)
			if isError(err) {
				return newError("invalid argument count for toBe function expect 1, got=%d", len(args))
			}
			evaluated, failed := left.ToBe(args[0])
			if failed != nil {
				return newError(failed.Error())
			}
			return evaluated
		}
		return fnObj
	case "unWrap":
		fnObj := &object.Builtin{}
		fnObj.Fn = func(args ...object.Object) object.Object {
			err := expectArgsCount(len(args) == 0)
			if isError(err) {
				return newError("invalid argument count for unWrap expect 0, got=%d", len(args))
			}
			return left.Value
		}
		return fnObj
	default:
		return newError("object %s doesn't support operation %s", left.Type(), rhsOpts.Member.Value)
	}
}

func responseInterface(left *object.Response, rhsOpts ast.RightOpts) object.Object {
	switch rhsOpts.Member.Value {
	case "status":
		return &object.Integer{Value: int64(left.Res.StatusCode)}
	case "get":
		fnObj := &object.Builtin{}
		fnObj.Fn = func(args ...object.Object) object.Object {
			err := expectArgsCount(len(args) == 1)
			if isError(err) {
				return newError("invalid argument count for get function expect 1, got=%d", len(args))
			}
			key, ok := args[0].(*object.String)
			if !ok {
				return newError("invalid argument type for %s expect string, got=%s", rhsOpts.Member.Value, args[0].Type())
			}
			return &object.String{Value: left.Res.Header.Get(key.Value)}
		}
		return fnObj
	case "statusText":
		return &object.String{Value: left.Res.Status}
	case "cookies":
		return newError("object %s doesn't support operation %s", left.Type(), rhsOpts.Member.Value)
	case "location":
		return newError("object %s doesn't support operation %s", left.Type(), rhsOpts.Member.Value)
	case "header":
		return newError("object %s doesn't support operation %s", left.Type(), rhsOpts.Member.Value)
	case "body":
		fnObj := &object.Builtin{}
		fnObj.Fn = func(args ...object.Object) object.Object {
			defer left.Res.Body.Close()
			bodyBytes, bodyErr := io.ReadAll(left.Res.Body)
			if bodyErr != nil {
				fmt.Println("Couldn't read body", bodyErr.Error())
			}
			bodyString := string(bodyBytes)
			return &object.String{Value: bodyString}
		}
		return fnObj
	default:
		return newError("object %s doesn't support operation %s", left.Type(), rhsOpts.Member.Value)
	}
}
