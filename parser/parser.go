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
	lexer.OPEN_PAREN: CALL,
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
	p.registerPrefix(lexer.TRUE, p.parseBoolean)
	p.registerPrefix(lexer.FALSE, p.parseBoolean)
	p.registerPrefix(lexer.OPEN_PAREN, p.parseGroupedExpression)
	p.registerPrefix(lexer.IF, p.parseIfExpression)
	p.registerPrefix(lexer.FN, p.parseFunctionLiteral)
}

func (p *Parser) initInfix() {
	p.infixParseFns = make(map[lexer.TokenKind]infixParseFn)
	p.registerInfix(lexer.PLUS, p.parseBinaryExpression)
	p.registerInfix(lexer.DASH, p.parseBinaryExpression)
	p.registerInfix(lexer.STAR, p.parseBinaryExpression)
	p.registerInfix(lexer.SLASH, p.parseBinaryExpression)
	p.registerInfix(lexer.EQUALS, p.parseBinaryExpression)
	p.registerInfix(lexer.NOT_EQUALS, p.parseBinaryExpression)
	p.registerInfix(lexer.LESS, p.parseBinaryExpression)
	p.registerInfix(lexer.GREATER, p.parseBinaryExpression)
	p.registerInfix(lexer.OPEN_PAREN, p.parseCallExpression)
}

func (p *Parser) noPrefixParseFnError(t lexer.TokenKind) {
	msg := fmt.Sprintf("no prefix parse function for %s found", lexer.TokenKindString(t))
	p.errors = append(p.errors, msg)
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
	if p, ok := bindingPowerLookup[p.curToken.Type]; ok {
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
	p.curToken = p.peekToken
	p.peekToken = p.l[p.pos]
	if p.pos+1 != len(p.l) {
		p.pos++
	}
}

func (p *Parser) curTokenIs(t lexer.TokenKind) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t lexer.TokenKind) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t lexer.TokenKind) bool {
	if p.peekTokenIs(t) {
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
	for p.curToken.Type != lexer.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
	return program
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

	p.nextToken()
	stmt.Value = p.parseExpression(LOWEST)

	for p.peekTokenIs(lexer.SEMI_COLON) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parserReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}
	p.nextToken()
	stmt.ReturnValue = p.parseExpression(LOWEST)
	for p.peekTokenIs(lexer.SEMI_COLON) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseExpressionStatment() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}
	stmt.Expression = p.parseExpression(LOWEST)
	if p.peekTokenIs(lexer.SEMI_COLON) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseExpression(bp bindingPower) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}

	leftExp := prefix() // lhs

	for !p.peekTokenIs(lexer.SEMI_COLON) && bp < p.peekBindingPower() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}
		p.nextToken()
		leftExp = infix(leftExp)
	}
	return leftExp
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	lit := &ast.FunctionLiteral{Token: p.curToken}
	if !p.expectPeek(lexer.OPEN_PAREN) {
		return nil
	}
	lit.Params = p.parseFunctionParams()
	if !p.expectPeek(lexer.OPEN_CURLY) {
		return nil
	}
	lit.Body = p.parseBlockStatement()
	return lit
}

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	exp := &ast.CallExpression{Token: p.curToken, Function: function}
	exp.Arguments = p.parseCallArguments()
	return exp
}

func (p *Parser) parseCallArguments() []ast.Expression {
	args := []ast.Expression{}
	if p.peekTokenIs(lexer.CLOSE_PAREN) {
		p.nextToken()
		return args
	}

	p.nextToken()
	arg := p.parseExpression(LOWEST)
	args = append(args, arg)

	for p.peekTokenIs(lexer.COMMA) {
		p.nextToken()
		p.nextToken()
		arg = p.parseExpression(LOWEST)
		args = append(args, arg)
	}

	if !p.expectPeek(lexer.CLOSE_PAREN) {
		return nil
	}
	return args
}

func (p *Parser) parseFunctionParams() []*ast.Identifier {
	idents := []*ast.Identifier{}

	if p.peekTokenIs(lexer.CLOSE_PAREN) {
		p.nextToken()
		return idents
	}

	p.nextToken()
	ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Value}
	idents = append(idents, ident)

	for p.peekTokenIs(lexer.COMMA) {
		p.nextToken()
		p.nextToken()
		ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Value}
		idents = append(idents, ident)
	}
	if !p.expectPeek(lexer.CLOSE_PAREN) {
		return nil
	}
	return idents
}

func (p *Parser) parseIfExpression() ast.Expression {
	expr := &ast.IfExpression{Token: p.curToken}
	if !p.expectPeek(lexer.OPEN_PAREN) {
		return nil
	}

	p.nextToken()

	expr.Condition = p.parseExpression(LOWEST)
	if !p.expectPeek(lexer.CLOSE_PAREN) {
		return nil
	}

	if !p.expectPeek(lexer.OPEN_CURLY) {
		return nil
	}

	expr.Consequence = p.parseBlockStatement()

	if p.peekTokenIs(lexer.ELSE) {
		p.nextToken()
		if !p.expectPeek(lexer.OPEN_CURLY) {
			return nil
		}
		expr.Alternative = p.parseBlockStatement()
	}

	return expr
}

func (p *Parser) parseBlockStatement() *ast.BlockStatment {
	block := &ast.BlockStatment{Token: p.curToken}
	block.Statements = []ast.Statement{}
	p.nextToken()
	for !p.curTokenIs(lexer.CLOSE_CURLY) && !p.curTokenIs(lexer.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}
	return block
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
	exp.Right = p.parseExpression(PREFIX)
	return exp
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

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()
	exp := p.parseExpression(LOWEST)
	if !p.expectPeek(lexer.CLOSE_PAREN) {
		return nil
	}
	return exp
}

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Token: p.curToken, Value: p.curTokenIs(lexer.TRUE)}
}
