package parser

import (
	"curgo/lexer"
	"curgo/types/ast"
	"curgo/types/tokens"
	"log"
	"strconv"
)

type Parser struct {
	tokens        *lexer.Lexer
	currentToken  token.Token
	peekToken     token.Token
	prefixLookup  map[token.TokenKind]prefixParseFn
	infixLookup   map[token.TokenKind]infixParseFn
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

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		tokens:      l,
	}
	p.initInfix()
	p.initPrefix()
	p.advanceTokens()
	p.advanceTokens()
	return p
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for !p.curTokenIs(token.EOF) {
		stmt := p.parseStmt()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.advanceTokens()
	}

	return program
}

func (p *Parser) initPrefix() {
	if p.prefixLookup == nil {
		p.prefixLookup = make(map[token.TokenKind]prefixParseFn)
	}
	p.registerPrefix(token.IDENTIFIER,  p.parseIdentifier)
	p.registerPrefix(token.STRING,      p.parseStringLiteral)
	p.registerPrefix(token.NUMBER,      p.parseNumberLiteral)
	p.registerPrefix(token.LPAREN,      p.parseGroupedExpression)
}

func (p *Parser) initInfix() {
	if p.infixLookup == nil {
		p.infixLookup = make(map[token.TokenKind]infixParseFn)
	}
	p.registerInfix(token.PLUS, p.parseBinaryExpression)
	p.registerInfix(token.LPAREN, p.parseCallExpression)
}

func (p *Parser) registerPrefix(k token.TokenKind, handler prefixParseFn) {
	p.prefixLookup[k] = handler
}

func (p *Parser) registerInfix(k token.TokenKind, handler infixParseFn) {
	p.infixLookup[k] = handler
}

func (p *Parser) peekTokenBindingPower() bindingPower {
	return bindingPowerLookup[p.peekToken.Value]
}

func (p *Parser) currentTokenBindingPower() bindingPower {
	return bindingPowerLookup[p.currentToken.Value]
}

func (p *Parser) peekTokenIs(k token.TokenKind) bool {
	if p.peekToken.Kind != k {
		return false
	}
	return true
}

func (p *Parser) curTokenIs(t token.TokenKind) bool {
	return p.currentToken.Kind == t
}

func (p *Parser) advanceTokens() {
	p.currentToken = p.peekToken
	p.peekToken = p.tokens.NextToken()
}

func (p *Parser) expectPeekToBe(k token.TokenKind) bool {
	if !p.peekTokenIs(k) { return false }
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
	case token.FETCH:
		return p.parseFetchStatment()
	case token.LET:
		return p.parseLetStmt()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseLetStmt() *ast.LetStatement {
	if !p.expectPeekToBe(token.IDENTIFIER) {return nil}
	ls := &ast.LetStatement{}
	ls.Identifier = &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Value}
	if !p.expectPeekToBe(token.ASSIGN) { return nil }
	p.advanceTokens()
	ls.Value = p.parseExpression(LOWEST)
	if !p.expectPeekToBe(token.SEMICOLON) {return nil}
	return ls
}

func (p *Parser) parseFetchStatment() *ast.FetchStmt {
	fs := &ast.FetchStmt{Token: p.currentToken}
	fs.Body = []ast.Statement{}
	fs.Arguments = []*ast.Identifier{}

	if !p.expectPeekToBe(token.IDENTIFIER) {
		return nil
	}

	fs.FetchIdentifier = &ast.Identifier{
		Token: p.currentToken,
		Value:  p.currentToken.Value,
	}

	p.advanceTokens()
	fs.Arguments = p.parseFetchArguments()

	// if !p.curTokenIs(token.RPAREN) { return nil }
	if !p.expectPeekToBe(token.COLON) {return nil}

	// TODO: if there is no token(endfet) its an infinite loop
	for !p.peekTokenIs(token.ENDFETCH) {
		fs.Body = append(fs.Body, p.parseFetchBody())
	}

	p.advanceTokens()
	return fs
}

func (p *Parser) parseFetchArguments() []*ast.Identifier {
	args := []*ast.Identifier{}
	if p.peekTokenIs(token.RPAREN) {
		p.advanceTokens()
		return args
	}
	p.advanceTokens()
	arg := &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Value}

	args = append(args, arg)
	for p.peekTokenIs(token.COMMA) {
		p.advanceTokens()
		p.advanceTokens()
		arg = &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Value}
		args = append(args, arg)
	}
	if !p.expectPeekToBe(token.RPAREN) { return nil }
	return args
}

func (p *Parser) parseFetchBody() ast.Statement {
	if !p.expectPeekToBe(token.IDENTIFIER) {
		return nil
	}
	ca := &ast.CurgoAssignStatment{}
	ca.Arg = &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Value}
	if !p.expectPeekToBe(token.TRANSPILEASSIGN) {
		return nil
	}
	p.advanceTokens()
	ca.Value = p.parseExpression(LOWEST)

	if !p.expectPeekToBe(token.SEMICOLON) {
		return nil
	}
	return ca
}

func (p *Parser) parseExpression(bp bindingPower) ast.Expression {
	prefix := p.prefixLookup[p.currentToken.Kind]
	if prefix == nil {
		// p.noPrefixFoundErr(p.currentToken) TODO:
		log.Fatalf("Prefix not founded: %s\n", p.currentToken.Value)
		return nil
	}
	left := prefix()

	for !p.peekTokenIs(token.SEMICOLON) && bp < p.peekTokenBindingPower() {
		infix := p.infixLookup[p.peekToken.Kind]
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
	if p.peekTokenIs(token.RPAREN) {
		p.advanceTokens()
		return args
	}
	p.advanceTokens()
	args = append(args, p.parseExpression(LOWEST))
	for p.peekTokenIs(token.COMMA) {
		p.advanceTokens()
		p.advanceTokens()
		args = append(args, p.parseExpression(LOWEST))
	}
	if !p.expectPeekToBe(token.RPAREN) { return nil }
	return args
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	es := &ast.ExpressionStatement{}
	es.Expression = p.parseExpression(LOWEST)
	if p.peekTokenIs(token.SEMICOLON) { p.advanceTokens() }
	return es
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.advanceTokens()
	exp := p.parseExpression(LOWEST)
	if !p.expectPeekToBe(token.RPAREN) { return nil }
	return exp
}
