package parser

import (
	"curgo/lexer"
	"curgo/types/ast"
	"testing"
)

func TestMapLiteral(t *testing.T) {
	source := `let cord = {x:2, y:4};`
	tokens := lexer.New(source)
	p := New(tokens)
	program, err := p.ParseProgram()
	if err != nil {
		t.Error(err.Error())
	}
	testNumberOfStatments(t, len(program.Statements), 1)
	letStmt, ok := program.Statements[0].(*ast.LetStatement)
	if !ok {
		t.Errorf("expect program.Statements[0] to be letstmt, got= %T", program.Statements[0])
	}
	mapLiteral, ok := letStmt.Value.(*ast.MapLiteral)
	if !ok {
		t.Errorf("expect letStmt value to be MapLiteral, got= %T", letStmt.Value)
	}
	val, ok := mapLiteral.Elements["x"]
	if !ok {
		t.Errorf("expect map value of 'x'")
	}
	vInt, _ := val.(*ast.NumberLiteral)
	if vInt.Value != 2 {
		t.Errorf("expect map values of 'x' to be 2, got=%d", vInt.Value)
	}

	val, ok = mapLiteral.Elements["y"]
	if !ok {
		t.Errorf("expect map value of 'y'")
	}
	vInt, _ = val.(*ast.NumberLiteral)
	if vInt.Value != 4 {
		t.Errorf("expect map values of 'y' to be 4, got=%d", vInt.Value)
	}
}

func TestIndexing(t *testing.T) {
	source := `x[0]`
	tokens := lexer.New(source)
	p := New(tokens)
	program, err := p.ParseProgram()
	if err != nil {
		t.Error(err.Error())
	}
	testNumberOfStatments(t, len(program.Statements), 1)
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("expect program.Statements[0] to be ast.ExpressionStatement, got=%T", program.Statements[0])
	}
	expr, ok := stmt.Expression.(*ast.Indexing)
	if expr.Ident.Stringify() != "x" {
		t.Errorf("expect identifier to be 'x', got=%s", expr.Ident.Stringify())
	}
	num, ok := expr.Target.(*ast.NumberLiteral)
	if !ok {
		t.Errorf("expect expr.Target to be NumberLiteral, got= %T", expr.Target)
	}
	if num.Value != 0 {
		t.Errorf("expect to access index of '0', got=%v", expr.Target)
	}
}

func TestArrayLieral(t *testing.T) {
	source := `let x = [1,2,3, a];`
	tokens := lexer.New(source)
	p := New(tokens)
	program, err := p.ParseProgram()
	if err != nil {
		t.Error(err.Error())
	}
	stmt, ok := program.Statements[0].(*ast.LetStatement)
	if !ok {
		t.Errorf("expect program.Statements[0] to be ast.LetStatement, got=%T", program.Statements[0])
	}
	expr, ok := stmt.Value.(*ast.ArrayLiteral)
	if !ok {
		t.Errorf("expect let stmt.Value to ast.ArrayLiteral, got=%T", stmt.Value)
	}
	testNumberOfStatments(t, len(expr.Elements), 4)
	num, ok := expr.Elements[0].(*ast.NumberLiteral)
	if !ok {
		t.Errorf("expect expr.Elements[0] to be NumberLiteral, got=%T", expr.Elements[0])
	}
	if num.Value != 1 {
		t.Errorf("expect first value in array to be %d, got=%d", 1, num.Value)
	}
}

func TestIfStmt(t *testing.T) {
	source := `if 1 == 1 {
		let x = 1;
		let y = x + 1;
	} else {
		let y = 2;
	}`
	tokens := lexer.New(source)
	p := New(tokens)
	program, err := p.ParseProgram()
	if err != nil {
		t.Error(err.Error())
	}
	testNumberOfStatments(t, len(program.Statements), 1)
	ifstmt, ok := program.Statements[0].(*ast.IfStmt)
	if !ok {
		t.Errorf("expect to get ast.ifstmt, got=%T", program.Statements[0])
	}
	testInfix(t, ifstmt.Cond, 1, "==", 1)
	testNumberOfStatments(t, len(ifstmt.Consequences.Statements), 2)
	testNumberOfStatments(t, len(ifstmt.Alternative.Statements), 1)
}

func TestMemberAccess(t *testing.T) {
	source := `a.b.c`
	tokens := lexer.New(source)
	p := New(tokens)
	program, err := p.ParseProgram()
	if err != nil {
		t.Error(err.Error())
	}
	expr, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("expect program.Statements[0] to exprStatement, got= %T", program.Statements[0])
	}
	rightExpr := expr.Expression.(*ast.SuffixExpression).Member
	testIdentifier(t, rightExpr,"c")
	leftExpr := expr.Expression.(*ast.SuffixExpression).Left
	testMemberAccess(t, leftExpr, "a", ".", "b")
}

func TestIncorrectFetchStatement(t *testing.T) {
	source := []struct{
		input string
	}{
		{input: "fetch"},
		{input: `fetch )`},
		{input: `fetch user() `},
		{input: `fetch user(): `},
		{input: `fetch user():
			header -> "Content-Type:application/json"endfet`},
	}

	for _, s := range source {
		tokens := lexer.New(s.input)
		p := New(tokens)
		program, err := p.ParseProgram()

		if err == nil {
			t.Errorf("Program should fail with syntax error but it didnt, got= %+v\n", program.Statements)
		}
	}
}

func TestProgram(t *testing.T) {
	source := `
	fetch user(id, payload): 
	  host    ->  "localhost:8888/" + id;
	  header  ->  "Content-Type:application/json";
	  method  ->  "POST"; 
	  data    ->  payload;
	endfet
	let foo = "bar";
	user(foo, "x");
	`
	tokens := lexer.New(source)
	p := New(tokens)
	program, _ := p.ParseProgram()
	testNumberOfStatments(t, len(program.Statements), 3)
	fs, ok := program.Statements[0].(*ast.FetchStmt)
	if !ok {t.Errorf("program.Statements[0] is not FetchStmt, got= %T", program.Statements[0])}
	b1 := fs.Body[0].(*ast.CurgoAssignStatment)
	testInfix(t, b1.Value, "localhost:8888/", "+", "id")
	b4 := fs.Body[3].(*ast.CurgoAssignStatment)
	testIdentifier(t, b4.Value, "payload")
}

func TestFetchStatement(t *testing.T) {
	source := `
	fetch user(id, payload):
	  host          ->  url;
	  header        ->  "Content-Type:application/json";
	  header        ->  "Accpet:json";
	  method        ->  "POST";
	  data          ->  ` + "`" + `{"fname":"yossef", "lname":"elshafey"}` + "`;" + "endfet"

	tokens := lexer.New(source)
	p := New(tokens)
	program, _ := p.ParseProgram()
	testNumberOfStatments(t, len(program.Statements), 1)
	fs, ok := program.Statements[0].(*ast.FetchStmt)
	if !ok {
		t.Errorf("Program.Statements[0] is not FetchStmt, got= %T", program.Statements[0])
	}
	if len( fs.Arguments ) != 2 {
		t.Errorf("Expect FetchStmt Arguments to be 2, got= %d", len(fs.Arguments))
	}
}

func TestLetStatment(t *testing.T) {
	source := `let i = 2 + 3;`
	tokens := lexer.New(source)
	p := New(tokens)
	program, _ := p.ParseProgram()
	testNumberOfStatments(t, len(program.Statements), 1)
	ls, ok := program.Statements[0].(*ast.LetStatement)
	if !ok {
		t.Errorf("program.Statements[0] is not LetStatement, got= %T", program.Statements[0])
	}
	testIdentifier(t, ls.Identifier, "i")
	testInfix(t, ls.Value, 2, "+", 3)
}

func TestCallExpression(t *testing.T) {
	source := `add(1,2 + 3, "foo" + "bar", "foo");`
	l := lexer.New(source)
	p := New(l)
	program, _ := p.ParseProgram()
	testNumberOfStatments(t, len(program.Statements), 1)
	es, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("Expect program.Statements[0] to be ExpressionStatement, got=%T", program.Statements[0])
	}
	ce := es.Expression.(*ast.CallExpression)
	if len(ce.Arguments) != 4 {
		t.Errorf("CallExpression.Arguments number does not match expected got=%d, want=%d", len(ce.Arguments), 4)
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
		t.Errorf("exp is not BinaryExpression, got=%T", exp)
	}

	if !testLiteralExpression(t, bExp.Left, lhs) {
		return false
	}

	if bExp.Operator != operator {
		t.Errorf("BinaryExpression.operator doesnt match expect=%s, got=%s", operator, bExp.Operator)
		return false
	}

	if !testLiteralExpression(t, bExp.Right, rhs) {
		return false
	}

	return true
}

func testMemberAccess(
	t *testing.T,
	exp ast.Expression,
	lhs interface{},
	operator string,
	rhs string,
) bool {
	maExpr, ok := exp.(*ast.SuffixExpression)
	if !ok {
		t.Errorf("exp is not MemberAccess, got=%T", exp)
	}

	if !testLiteralExpression(t, maExpr.Left, lhs) {
		return false
	}

	if maExpr.Operator != operator {
		t.Errorf("MemberAccess.operator doesnt match expect=%s, got=%s", operator, maExpr.Operator)
		return false
	}

	if !testIdentifier(t, maExpr.Member, rhs) {
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
		t.Errorf("expect exp to be NumberLiteral, got= %T", exp)
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
		if !testIdentifier(t, exp, value) {
			t.Errorf("expect exp to be StringLiteral, got= %T", exp)
			return false
		}
		return true
	}
	if sl.Value != value {
		t.Errorf("Expect NumberLiteral to be %s, got= %s", value, sl.Value)
	}
	return true
}

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("expect exp to be Identifier, got= %T", exp)
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
