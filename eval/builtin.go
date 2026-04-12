package eval

import (
	"curgo/types/object"
	"fmt"
)


var builtins = map[string]*object.Builtin{
	"print": {
		Fn: func(args ...object.Object) object.Object {
			for _, arg := range args {
				fmt.Println(arg.(*object.Response).Res)
			}
			return nil
		},
	},
}
