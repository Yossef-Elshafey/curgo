package parser

import (
	"curgo/types/tokens"
	"fmt"
)

type Parser struct {
	que          []tokens.Token
	pos          int
	currentToken tokens.Token
	peekToken    tokens.Token
}

func NewParser() *Parser {
	p := &Parser{}
	p.pos = 0
	p.que = make([]tokens.Token, 0)
	return p
}

func (p *Parser) Parse(token <-chan tokens.Token) {
	for tk := range token {
		fmt.Printf("Reading: %+v\n", tk)
		p.que = append(p.que, tk)
		if len(p.que) >= 2 {
			fmt.Printf("Begin Processing\n")
			p.process()
		}
	}
}

func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	p.peekToken = p.que[p.pos]
	if p.pos+1 != len(p.que) {
		p.pos++
	}
}

func (p *Parser) process() {
	p.nextToken()
	p.nextToken()
	fmt.Printf("Current: %+v\n", p.currentToken)
	fmt.Printf("Peek: %+v\n", p.peekToken)
}
