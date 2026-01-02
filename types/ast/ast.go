package ast

import (
	"bytes"
	"curgo/lexer"
)

type Node interface {
	Stringify() string
}

type Statement interface {
	Node
	iStmt()
}

type Expression interface {
	Node
	iExpr()
}

type Program struct {
	Statements []Statement
}

func (p *Program) Stringify() string {
	return "Program"
}

type Identifier struct {
	Token lexer.Token
	Value  string
}

func (i *Identifier) iExpr() {}
func (i *Identifier) Stringify() string {
	var out bytes.Buffer
	out.WriteString(i.Token.Value)
	return out.String()
}

type FetchStmt struct {
	Token lexer.Token
	FetchIdentifier  *Identifier
	Body  []Statement
}

func (f *FetchStmt) iStmt() {}
func (fs *FetchStmt) Stringify() string {
	var out bytes.Buffer
	out.WriteString(fs.FetchIdentifier.Value)
	for _, stmt := range fs.Body {
		out.WriteString(stmt.Stringify())
	}
	return out.String()
}

type StringLiteral struct {
	Token lexer.Token
	Value string
}

func (sl *StringLiteral) Stringify() string {
	var out bytes.Buffer
	out.WriteString(sl.Token.Value)
	return out.String()
}

func (sl *StringLiteral) iExpr() {}

type CurgoAssignStatment struct {
	Arg   *Identifier
	Value Expression
}

func (ca *CurgoAssignStatment) iStmt() {}
func (ca *CurgoAssignStatment) Stringify() string {
	var out bytes.Buffer
	out.WriteString(ca.Arg.Stringify())
	out.WriteString("->")
	out.WriteString(ca.Value.Stringify())
	return out.String()
}
