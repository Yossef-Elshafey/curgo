package object

import (
	"curgo/types/ast"
	"curgo/types/tokens"
	"fmt"
	"strconv"
	"strings"
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
	BOOLEAN       =  "BOOLEAN"
	ARRAY         =  "ARRAY"
	MAP           =  "MAP"
	NULL          =  "NULL"
	EXPECT        =  "EXPECT"
)

type Object interface { 
	Type() ObjectType
	Visit() string
}

type FetchFunction struct {
	Token   *token.Token
	Params  []*ast.Identifier
	Body    []ast.Statement
	Env     *Env
}

func (ff *FetchFunction) Type() ObjectType { return FUNCTION_OBJ}
func (ff *FetchFunction) Visit() string { 
	return fmt.Sprintf("fetchFunction %s", ff.Token.Value)
}

type Error struct {
	Message string
}

func (e *Error) Type() ObjectType { return ERROR_OBJ }
func (e *Error) Visit() string {return e.Message}

type String struct {
	Value string
}

func (s *String) Type() ObjectType { return STRING_OBJ }
func (s *String) Visit() string { return s.Value }

type Integer struct {
	Value int64
}

func (i *Integer) Type() ObjectType { return INTEGER_OBJ }
func (i *Integer) Visit() string { return strconv.Itoa(int(i.Value)) }

type CurgoCall struct {
	Key    string
	Value  Object
}

func (cc *CurgoCall) Type() ObjectType { return CALLPARAM}
func (cc *CurgoCall) Visit() string { return fmt.Sprintf("%s -> %s", cc.Key, cc.Value )}


type BuiltinFunction func(args ...Object) Object

type Builtin struct {
	Fn BuiltinFunction
}

func (b *Builtin) Type() ObjectType { return BUILTIN }
func (b *Builtin) Visit() string { return fmt.Sprintf("fn(%+v)", b.Fn) }

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType { return BOOLEAN }
func (b *Boolean) Visit() string { return fmt.Sprintf("%t", b.Value)}

type Null struct {}

func (n *Null) Type() ObjectType { return NULL }
func (n *Null) Visit() string { return fmt.Sprintf("%s", "Null")}

type Array struct {
	Elements []Object
}

func (a *Array) Type() ObjectType { return ARRAY }
func (a *Array) Visit() string {
    var out []string
    for _, el := range a.Elements {
        out = append(out, el.Visit())
    }
    return "[" + strings.Join(out, ", ") + "]"
}

type Map struct {
	Elements map[string]Object
}

func (m *Map) Type() ObjectType { return MAP }
func (m *Map) Visit() string { return fmt.Sprintf("%+v", m.Elements )}

type ExpectContext struct {
	Value Object
}

func (ec *ExpectContext) Type() ObjectType { return EXPECT }
func (ec *ExpectContext) Visit() string { return fmt.Sprintf("%+v", ec.Value)}

func (ec *ExpectContext) ToBe(v Object) (*ExpectContext, error) {
	fmt.Println(strings.Repeat("-", 20))
	fmt.Printf("Running\nexpect(%s, %s).toBe(%s, %s)\n",
					ec.Value.Type(), ec.Value.Visit(), v.Type(), v.Visit())
	if ec.Value.Visit() != v.Visit() {
		return nil ,fmt.Errorf("Expect failed, got=expect(%s, %s) toBe(%s, %s)",
					ec.Value.Type(), ec.Value.Visit(), v.Type(), v.Visit())
	}
	fmt.Printf("Passed\n")
	fmt.Println(strings.Repeat("-", 20))
	ec.Value = v
	return ec, nil
}
