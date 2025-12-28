package parser

import (
	"curgo/eval"
	"curgo/types/ast"
	"curgo/types/tokens"
	"fmt"
	"log"
)

type Parser struct {
	Tokens       []tokens.Token
	currentToken tokens.Token
	peekToken    tokens.Token
	position     int
	evaluater    *eval.Evaluater
}

type (
	bindingPower  int
	infixParseFn  func(ast.Expression) ast.Expression
	prefixParseFn func() ast.Expression
)

const (
	LOWEST bindingPower = iota
	SUM
	PRODUCT
)

// NOTE: prefix and infix in not used currently since there is no response(curl output) operations

var bindingPowerLookup = map[string]bindingPower{
	"+": SUM,
	"-": SUM,
	"*": SUM,
	"/": SUM,
}
var prefixLookup = map[tokens.TokenKind]prefixParseFn{}
var infixLookup = map[tokens.TokenKind]infixParseFn{}

func (p *Parser) initParser() {
	p.alignTokens()
	p.initPrefix()
	p.initInfix()
	p.evaluater = &eval.Evaluater{}
}

func (p *Parser) Parse() *eval.Evaluater {
	p.initParser()
	program := ast.Program{}
	for !p.peekTokenIs(tokens.EOF) {
		stmt := p.parseStmt()
		program.Statements = append(program.Statements, stmt)
		p.advanceTokens()
	}
	p.evaluater.Program = program
	return p.evaluater
}

func (p *Parser) initPrefix() {
	p.registerPrefix(tokens.IDENTIFIER, p.parseIdentifier)
	p.registerPrefix(tokens.STRING, p.parseStringLiteral)
	p.registerPrefix(tokens.BACKTICK, p.parseStringLiteral)
}

func (p *Parser) initInfix() {
	p.registerInfix(tokens.PLUS, p.parseBinaryExpression)
}

func (p *Parser) registerPrefix(k tokens.TokenKind, handler prefixParseFn) {
	prefixLookup[k] = handler
}

func (p *Parser) registerInfix(k tokens.TokenKind, handler infixParseFn) {
	infixLookup[k] = handler
}

func (p *Parser) peekTokenBindingPower() bindingPower {
	return bindingPowerLookup[p.peekToken.Value]
}

func (p *Parser) peekTokenIs(k tokens.TokenKind) bool {
	if p.peekToken.Kind != k {
		return false
	}
	return true
}

func (p *Parser) advanceTokens() {
	p.currentToken = p.peekToken
	p.peekToken = p.Tokens[p.position]
	if p.position+1 != len(p.Tokens) {
		p.position++
	}
}

func (p *Parser) alignTokens() {
	p.advanceTokens()
	p.advanceTokens()
}

func (p *Parser) expectPeekToBe(k tokens.TokenKind) bool {
	if p.peekToken.Kind != k {
		log.Fatalf("Parser:%d:%d: Expected %s to be %s",
			p.peekToken.Pos.Line,
			p.peekToken.Pos.Start,
			tokens.TokenKindStringify(p.peekToken.Kind),
			tokens.TokenKindStringify(k))
	}
	p.advanceTokens()
	return true
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Value}
}

func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{Token: p.currentToken, Value: p.currentToken.Value}
}

func (p *Parser) parseBinaryExpression(lhs ast.Expression) ast.Expression {
	fmt.Printf("Infix is not supported\n")
	p.advanceTokens()
	return nil
}

func (p *Parser) parseStmt() ast.Statement {
	switch p.currentToken.Kind {
	case tokens.FETCH:
		return p.parseFetchStatment()
	default:
		return nil
	}
}

func (p *Parser) parseFetchStatment() *ast.FetchStmt {
	fs := &ast.FetchStmt{Token: p.currentToken}
	fs.Body = []ast.Statement{}
	if !p.expectPeekToBe(tokens.IDENTIFIER) {
		return nil
	}

	fs.FetchIdentifier = &ast.Identifier{
		Token: p.currentToken,
		Value:  p.currentToken.Value,
	}

	if !p.expectPeekToBe(tokens.COLON) {
		return nil
	}

	for !p.peekTokenIs(tokens.ENDFETCH) {
		fs.Body = append(fs.Body, p.parseFetchBody())
	}
	p.advanceTokens()
	return fs
}

func (p *Parser) parseFetchBody() ast.Statement {
	if !p.expectPeekToBe(tokens.IDENTIFIER) {
		return nil
	}
	ca := &ast.CurgoAssignStatment{}
	ca.Arg = &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Value}
	if !p.expectPeekToBe(tokens.TRANSPILE_ASSIGN) {
		return nil
	}

	p.advanceTokens()

	ca.Value = p.parseExpression(LOWEST)

	if !p.expectPeekToBe(tokens.SEMI_COLON) {
		return nil
	}
	return ca
}

func (p *Parser) parseExpression(bp bindingPower) ast.Expression {
	prefix := prefixLookup[p.currentToken.Kind]
	left := prefix()

	for !p.peekTokenIs(tokens.SEMI_COLON) && bp < p.peekTokenBindingPower() {
		infix := infixLookup[p.peekToken.Kind]
		left = infix(left)
	}

	return left
}
