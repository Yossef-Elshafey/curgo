package parser

import (
	"log"
	"parser/ast"
	"parser/lexer"
	"strconv"
)

type bindingPower int

const (
	// NOTE: enum order is related to higher precedence (opposite)
	defalt bindingPower = iota
	comma
	assignment
	logical
	relational
	additive
	multiplicative
	unary
	call
	member
	primary
)

type stmtHandler func(p *parser) ast.Stmt

type nudHandler func(p *parser) ast.Expr

type ledHandler func(p *parser, left ast.Expr, bp bindingPower) ast.Expr

type stmtLookup map[lexer.TokenKind]stmtHandler
type nudLookUp map[lexer.TokenKind]nudHandler
type ledLookup map[lexer.TokenKind]ledHandler
type bpLookUp map[lexer.TokenKind]bindingPower

var stmtLu = stmtLookup{}
var nudLu = nudLookUp{}
var ledLu = ledLookup{}
var bpLu = bpLookUp{}

func led(t lexer.TokenKind, bp bindingPower, ledFn ledHandler) {
	bpLu[t] = bp
	ledLu[t] = ledFn
}

func nud(t lexer.TokenKind, bp bindingPower, nudFn nudHandler) {
	bpLu[t] = bp
	nudLu[t] = nudFn
}

func stmt(t lexer.TokenKind, stmtFn stmtHandler) {
	bpLu[t] = defalt
	stmtLu[t] = stmtFn
}

func createLookups() {
	// -------------- logical
	led(lexer.AND, logical, parseBinary)
	led(lexer.OR, logical, parseBinary)
	// --------------

	// -------------- bit
	led(lexer.LESS, relational, parseBinary)
	led(lexer.LESS_EQUALS, relational, parseBinary)
	led(lexer.GREATER, relational, parseBinary)
	led(lexer.GREATER_EQUALS, relational, parseBinary)
	led(lexer.EQUALS, relational, parseBinary)
	led(lexer.NOT_EQUALS, relational, parseBinary)
	// --------------

	// -------------- arithmetic
	led(lexer.PLUS, additive, parseBinary)
	led(lexer.DASH, additive, parseBinary)

	led(lexer.STAR, multiplicative, parseBinary)
	led(lexer.PERCENT, multiplicative, parseBinary)
	led(lexer.SLASH, multiplicative, parseBinary)
	// --------------

	// -------------- primary
	nud(lexer.NUMBER, primary, parsePrimary)
	nud(lexer.STRING, primary, parsePrimary)
	nud(lexer.IDENTIFIER, primary, parsePrimary)
	// --------------
}

func parsePrimary(p *parser) ast.Expr {
	switch p.currentToken().Type {
	case lexer.NUMBER:
		number, err := strconv.ParseFloat(p.eat().Value, 64)
		if err != nil {
			log.Fatalf("Parser: cannot convert %s to a number", p.eat().Value)
		}
		return ast.NumberExpr{Value: number}

	case lexer.STRING:
		return ast.StringExpr{
			Value: p.eat().Value,
		}
	case lexer.IDENTIFIER:
		return ast.SymbolExpr{
			Value: p.eat().Value,
		}
	default:
		log.Fatalf("Cannot build expression %s", lexer.TokenKindString(p.currentToken().Type))
		return ast.SymbolExpr{}
	}
}

func parseBinary(p *parser, left ast.Expr, bp bindingPower) ast.Expr {
	operator := p.eat()
	rhs := parseExpr(p, bp)
	return ast.BinaryExpr{
		Left:     left,
		Operator: operator,
		Right:    rhs,
	}
}

func parseExpr(p *parser, bp bindingPower) ast.Expr {
	token := p.currentToken()
	nudFn, exists := nudLu[token.Type]
	if !exists {
		log.Fatalf("Parser: cannot find handler for token %s", lexer.TokenKindString(token.Type))
	}
	lhs := nudFn(p)
	for bpLu[p.currentToken().Type] > bp {
		t := p.currentToken().Type
		ledFn, exists := ledLu[t]

		if !exists {
			log.Fatalf("Parser: cannot find handler for token %s", lexer.TokenKindString(token.Type))
		}

		// Traverse
		lhs = ledFn(p, lhs, bp)
	}
	return lhs
}
