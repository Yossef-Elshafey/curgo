package ast

import "curgo/types/tokens"

type Expression interface {
	expressionNode()
}

type Statement interface {
	statementNode()
}

type Program struct {
	Stmts []Statement
}

type LetStatment struct {
	Token tokens.Token
	Name  *Identifier
	Value string
}

func (ls *LetStatment) statementNode() {}

type Identifier struct {
	Token tokens.Token
	Value string
}

func (i *Identifier) expressionNode() {}
