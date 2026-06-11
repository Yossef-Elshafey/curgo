package eval

import (
	"curgo/types/object"
	"fmt"
)


var builtins = map[string]*object.Builtin{
	"print": {
		Fn: func(args ...object.Object) object.Object {
			for _, arg := range args {
				fmt.Println(arg.Visit())
			}
			return nil
		},
	},
	"expect": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("assert expected one argument, got=%d", len(args))
			}
			ec := &object.ExpectContext{}
			ec.Value = args[0]
			return ec
		},
	},
}
