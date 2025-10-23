package parser

import (
	"curgo/ast"
	"curgo/lexer"
	"fmt"
)

type Parser struct {
	l         []lexer.Token
	curToken  lexer.Token
	peekToken lexer.Token
	pos       int
	errors    []string
}

func New(l []lexer.Token) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	// set the two token pointers appropriately
	p.nextToken()
	p.nextToken()
	return p
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
		return nil
	}
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
