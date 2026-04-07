package ast

import (
	"bytes"
	"curgo/lexer"
	"fmt"
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
	Arguments []*Identifier
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

type BinaryExpression struct {
	Left Expression
	Operator lexer.Token
	Right Expression
}

func (br *BinaryExpression) iExpr() {}
func (br *BinaryExpression) Stringify() string {
	var out bytes.Buffer
	out.WriteString(br.Left.Stringify())
	out.WriteString(br.Operator.Value)
	out.WriteString(br.Right.Stringify())
	return out.String()
}

type LetStatement struct {
	Identifier *Identifier
	Value Expression
}

func (ls *LetStatement) iStmt() {}
func (ls *LetStatement) Stringify() string {
	var out bytes.Buffer
	out.WriteString("let")
	out.WriteString(ls.Identifier.Stringify())
	out.WriteString("=")
	out.WriteString(ls.Value.Stringify())
	return out.String()
}

type ExpressionStatement struct {
	Token lexer.Token
	Expression Expression
}

func (es *ExpressionStatement) iStmt() {}
func (es *ExpressionStatement) Stringify() string {
	var out bytes.Buffer
	out.WriteString(es.Expression.Stringify())
	return out.String()
}

type CallExpression struct {
	Token lexer.Token
	Function Expression
	Arguments []Expression
}

func (ce *CallExpression) iExpr() {}
func (ce *CallExpression) Stringify() string {
	var out bytes.Buffer
	out.WriteString(ce.Function.Stringify())
	out.WriteString("(")
	for arg := range ce.Arguments {
		out.WriteString(ce.Arguments[arg].Stringify())
	}
	out.WriteString(")")
	return out.String()
}

type NumberLiteral struct {
	Token lexer.Token
	Value int64
}

func (nl *NumberLiteral) iExpr() {}
func (nl *NumberLiteral) Stringify() string {
	var out bytes.Buffer
	s := fmt.Sprintf("%d", nl.Value)
	out.WriteString(s)
	return out.String()
}
