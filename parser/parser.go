package parser

import (
	"curgo/lexer"
	"curgo/types/ast"
	"curgo/types/tokens"
	"curgo/utils"
	"fmt"
	"log"
	"strconv"
)

type Parser struct {
	currentToken  lexer.Token
	tokens        []lexer.Token
	peekToken     lexer.Token
	position      int
}

type (
	bindingPower   int
	infixParseFn   func(ast.Expression)  ast.Expression
	prefixParseFn  func()                ast.Expression
)

const (
	LOWEST bindingPower = iota
	SUM
	PRODUCT
	CALL
)

var bindingPowerLookup = map[string]bindingPower{
	"+": SUM,
	"-": SUM,
	"*": PRODUCT,
	"/": PRODUCT,
	"(": CALL,
}

var prefixLookup = map[tokens.TokenKind]prefixParseFn{}
var infixLookup = map[tokens.TokenKind]infixParseFn{}


func Parse(t []lexer.Token) *ast.Program {
	p := &Parser{}
	p.tokens = t
	p.initParser()
	program := &ast.Program{}
	for !p.peekTokenIs(tokens.EOF) {
		stmt := p.parseStmt()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
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
	p.registerPrefix(tokens.IDENTIFIER,  p.parseIdentifier)
	p.registerPrefix(tokens.STRING,      p.parseStringLiteral)
	p.registerPrefix(tokens.BACKTICK,    p.parseStringLiteral)
	p.registerPrefix(tokens.NUMBER,      p.parseNumberLiteral)
	p.registerPrefix(tokens.OPEN_PAREN,  p.parseGroupedExpression)
}

func (p *Parser) initInfix() {
	p.registerInfix(tokens.PLUS, p.parseBinaryExpression)
	p.registerInfix(tokens.OPEN_PAREN, p.parseCallExpression)
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

func (p *Parser) currentTokenBindingPower() bindingPower {
	return bindingPowerLookup[p.currentToken.Value]
}

// NOTE: peekTokenIs, expectPeekToBe the way peek check handled is foolish
// implement peek token check such that errors, token advancing is clear
// Seperate the error messages in different function
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
		fmt.Printf("Parser:%d:%d: encouter error at line:\n %s\n", line-1,
			p.currentToken.Pos.End,
			lineIssue)

		log.Fatalf("Expect to find %s after '%s', got=%s", tokens.TokenKindStringify(k),
			p.currentToken.Value,
			tokens.TokenKindStringify(p.peekToken.Kind))

		return false
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
	expr := &ast.BinaryExpression{
		Operator: p.currentToken,
		Left: lhs,
	}
	bp := p.currentTokenBindingPower()
	p.advanceTokens()
	expr.Right = p.parseExpression(bp)
	return expr
}

func (p *Parser) parseStmt() ast.Statement {
	switch p.currentToken.Kind {
	case tokens.FETCH:
		return p.parseFetchStatment()
	case tokens.LET:
		return p.parseLetStmt()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseLetStmt() *ast.LetStatement {
	if !p.expectPeekToBe(tokens.IDENTIFIER) {return nil}
	ls := &ast.LetStatement{}
	ls.Identifier = &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Value}
	if !p.expectPeekToBe(tokens.EQUAL) { }
	p.advanceTokens()
	ls.Value = p.parseExpression(LOWEST)
	if !p.expectPeekToBe(tokens.SEMI_COLON) {return nil}
	return ls
}

func (p *Parser) parseFetchStatment() *ast.FetchStmt {
	fs := &ast.FetchStmt{Token: p.currentToken}
	fs.Body = []ast.Statement{}
	fs.Arguments = []*ast.Identifier{}

	if !p.expectPeekToBe(tokens.IDENTIFIER) {
		return nil
	}

	fs.FetchIdentifier = &ast.Identifier{
		Token: p.currentToken,
		Value:  p.currentToken.Value,
	}

	p.advanceTokens()
	fs.Arguments = p.parseFetchArguments()

	if !p.expectPeekToBe(tokens.COLON) { return nil }

	for !p.peekTokenIs(tokens.ENDFETCH) {
		fs.Body = append(fs.Body, p.parseFetchBody())
	}

	p.advanceTokens()
	return fs
}

func (p *Parser) parseFetchArguments() []*ast.Identifier {
	args := []*ast.Identifier{}
	if p.peekTokenIs(tokens.CLOSE_PAREN) {
		p.advanceTokens()
		return args
	}
	p.advanceTokens()
	arg := &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Value}

	args = append(args, arg)
	for p.peekTokenIs(tokens.COMMA) {
		p.advanceTokens()
		p.advanceTokens()
		arg = &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Value}
		args = append(args, arg)
	}
	if !p.expectPeekToBe(tokens.CLOSE_PAREN) { return nil }
	return args
}

func (p *Parser) parseFetchBody() ast.Statement {
	if !p.expectPeekToBe(tokens.STRING) {
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
	if prefix == nil {
		// p.noPrefixFoundErr(p.currentToken) TODO:
		log.Fatalf("Prefix not founded: %s\n", p.currentToken.Value)
		return nil
	}
	left := prefix()

	for !p.peekTokenIs(tokens.SEMI_COLON) && bp < p.peekTokenBindingPower() {
		infix := infixLookup[p.peekToken.Kind]
		if infix == nil {
			return left
		}
		p.advanceTokens()
		left = infix(left)
	}
	return left
}

func (p *Parser) parseNumberLiteral() ast.Expression {
	nl := &ast.NumberLiteral{}
	i, err := strconv.Atoi(p.currentToken.Value)
	if err != nil {
		log.Fatalf("Failed to convert %v to string", p.currentToken.Value)
	}
	nl.Value = int64(i)
	nl.Token = p.currentToken
	return nl
}

func (p *Parser) parseCallExpression(fs ast.Expression) ast.Expression {
	ce := &ast.CallExpression{Token: p.currentToken, Function: fs}
	ce.Arguments = p.parseCallArgument()
	return ce
}

func (p *Parser) parseCallArgument() []ast.Expression {
	args := []ast.Expression{}
	if p.peekTokenIs(tokens.CLOSE_PAREN) {
		p.advanceTokens()
		return args
	}
	p.advanceTokens()
	args = append(args, p.parseExpression(LOWEST))
	for p.peekTokenIs(tokens.COMMA) {
		p.advanceTokens()
		p.advanceTokens()
		args = append(args, p.parseExpression(LOWEST))
	}
	if !p.expectPeekToBe(tokens.CLOSE_PAREN) { return nil }
	return args
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	es := &ast.ExpressionStatement{}
	es.Expression = p.parseExpression(LOWEST)
	if p.peekTokenIs(tokens.SEMI_COLON) { p.advanceTokens() }
	return es
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.advanceTokens()
	exp := p.parseExpression(LOWEST)
	if !p.expectPeekToBe(tokens.CLOSE_PAREN) { return nil }
	return exp
}
