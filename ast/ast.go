package ast

import (
	"bytes"
	"curgo/lexer"
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
// Then i can use lexer.TokenKindStr(token) to represent its value as string

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
