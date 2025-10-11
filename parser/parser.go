package parser

import (
	// "go/ast" // NOTE:
	"log"
	"parser/ast"
	"parser/lexer"
)

type parser struct {
	tokens []lexer.Token
	pos    int
}

func createParser(tokens []lexer.Token) *parser {
	createLookups()
	return &parser{
		tokens: tokens,
		pos:    0,
	}
}

func Parse(tokens []lexer.Token) ast.BlockStmt {
	body := make([]ast.Stmt, 0)
	p := createParser(tokens)

	for p.hasTokens() {
		body = append(body, parseStmt(p))
	}
	return ast.BlockStmt{
		Body: body,
	}
}

func (p *parser) currentToken() lexer.Token {
	return p.tokens[p.pos]
}

func (p *parser) eat() lexer.Token {
	token := p.currentToken()
	p.pos++
	return token
}

func (p *parser) hasTokens() bool {
	return p.pos <= len(p.tokens) && p.currentToken().Type != lexer.EOF
}

func (p *parser) expect(expectedType lexer.TokenKind) lexer.Token {
	return p.expectError(expectedType, nil)
}

func (p *parser) expectError(expectedType lexer.TokenKind, err any) lexer.Token {
	tk := p.currentToken().Type

	if tk != expectedType {
		if err != nil {
			log.Fatalf("Expected %s but got %s instead.",
				lexer.TokenKindString(expectedType), lexer.TokenKindString(tk))
		}
	}
	return p.eat()
}
