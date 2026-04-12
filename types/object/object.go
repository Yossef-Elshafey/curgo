package object

import (
	"curgo/lexer"
	"curgo/types/ast"
	"net/http"
)


type ObjectType string

const (
	INTEGER_OBJ   =  "INTEGER"
	ERROR_OBJ     =  "ERROR"
	FUNCTION_OBJ  =  "FUNCTION"
	STRING_OBJ    =  "STRING"
	CALLPARAM     =  "CALLPARAM"
	RESPONSE      =  "RESPONSE"
	BUILTIN       =  "BUILTIN"
)

type Object interface { 
	Type() ObjectType
}

type FetchFunction struct {
	Token   *lexer.Token
	Params  []*ast.Identifier
	Body    []ast.Statement
	Env     *Env
}

func (ff *FetchFunction) Type() ObjectType { return FUNCTION_OBJ}

type Error struct {
	Message string
}

func (e *Error) Type() ObjectType { return ERROR_OBJ }

type String struct {
	Value string
}

func (s *String) Type() ObjectType { return STRING_OBJ }

type Integer struct {
	Value int64
}

func (s *Integer) Type() ObjectType { return INTEGER_OBJ }

type CurgoCall struct {
	Key    string
	Value  Object
}

func (cc *CurgoCall) Type() ObjectType { return CALLPARAM}

type Response struct {
	Res *http.Response
}

func (r *Response) Type() ObjectType { return RESPONSE}


type BuiltinFunction func(args ...Object) Object

type Builtin struct {
	Fn BuiltinFunction
}

func (b *Builtin) Type() ObjectType { return BUILTIN }
