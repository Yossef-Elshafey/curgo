package ast

import (
	"bytes"
	"curgo/types/tokens"
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
	Token token.Token
	Value  string
}

func (i *Identifier) iExpr() {}
func (i *Identifier) Stringify() string {
	var out bytes.Buffer
	out.WriteString(i.Token.Value)
	return out.String()
}

type FetchStmt struct {
	Token token.Token
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
	Token token.Token
	Value string
}

func (sl *StringLiteral) Stringify() string {
	var out bytes.Buffer
	out.WriteString(sl.Token.Value)
	return out.String()
}

func (sl *StringLiteral) iExpr() {}

type CurgoAssignStatment struct {
	Token token.Token
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
	Token token.Token
	Left Expression
	Operator string
	Right Expression
}

func (br *BinaryExpression) iExpr() {}
func (br *BinaryExpression) Stringify() string {
	var out bytes.Buffer
	out.WriteString(br.Left.Stringify())
	out.WriteString(br.Operator)
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
	Token token.Token
	Expression Expression
}

func (es *ExpressionStatement) iStmt() {}
func (es *ExpressionStatement) Stringify() string {
	var out bytes.Buffer
	out.WriteString(es.Expression.Stringify())
	return out.String()
}

type CallExpression struct {
	Token token.Token
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
	Token token.Token
	Value int64
}

func (nl *NumberLiteral) iExpr() {}
func (nl *NumberLiteral) Stringify() string {
	var out bytes.Buffer
	s := fmt.Sprintf("%d", nl.Value)
	out.WriteString(s)
	return out.String()
}

type SuffixExpression struct {
	Left Expression
	Operator string
	Member *Identifier
}

func (se *SuffixExpression) iExpr() {}
func (se *SuffixExpression) Stringify() string {
	var out bytes.Buffer
	out.WriteString(se.Left.Stringify())
	out.WriteString(se.Operator)
	out.WriteString(se.Member.Stringify())
	return out.String()
}

type BlockStatement struct {
	Token      token.Token
	Statements []Statement
}

func (bs *BlockStatement) iStmt()       {}
func (bs *BlockStatement) Stringify() string {
	var out bytes.Buffer

	for _, s := range bs.Statements {
		out.WriteString(s.Stringify())
	}

	return out.String()
}

type IfStmt struct {
	Token         token.Token
	Cond          Expression
	Consequences  *BlockStatement
	Alternative   *BlockStatement
}

func (is *IfStmt) iStmt() {}
func (is *IfStmt) Stringify() string {
	var out bytes.Buffer
	out.WriteString(is.Cond.Stringify())
	out.WriteString(is.Consequences.Stringify())
	return out.String()
}

type Indexing struct {
	Token   token.Token
	Ident   Expression
	Target  Expression
}

func (i *Indexing) iExpr() {}
func (i *Indexing) Stringify() string {
	var out bytes.Buffer
	out.WriteString(i.Ident.Stringify())
	out.WriteString(i.Target.Stringify())
	return out.String()
}


type ArrayLiteral struct {
	Token   token.Token
	Elements []Expression
}

func (al *ArrayLiteral) iExpr() {}
func (al *ArrayLiteral) Stringify() string {
	var out bytes.Buffer
	out.WriteString(al.Token.Value)
	for _, elm := range al.Elements {
		out.WriteString(elm.Stringify())
	}
	return out.String()
}

type MapLiteral struct {
	Token   token.Token
	Elements map[string]Expression
}

func (ml *MapLiteral) iExpr() {}
func (ml *MapLiteral) Stringify() string {
	var out bytes.Buffer
	out.WriteString(ml.Token.Value)
	for _, elm := range ml.Elements {
		out.WriteString(elm.Stringify())
	}
	return out.String()
}
