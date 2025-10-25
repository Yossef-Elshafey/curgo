package parser

import (
	"curgo/ast"
	"curgo/lexer"
	"fmt"
	"strconv"
)

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

type bindingPower int

const (
	_ bindingPower = iota
	LOWEST
	EQUALS
	NOT_EQUALS
	LESSGREATER
	SUM
	PRODUCT
	PREFIX
	CALL
)

var bindingPowerLookup = map[lexer.TokenKind]bindingPower{
	lexer.EQUALS:     EQUALS,
	lexer.NOT_EQUALS: NOT_EQUALS,
	lexer.LESS:       LESSGREATER,
	lexer.GREATER:    LESSGREATER,
	lexer.PLUS:       SUM,
	lexer.DASH:       SUM,
	lexer.SLASH:      PRODUCT,
	lexer.STAR:       PRODUCT,
}

type Parser struct {
	l              []lexer.Token
	curToken       lexer.Token
	peekToken      lexer.Token
	pos            int
	errors         []string
	prefixParseFns map[lexer.TokenKind]prefixParseFn
	infixParseFns  map[lexer.TokenKind]infixParseFn
}

func New(l []lexer.Token) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	// set the two token pointers appropriately
	p.nextToken()
	p.nextToken()
	p.initPrefix()
	p.initInfix()
	return p
}

func (p *Parser) initPrefix() {
	p.prefixParseFns = make(map[lexer.TokenKind]prefixParseFn)
	p.registerPrefix(lexer.IDENTIFIER, p.parseIdentifier)
	p.registerPrefix(lexer.NUMBER, p.parseIntegerLiteral)
	p.registerPrefix(lexer.NOT, p.parsePrefixExpression)
	p.registerPrefix(lexer.DASH, p.parsePrefixExpression)
}

func (p *Parser) initInfix() {
	p.infixParseFns = make(map[lexer.TokenKind]infixParseFn)
	p.registerInfix(lexer.PLUS, p.parseBinaryExpression)
	p.registerInfix(lexer.MINUS_MINUS, p.parseBinaryExpression)
	p.registerInfix(lexer.SLASH, p.parseBinaryExpression)
	p.registerInfix(lexer.STAR, p.parseBinaryExpression)
	p.registerInfix(lexer.EQUALS, p.parseBinaryExpression)
	p.registerInfix(lexer.NOT_EQUALS, p.parseBinaryExpression)
	p.registerInfix(lexer.LESS, p.parseBinaryExpression)
	p.registerInfix(lexer.GREATER, p.parseBinaryExpression)
}

func (p *Parser) registerPrefix(kind lexer.TokenKind, fn prefixParseFn) {
	p.prefixParseFns[kind] = fn
}

func (p *Parser) peekBindingPower() bindingPower {
	if p, ok := bindingPowerLookup[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) currentBindingPower() bindingPower {
	if p, ok := bindingPowerLookup[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) registerInfix(kind lexer.TokenKind, fn infixParseFn) {
	p.infixParseFns[kind] = fn
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) PeekError(t lexer.TokenKind) {
	msg := fmt.Sprintf("Expected next token to be %s, got %s instead ",
		lexer.TokenKindString(t), lexer.TokenKindString(p.peekToken.Type))
	p.errors = append(p.errors, msg)
}

func (p *Parser) nextToken() {
	fmt.Printf("Got call To Advance Tokens. before\n Current:%+v, Peek:%+v\n", p.curToken, p.peekToken)
	p.curToken = p.peekToken
	p.peekToken = p.l[p.pos]
	if p.pos+1 != len(p.l) {
		p.pos++
	}
	fmt.Printf("After\n Current:%+v, Peek:%+v\n", p.curToken, p.peekToken)
}

func (p *Parser) curTokenIs(t lexer.TokenKind) bool {
	return p.curToken.Type == t
}

func (p *Parser) PeekTokenIs(t lexer.TokenKind) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t lexer.TokenKind) bool {
	if p.PeekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.PeekError(t)
		return false
	}
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}
	fmt.Printf("Parsing Tokens: %+v\n", p.l)
	for p.curToken.Type != lexer.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
	return program
}

func (p *Parser) parseBinaryExpression(left ast.Expression) ast.Expression {
	exp := &ast.BinaryExpression{
		Token:    p.curToken,
		Operator: p.curToken.Value,
		Left:     left,
	}
	bp := p.currentBindingPower()
	p.nextToken()
	exp.Right = p.parseExpression(bp)
	return exp
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	exp := &ast.UnaryExpression{
		Token:    p.curToken,
		Operator: p.curToken.Value,
	}
	p.nextToken()
	fmt.Printf("Current Token after parsePrefixExpression advance: %+v\n", p.curToken)
	exp.Right = p.parseExpression(PREFIX)
	return exp
}

func (p *Parser) noPrefixParseFnError(t lexer.TokenKind) {
	msg := fmt.Sprintf("no prefix parse function for %s found", lexer.TokenKindString(t))
	p.errors = append(p.errors, msg)
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Value}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}
	value, err := strconv.ParseInt(p.curToken.Value, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Value)
		p.errors = append(p.errors, msg)
		return nil
	}
	lit.Value = value
	return lit
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case lexer.LET:
		return p.parseLetStatement()
	case lexer.RETURN:
		return p.parserReturnStatement()
	default:
		return p.parseExpressionStatment()
	}
}

func (p *Parser) parseExpressionStatment() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}
	stmt.Expression = p.parseExpression(LOWEST)
	if p.PeekTokenIs(lexer.SEMI_COLON) {
		p.nextToken()
	}
	// TODO: else throw error
	return stmt
}

func (p *Parser) parseExpression(bp bindingPower) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}

	leftExp := prefix()
	return leftExp
}

func (p *Parser) parserReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}
	p.nextToken()
	for !p.curTokenIs(lexer.SEMI_COLON) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseLetStatement() *ast.LetStatment {
	stmt := &ast.LetStatment{Token: p.curToken}
	if !p.expectPeek(lexer.IDENTIFIER) {
		return nil
	}

	stmt.Name = &ast.Identifier{
		Token: p.curToken,
		Value: p.curToken.Value,
	}

	if !p.expectPeek(lexer.ASSIGNMENT) {
		return nil
	}
	// Skip parsing the expression for now
	for !p.curTokenIs(lexer.SEMI_COLON) {
		p.nextToken()
	}
	return stmt
}
