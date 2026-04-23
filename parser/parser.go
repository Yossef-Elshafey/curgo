package parser

import (
	"curgo/lexer"
	"curgo/types/ast"
	"curgo/types/tokens"
	"fmt"
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
	infixParseFn   func(ast.Expression)  ( ast.Expression, error )
	prefixParseFn  func()                ( ast.Expression, error )
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

func (p *Parser) ParseProgram() ( *ast.Program, error ) {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for !p.curTokenIs(token.EOF) {
		stmt, err := p.parseStmt()
		if err != nil {
			return nil, err
		}
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.advanceTokens()
	}
	return program, nil
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

func (p *Parser) syntaxError(msg string) error {
	return fmt.Errorf("Line %d\n %s\n", p.currentToken.Line, msg)
}

func (p *Parser) debugToken() {
	fmt.Printf("Current: %+v\n", p.currentToken)
	fmt.Printf("peek: %+v\n", p.peekToken)
}

func (p *Parser) expectPeekToBe(k token.TokenKind) bool {
	if !p.peekTokenIs(k) { return false }
	p.advanceTokens()
	return true
}

func (p *Parser) parseIdentifier() ( ast.Expression, error ) {
	return &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Value}, nil
}

func (p *Parser) parseStringLiteral() ( ast.Expression, error ) {
	return &ast.StringLiteral{Token: p.currentToken, Value: p.currentToken.Value}, nil
}

func (p *Parser) parseBinaryExpression(lhs ast.Expression) ( ast.Expression, error ) {
	expr := &ast.BinaryExpression{
		Token: p.currentToken,
		Operator: p.currentToken.Value,
		Left: lhs,
	}
	bp := p.currentTokenBindingPower()
	p.advanceTokens()
	parsedExpr, err := p.parseExpression(bp)
	if err != nil {
		return nil, err
	}
	expr.Right = parsedExpr
	return expr, nil
}

func (p *Parser) parseStmt() ( ast.Statement, error ) {
	switch p.currentToken.Kind {
	case token.FETCH:
		return p.parseFetchStatment()
	case token.LET:
		return p.parseLetStmt()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseLetStmt() ( *ast.LetStatement, error ) {
	if !p.expectPeekToBe(token.IDENTIFIER) {
		return nil, p.syntaxError("expect identifier after let")
	}
	ls := &ast.LetStatement{}
	ls.Identifier = &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Value}
	if !p.expectPeekToBe(token.ASSIGN) {
		return nil, p.syntaxError("expect '=' after let identifier") 
	}
	p.advanceTokens()
	expr, err := p.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}
	ls.Value = expr
	if !p.expectPeekToBe(token.SEMICOLON) {
		return nil, p.syntaxError("expect ';' after let stmt")
	}
	return ls, nil
}

func (p *Parser) parseFetchStatment() (*ast.FetchStmt, error) {
	fs := &ast.FetchStmt{Token: p.currentToken}
	fs.Body = []ast.Statement{}
	fs.Arguments = []*ast.Identifier{}

	if !p.expectPeekToBe(token.IDENTIFIER) {
		return nil, p.syntaxError("Expect Identifier after fetch")	
	}

	fs.FetchIdentifier = &ast.Identifier{
		Token: p.currentToken,
		Value:  p.currentToken.Value,
	}

	if !p.expectPeekToBe(token.LPAREN) {
		return nil, p.syntaxError("Expect '(' after fetch Identifier")	
	}
	args, err := p.parseFetchArguments()
	if err != nil {
		return nil, err
	}
	fs.Arguments = args
	if !p.expectPeekToBe(token.COLON) {
		return nil, p.syntaxError("expect ':' after fetch stmt arguments") 
	}

	for !p.peekTokenIs(token.ENDFETCH) {
		body, err := p.parseFetchBody()
		if err != nil {
			return nil, err
		}
		fs.Body = append(fs.Body, body)
	}

	p.advanceTokens()
	return fs, nil
}

func (p *Parser) parseFetchArguments() ( []*ast.Identifier, error ) {
	args := []*ast.Identifier{}
	if p.peekTokenIs(token.RPAREN) {
		p.advanceTokens()
		return args, nil
	}
	if !p.expectPeekToBe(token.IDENTIFIER) {
		return nil, p.syntaxError("Expect Identifier after (")
	}
	arg := &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Value}

	args = append(args, arg)
	for p.peekTokenIs(token.COMMA) {
		p.advanceTokens()
		p.advanceTokens()
		arg = &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Value}
		args = append(args, arg)
	}
	if !p.expectPeekToBe(token.RPAREN) {
		return nil, p.syntaxError("expect '(' after fetch stmt arguments")
	}
	return args, nil
}

func (p *Parser) parseFetchBody() ( ast.Statement, error ) {
	if !p.expectPeekToBe(token.IDENTIFIER) {
		fmt.Printf("%s\n", "use 'endfet' keyword to close fetch stmt")
		return nil, p.syntaxError("expect fetch stmt body to be include an argument of network call")
	}
	ca := &ast.CurgoAssignStatment{}
	ca.Token = p.currentToken
	ca.Arg = &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Value}

	if !p.expectPeekToBe(token.TRANSPILEASSIGN) {
		fmt.Printf("%s\n", "use 'endfet' keyword to close fetch stmt")
		return nil, p.syntaxError("expect '->' after fetch stmt body identifier")
	}

	p.advanceTokens()
	expr, err := p.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}
	ca.Value = expr

	if !p.expectPeekToBe(token.SEMICOLON) {
		return nil, p.syntaxError("expect ';' at the end of fetch stmt body argument")
	}
	return ca, nil
}

func (p *Parser) parseExpression(bp bindingPower) ( ast.Expression, error ) {
	prefix := p.prefixLookup[p.currentToken.Kind]
	if prefix == nil {
		return nil, p.syntaxError("prefix not supported" + p.currentToken.Value)
	}
	left, err := prefix()
	if err != nil {
		return nil, err
	}
	for !p.peekTokenIs(token.SEMICOLON) && bp < p.peekTokenBindingPower() {
		infix := p.infixLookup[p.peekToken.Kind]
		if infix == nil {
			return left, nil
		}
		p.advanceTokens()
		parsed, err := infix(left)
		if err != nil {
			return nil, err
		}
		left = parsed
	}
	return left, nil
}

func (p *Parser) parseNumberLiteral() ( ast.Expression, error ) {
	nl := &ast.NumberLiteral{}
	i, err := strconv.Atoi(p.currentToken.Value)
	if err != nil {
		return nil, p.syntaxError("Failed to convert %v to string" + p.currentToken.Value)
	}
	nl.Value = int64(i)
	nl.Token = p.currentToken
	return nl, nil
}

func (p *Parser) parseCallExpression(fs ast.Expression) ( ast.Expression, error ) {
	ce := &ast.CallExpression{Token: p.currentToken, Function: fs}
	ca, err := p.parseCallArgument()
	if err != nil {
		return nil, err
	}
	ce.Arguments = ca
	return ce, nil
}

func (p *Parser) parseCallArgument() ( []ast.Expression, error ) {
	args := []ast.Expression{}
	if p.peekTokenIs(token.RPAREN) {
		p.advanceTokens()
		return args, nil
	}
	p.advanceTokens()
	arg, err := p.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}
	args = append(args, arg)
	for p.peekTokenIs(token.COMMA) {
		p.advanceTokens()
		p.advanceTokens()
		arg, err = p.parseExpression(LOWEST)
		if err != nil {
			return nil, err
		}
		args = append(args, arg)
	}
	if !p.expectPeekToBe(token.RPAREN) {
		return nil, p.syntaxError("expect ')' after call expression")
	}
	// if !p.expectPeekToBe(token.SEMICOLON) {
	// 	return nil, p.syntaxError("expect ';' after call expression")
	// }
	return args, nil
}

func (p *Parser) parseExpressionStatement() ( *ast.ExpressionStatement, error ) {
	es := &ast.ExpressionStatement{}
	expr, err := p.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}
	es.Expression = expr
	if p.peekTokenIs(token.SEMICOLON) { p.advanceTokens() }
	return es, nil
}

func (p *Parser) parseGroupedExpression() ( ast.Expression, error ) {
	p.advanceTokens()

	expr, err := p.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}
	if !p.expectPeekToBe(token.RPAREN) {
		return nil, p.syntaxError("expect ')' after group expressions")
	}
	return expr, nil
}
