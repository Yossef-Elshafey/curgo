package main

import (
	"bytes"
	"fmt"
	"log"
	"strconv"
	"strings"
)

func main() {
	p := New("1+1+1+1+1*2;")
	p.advanceTokens()
	p.advanceTokens()
	p.printTokens("Main")
	program := p.parseExpression(LOWEST)
	fmt.Printf("Expression: %+v\n", program.(*BinaryExpression).stringify())
	fmt.Printf("%s\n", strings.Repeat("--", 60))
	ev := Eval(program)
	fmt.Printf("Evaluate: %d\n", ev.(*Number).value)
}

func Eval(program Expression) Expression {
	switch node := program.(type) {
	case *Number:
		return &Number{value: node.value}
	case *BinaryExpression:
		left := Eval(node.Left)
		right := Eval(node.Right)
		return evalBinaryExpression(node.Operator, left, right)
	}
	return nil
}

func evalBinaryExpression(op string, l Expression, r Expression) Expression {
	switch op {
	case "+":
		n1, n2 := l.(*Number), r.(*Number)
		return &Number{n1.value + n2.value}

	case "*":
		n1, n2 := l.(*Number), r.(*Number)
		return &Number{n1.value * n2.value}
	}
	return nil
}

const (
	LOWEST = iota
	SUM
	PRODUCT
)

var bindingPowers = map[string]int{
	"+": SUM,
	"*": PRODUCT,
}

type Expression interface {
	expressionNode()
	stringify() string
}

type Number struct {
	value int
}

func (n *Number) expressionNode() {}
func (n *Number) stringify() string {
	var node bytes.Buffer
	node.WriteString(strconv.Itoa(n.value))
	return node.String()
}

type BinaryExpression struct {
	Left     Expression
	Operator string
	Right    Expression
}

func (b *BinaryExpression) expressionNode() {}
func (b *BinaryExpression) stringify() string {
	var node bytes.Buffer
	if _, ok := b.Left.(*BinaryExpression); ok {
		node.WriteString("Left ")
	} else {
		node.WriteString(b.Left.stringify())
	}

	node.WriteString(b.Operator)
	if _, ok := b.Right.(*BinaryExpression); ok {
		node.WriteString(" Right")
	} else {
		node.WriteString(b.Right.stringify())
	}
	return node.String()
}

type Parser struct {
	input        string
	currentToken string
	peekToken    string
	position     int
}

func New(i string) *Parser {
	return &Parser{input: i}
}

func (t *Parser) advanceTokens() {
	t.currentToken = t.peekToken
	t.peekToken = string(t.input[t.position])
	if t.position+1 < len(t.input) {
		t.position += 1
	}
}

func (p *Parser) printTokens(prefix string) {
	fmt.Printf("%s, Current: %+v, Peek: %+v\n", prefix, p.currentToken, p.peekToken)
}

func (p *Parser) parseExpression(bp int) Expression {
	var lhs Expression
	for p.currentToken != ";" {
		lhs = p.handlePrefix()
		p.printTokens("bp < p.peekBindingPower")
		for bp < p.peekBindingPower() {
			p.advanceTokens()
			lhs = p.parseBinaryExpression(lhs)
		}
		return lhs
	}
	return lhs //
}

func (p *Parser) peekBindingPower() int {
	return bindingPowers[p.peekToken]
}

func (p *Parser) parseBinaryExpression(lhs Expression) Expression {
	exp := &BinaryExpression{
		Left:     lhs,
		Operator: p.currentToken,
	}
	currentBindingPower := bindingPowers[p.currentToken]
	p.advanceTokens()
	exp.Right = p.parseExpression(currentBindingPower)
	return exp
}

func (p *Parser) handlePrefix() Expression {
	num, error := strconv.Atoi(p.currentToken)
	if error != nil {
		log.Fatalf("Cannot convert Atoi got=%s", p.currentToken)
	}
	return &Number{value: num}
}
