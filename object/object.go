package object

import "fmt"

type ObjectType string

const (
	INTEGER          = "INTEGER"
	BOOLEAN          = "BOOLEAN"
	NULL             = "NULL"
	RETURN_VALUE_OBJ = "RETURN_VALUE"
	ERROR_OBJ        = "ERROR"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Integer struct {
	// TODO: use ceil(log2(decimal)) to fit decimal numbers
	Value int64
}

func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }
func (i *Integer) Type() ObjectType { return INTEGER }

type Boolean struct {
	Value bool
}

func (b *Boolean) Inspect() string  { return fmt.Sprintf("%t", b.Value) }
func (b *Boolean) Type() ObjectType { return BOOLEAN }

type Null struct{}

func (n *Null) Type() ObjectType { return NULL }
func (n *Null) Inspect() string  { return "null" }

type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJ }
func (rv *ReturnValue) Inspect() string  { return rv.Value.Inspect() }

type Error struct {
	Message string
}

func (e *Error) Type() ObjectType { return ERROR_OBJ }
func (e *Error) Inspect() string  { return "ERROR:" + e.Message }
