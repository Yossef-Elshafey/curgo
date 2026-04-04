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
	bindingPower   int
	infixParseFn   func(ast.Expression)  ast.Expression
	prefixParseFn  func()                ast.Expression
	PARSERSIG      int
)

const (
	TERM PARSERSIG = iota
	IGNORE PARSERSIG = iota
)

const (
	LOWEST bindingPower = iota
	SUM
	PRODUCT
)

var bindingPowerLookup = map[string]bindingPower{
	"+": SUM,
	"-": SUM,
	"*": PRODUCT,
	"/": PRODUCT,
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
		program.Statements = append(program.Statements, stmt)
		p.advanceTokens()
	}
	fmt.Printf("%+v\n", program.Statements[0].(*ast.LetStatement).Value)
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

func (p *Parser) currentTokenBindingPower() bindingPower {
	return bindingPowerLookup[p.currentToken.Value]
}

// NOTE: peekTokenIs, expectPeekToBe the way peek check handled is foolish
// implement peek token check such that errors, token advancing is clear
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

// 
func (p *Parser) expectPeekToBe(k tokens.TokenKind, sig PARSERSIG) bool {
	if p.peekToken.Kind != k {
		line := p.peekToken.Pos.Line
		lineIssue := utils.ReadSourceAsLines(line)
		p.isEndOfFetch()
		if sig == TERM {
			fmt.Printf("Parser:%d:%d: encouter error at line:\n %s\n", line-1,
				p.currentToken.Pos.End,
				lineIssue)

			log.Fatalf("Expect to find %s after '%s', got=%s", tokens.TokenKindStringify(k),
				p.currentToken.Value,
				tokens.TokenKindStringify(p.peekToken.Kind))

		}
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
		return nil
	}
}

func (p *Parser) parseLetStmt() *ast.LetStatement {
	if !p.expectPeekToBe(tokens.IDENTIFIER, TERM) {return nil}
	ls := &ast.LetStatement{}
	ls.Identifier = &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Value}
	if !p.expectPeekToBe(tokens.EQUAL, IGNORE) {
		if !p.expectPeekToBe(tokens.SEMI_COLON, TERM) {return nil}
		return ls
	}
	p.advanceTokens()
	ls.Value = p.parseExpression(LOWEST)
	if !p.expectPeekToBe(tokens.SEMI_COLON, TERM) {return nil}
	return ls
}

func (p *Parser) parseFetchStatment() *ast.FetchStmt {
	fs := &ast.FetchStmt{Token: p.currentToken}
	fs.Body = []ast.Statement{}
	if !p.expectPeekToBe(tokens.IDENTIFIER, TERM) {
		return nil
	}

	fs.FetchIdentifier = &ast.Identifier{
		Token: p.currentToken,
		Value:  p.currentToken.Value,
	}

	if !p.expectPeekToBe(tokens.COLON, TERM) {
		return nil
	}

	for !p.peekTokenIs(tokens.ENDFETCH) {
		fs.Body = append(fs.Body, p.parseFetchBody())
	}

	p.advanceTokens()
	return fs
}

func (p *Parser) parseFetchBody() ast.Statement {
	if !p.expectPeekToBe(tokens.IDENTIFIER, TERM) {
		return nil
	}
	ca := &ast.CurgoAssignStatment{}
	ca.Arg = &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Value}
	if !p.expectPeekToBe(tokens.TRANSPILE_ASSIGN, TERM) {
		return nil
	}

	if !p.expectPeekToBe(tokens.STRING, TERM) {
		return nil
	}

	ca.Value = p.parseExpression(LOWEST)

	if !p.expectPeekToBe(tokens.SEMI_COLON, TERM) {
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
