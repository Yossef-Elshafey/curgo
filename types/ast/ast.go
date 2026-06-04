package ast

import (
	"bytes"
	"curgo/types/tokens"
	"fmt"
)

type Node interface {
	GetLine() int
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

func (p *Program) GetLine() int {
	return -1
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
func (i *Identifier) GetLine() int { return i.Token.Line }

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

func (fs *FetchStmt) GetLine() int { return fs.Token.Line }

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
func (sl *StringLiteral) GetLine() int {return sl.Token.Line}

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

func (ca *CurgoAssignStatment) GetLine() int {return ca.Token.Line}

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

func (br *BinaryExpression) GetLine() int {return br.Token.Line}

type LetStatement struct {
	Token token.Token
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

func (ls *LetStatement) GetLine() int {return ls.Token.Line}

type ExpressionStatement struct {
	Token token.Token
	Expression Expression
}

func (es *ExpressionStatement) GetLine() int {return es.Token.Line}
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

func (ce *CallExpression) GetLine() int {return ce.Token.Line}
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

func (nl *NumberLiteral) GetLine() int {return nl.Token.Line}
func (nl *NumberLiteral) iExpr() {}
func (nl *NumberLiteral) Stringify() string {
	var out bytes.Buffer
	s := fmt.Sprintf("%d", nl.Value)
	out.WriteString(s)
	return out.String()
}

type RightOpts struct {
	Member *Identifier
	Callable bool
}

type SuffixExpression struct {
	Token token.Token
	Left Expression
	Operator string
	Right RightOpts
}

func (se *SuffixExpression) GetLine() int {return se.Token.Line}
func (se *SuffixExpression) iExpr() {}
func (se *SuffixExpression) Stringify() string {
	var out bytes.Buffer
	out.WriteString(se.Left.Stringify())
	out.WriteString(se.Operator)
	out.WriteString(se.Right.Member.Stringify())
	return out.String()
}

type BlockStatement struct {
	Token      token.Token
	Statements []Statement
}

func (bs *BlockStatement) Stringify() string {
	var out bytes.Buffer

	for _, s := range bs.Statements {
		out.WriteString(s.Stringify())
	}

	return out.String()
}
func (bs *BlockStatement) GetLine() int {return bs.Token.Line}
func (bs *BlockStatement) iStmt()       {}

type IfStmt struct {
	Token         token.Token
	Cond          Expression
	Consequences  *BlockStatement
	Alternative   *BlockStatement
}

func (is *IfStmt) GetLine() int {return is.Token.Line}
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

func (i *Indexing) GetLine() int {return i.Token.Line}
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

func (al *ArrayLiteral) GetLine() int {return al.Token.Line}
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

func (ml *MapLiteral) GetLine() int {return ml.Token.Line}
func (ml *MapLiteral) iExpr() {}
func (ml *MapLiteral) Stringify() string {
	var out bytes.Buffer
	out.WriteString(ml.Token.Value)
	for _, elm := range ml.Elements {
		out.WriteString(elm.Stringify())
	}
	return out.String()
}

type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) iExpr()      {}
func (pe *PrefixExpression) Stringify() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.Stringify())
	out.WriteString(")")

	return out.String()
}
func (pe *PrefixExpression) GetLine() int {return pe.Token.Line}
