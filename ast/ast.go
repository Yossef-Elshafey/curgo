package ast

import (
	"bytes"
	"curgo/lexer"
	"strings"
)

type Node interface {
	TokenLiteral() string
	Stringify() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

func (p *Program) Stringify() string {
	var out bytes.Buffer
	for _, s := range p.Statements {
		out.WriteString(s.Stringify())
	}
	return out.String()
}

// TODO: change Token to be lexer.TokenKind instead of the struct

type LetStatment struct {
	Token lexer.Token
	Name  *Identifier
	Value Expression
}

func (ls *LetStatment) statementNode()       {}
func (ls *LetStatment) TokenLiteral() string { return ls.Token.Value }
func (ls *LetStatment) Stringify() string {
	var out bytes.Buffer
	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.Stringify())
	out.WriteString(" = ")
	if ls.Value != nil {
		out.WriteString(ls.Value.Stringify())
	}
	out.WriteString(";")
	return out.String()
}

type Identifier struct {
	Token lexer.Token
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Value }
func (i *Identifier) Stringify() string {
	return i.Value
}

type ReturnStatement struct {
	Token       lexer.Token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Value }
func (rs *ReturnStatement) Stringify() string {
	var out bytes.Buffer
	out.WriteString(rs.TokenLiteral() + " ")
	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.Stringify())
	}
	out.WriteString(";")
	return out.String()
}

type ExpressionStatement struct {
	Token      lexer.Token
	Expression Expression
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Value }
func (es *ExpressionStatement) Stringify() string {
	if es.Expression != nil {
		return es.Expression.Stringify()
	}
	return ""
}

type IntegerLiteral struct {
	Token lexer.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Value }
func (il *IntegerLiteral) Stringify() string    { return il.Token.Value }

type UnaryExpression struct {
	Token    lexer.Token
	Operator string
	Right    Expression
}

func (ue *UnaryExpression) expressionNode()      {}
func (ue *UnaryExpression) TokenLiteral() string { return ue.Token.Value }
func (ue *UnaryExpression) Stringify() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(ue.Operator)
	out.WriteString(ue.Right.Stringify())
	out.WriteString(")")
	return out.String()
}

type BinaryExpression struct {
	Token    lexer.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (be *BinaryExpression) expressionNode()      {}
func (be *BinaryExpression) TokenLiteral() string { return be.Token.Value }
func (be *BinaryExpression) Stringify() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(be.Left.Stringify())
	out.WriteString(" " + be.Operator + " ")
	out.WriteString(be.Right.Stringify())
	out.WriteString(")")
	return out.String()
}

type Boolean struct {
	Token lexer.Token
	Value bool
}

func (bl *Boolean) expressionNode()      {}
func (bl *Boolean) TokenLiteral() string { return bl.Token.Value }
func (bl *Boolean) Stringify() string    { return bl.Token.Value }

type IfExpression struct {
	Token       lexer.Token
	Condition   Expression
	Consequence *BlockStatment
	Alternative *BlockStatment
}

func (ie *IfExpression) expressionNode()      {}
func (ie *IfExpression) TokenLiteral() string { return ie.Token.Value }
func (ie *IfExpression) Stringify() string {
	var out bytes.Buffer
	out.WriteString("if")
	out.WriteString(ie.Condition.Stringify())
	out.WriteString(" ")
	out.WriteString(ie.Consequence.Stringify())
	if ie.Alternative != nil {
		out.WriteString("else ")
		out.WriteString(ie.Alternative.Stringify())
	}
	return out.String()
}

type BlockStatment struct {
	Token      lexer.Token
	Statements []Statement
}

func (bs *BlockStatment) expressionNode()      {}
func (bs *BlockStatment) TokenLiteral() string { return bs.Token.Value }
func (bs *BlockStatment) Stringify() string {
	var out bytes.Buffer
	for _, s := range bs.Statements {
		out.WriteString(s.Stringify())
	}
	return out.String()
}

type FunctionLiteral struct {
	Token  lexer.Token
	Params []*Identifier
	Body   *BlockStatment
}

func (fl *FunctionLiteral) expressionNode()      {}
func (fl *FunctionLiteral) TokenLiteral() string { return fl.Token.Value }
func (fl *FunctionLiteral) Stringify() string {
	var out bytes.Buffer
	params := []string{}
	for _, p := range fl.Params {
		params = append(params, p.Stringify())
	}
	out.WriteString(fl.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString(fl.Body.Stringify())
	return out.String()
}

type CallExpression struct {
	Token     lexer.Token
	Function  Expression
	Arguments []Expression
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Value }
func (ce *CallExpression) Stringify() string {
	var out bytes.Buffer
	args := []string{}
	for _, a := range ce.Arguments {
		args = append(args, a.Stringify())
	}
	out.WriteString(ce.Function.Stringify())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")
	return out.String()
}
