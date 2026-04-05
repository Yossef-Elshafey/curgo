package parser

import (
	"curgo/lexer"
	"curgo/types/ast"
	"testing"
)


func TestProgram(t *testing.T) {
	source := `
	fetch user: 
	endfet `
	tokens := lexer.Tokenize(source)
	p := Parse(tokens)
	testNumberOfStatments(t, len(p.Statements), 1)
	fs, ok := p.Statements[0].(*ast.FetchStmt)
	if !ok {
		t.Errorf("Program.Statements[0] is not ast.FetchStmt, got= %T", p.Statements[0])
	}
	if fs.FetchIdentifier.Value != "user" {
		t.Errorf("FetchStmt.FetchIdentifier != user, got=%s", fs.FetchIdentifier.Value)
	}
}

func TestLetStatment(t *testing.T) {
	source := `let i = 2 + 3;`
	tokens := lexer.Tokenize(source)
	p := Parse(tokens)
	testNumberOfStatments(t, len(p.Statements), 1)
	ls, ok := p.Statements[0].(*ast.LetStatement)
	if !ok {
		t.Errorf("program.Statements[0] is not LetStatement, got= %t", p.Statements[0])
	}
	testIdentifier(t, ls.Identifier, "i")
	testInfix(t, ls.Value, 2, "+", 3)
}

func TestCallExpression(t *testing.T) {
	source := `add(1,2 + 3, "foo" + "bar", "foo");`
	l := lexer.Tokenize(source)
	p := Parse(l)
	testNumberOfStatments(t, len(p.Statements), 1)
	es, ok := p.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("Expect program.Statements[0] to be ExpressionStatement, got=%t", p.Statements[0])
	}
	ce := es.Expression.(*ast.CallExpression)
	if len(ce.Arguments) != 4 {
		t.Errorf("CallExpression.Arguments number does not match expected got=%d, want=%d", len(ce.Arguments), 3)
	}
	testLiteralExpression(t, ce.Arguments[0], 1) 
	testInfix(t, ce.Arguments[1], 2, "+", 3) 
	testInfix(t, ce.Arguments[2], "foo", "+", "bar")
	testLiteralExpression(t, ce.Arguments[3], "foo")
}

func testInfix(
	t *testing.T,
	exp ast.Expression,
	lhs interface{},
	operator string,
	rhs interface{},
) bool {
	bExp, ok := exp.(*ast.BinaryExpression)
	if !ok {
		t.Errorf("exp is not BinaryExpression, got=%t", exp)
	}

	if !testLiteralExpression(t, bExp.Left, lhs) {
		return false
	}

	if bExp.Operator.Value != operator {
		t.Errorf("BinaryExpression.operator doesnt match expect=%s, got=%s", operator, bExp.Operator.Value)
		return false
	}

	if !testLiteralExpression(t, bExp.Right, rhs) {
		return false
	}

	return true
}

func testLiteralExpression(
	t *testing.T,
	exp ast.Expression,
	expected interface{},
) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case string:
		return testStringLiteral(t, exp, v)
	}
	t.Errorf("type of exp not handled. got=%T", exp)
	return false
}

func testIntegerLiteral(t *testing.T, exp ast.Expression, value int64) bool {
	il, ok := exp.(*ast.NumberLiteral)
	if !ok {
		t.Errorf("expect exp to be NumberLiteral, got= %t", exp)
		return false
	}
	if il.Value != value {
		t.Errorf("Expect NumberLiteral to be %d, got= %d", value, il.Value)
	}
	return true
}

func testStringLiteral(t *testing.T, exp ast.Expression, value string) bool {
	sl, ok := exp.(*ast.StringLiteral)
	if !ok {
		t.Errorf("expect exp to be StringLiteral, got= %t", exp)
		return false
	}
	if sl.Value != value {
		t.Errorf("Expect NumberLiteral to be %s, got= %s", value, sl.Value)
	}
	return true
}

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("expect exp to be Identifier, got= %t", exp)
		return false
	}
	if ident.Value != value {
		t.Errorf("Expect NumberLiteral to be %s, got= %s", value, ident.Value)
	}
	return true
}

func testNumberOfStatments(t *testing.T, got, expected int) {
	if got != expected {
		t.Errorf("expect program.Statements to be %d, got=%d", expected, got)
	}
}
