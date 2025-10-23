package parser

import (
	"curgo/ast"
	"curgo/lexer"
	"testing"
)

func TestLetStatement(t *testing.T) {
	input := `
	let x = 5;
	let y = 10;
	let foobar =   563;
	`
	l := lexer.Tokenize(input)
	p := New(l)
	program := p.ParseProgram()

	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	checkParseErrors(t, p)
	if len(program.Statements) != 3 {
		t.Fatalf("Program Statements does not contain 3 statements. got=%d",
			len(program.Statements))
	}
	tests := []struct {
		exepectedIdent string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}
	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testLetStatement(t, stmt, tt.exepectedIdent) {
			return
		}
	}
}

func TestReturnStatements(t *testing.T) {
	inp := `
	return 5;
	return 10;
	return 341;
	`
	l := lexer.Tokenize(inp)
	p := New(l)
	program := p.ParseProgram()
	checkParseErrors(t, p)
	if len(program.Statements) != 3 {
		t.Fatalf("Program.Statements does not contain 3 statements. got=%d",
			len(program.Statements))
	}
	for _, stmt := range program.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("Stmt not casted to ast.ReturnStatement. got=%T", stmt)
		}
		if returnStmt.TokenLiteral() != "return" {
			t.Errorf("returnStmt.TokenLiteral() not 'return' got= %q", returnStmt.TokenLiteral())
		}
	}
}

func checkParseErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("Parser encoutered an error, Error Length: %d", len(p.Errors()))
	for _, msg := range errors {
		t.Errorf("Parser Error: %q", msg)
	}
	t.FailNow()
}

func testLetStatement(t *testing.T, stmt ast.Statement, name string) bool {
	if stmt.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral not 'let' got=%q", stmt.TokenLiteral())
		return false
	}

	letStmt, ok := stmt.(*ast.LetStatment)
	if !ok {
		t.Errorf("S not *ast.LetStatment. got=%T", stmt)
		return false
	}

	if letStmt.Name.Value != name {
		t.Errorf("letstmt.Name.Value not '%s'. got=%s", name, letStmt.Name.Value)
		return false
	}

	if letStmt.Name.Value != name {
		t.Errorf("s.Name not '%s'. got=%s", name, letStmt.Name.Value)
	}
	return true
}
