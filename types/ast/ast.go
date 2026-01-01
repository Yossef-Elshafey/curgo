package ast

import (
	"curgo/lexer"
)

type Node interface {
	GetPositions() lexer.Position
	Stringify() string
}

type Statement interface {
	iStmt()
}

type Expression interface {
	iExpr()
}

type Program struct {
	Statements []Statement
}

type Identifier struct {
	Token lexer.Token
	Value  string
}

func (i *Identifier) iExpr() {}

type FetchStmt struct {
	Token lexer.Token
	FetchIdentifier  *Identifier
	Body  []Statement
}

type StringLiteral struct {
	Token lexer.Token
	Value string
}

func (sl *StringLiteral) iExpr() {}

type CurgoAssignStatment struct {
	Arg   *Identifier
	Value Expression
}

func (ca *CurgoAssignStatment) iStmt() {}

func (f *FetchStmt) iStmt() {}
