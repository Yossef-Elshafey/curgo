package parser

import (
	"curgo/lexer"
	"curgo/types/ast"
	"curgo/types/tokens"
	"curgo/utils"
	"fmt"
	"log"
)

type Parser struct {
	currentToken  lexer.Token
	tokens        []lexer.Token
	peekToken     lexer.Token
	position      int
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
	"*": PRODUCT,
	"/": PRODUCT,
}

var prefixLookup = map[tokens.TokenKind]prefixParseFn{}
var infixLookup = map[tokens.TokenKind]infixParseFn{}


func Parse(t []lexer.Token) ast.Program {
	p := &Parser{}
	p.tokens = t
	p.initParser()
	program := ast.Program{}
	for !p.peekTokenIs(tokens.EOF) {
		stmt := p.parseStmt()
		program.Statements = append(program.Statements, stmt)
		p.advanceTokens()
	}
	return program
}

func (p *Parser) initParser() {
	p.alignTokens()
	p.initPrefix()
	p.initInfix()
}

func (p *Parser) initPrefix() {
	p.registerPrefix(tokens.IDENTIFIER, p.parseIdentifier)
	p.registerPrefix(tokens.STRING,     p.parseStringLiteral)
	p.registerPrefix(tokens.BACKTICK,   p.parseStringLiteral)
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
	p.peekToken = p.tokens[p.position]
	if p.position+1 != len(p.tokens) {
		p.position++
	}
}

func (p *Parser) alignTokens() {
	p.advanceTokens()
	p.advanceTokens()
}

func (p *Parser) isEndOfFetch() {
	if p.peekTokenIs(tokens.EOF) && p.currentToken.Kind == tokens.SEMI_COLON {
		fmt.Printf("hint: use 'endfet' keyword to close fetch statement\n")
	}
}

func (p *Parser) expectPeekToBe(k tokens.TokenKind) bool {
	if p.peekToken.Kind != k {
		line := p.peekToken.Pos.Line
		lineIssue := utils.ReadSourceAsLines(line)
		p.isEndOfFetch()

    // fmt.Printf("%c[%dmHELLO!\n", 0x1B, 32);
		fmt.Printf("Parser:%d:%d: encouter error at line:\n %s\n", line-1,
			p.currentToken.Pos.End,
			lineIssue)

		log.Fatalf("Expect to find %s after '%s', got=%s", tokens.TokenKindStringify(k),
			p.currentToken.Value,
			tokens.TokenKindStringify(p.peekToken.Kind))
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
	fmt.Printf("Infix is not supported\n") // NOTE:
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

	if !p.expectPeekToBe(tokens.STRING) {
		return nil
	}

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
