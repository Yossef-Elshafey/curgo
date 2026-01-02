package parser

import (
	"curgo/lexer"
	"curgo/types/ast"
	"curgo/utils"
	"testing"
)

func TestSyntaxErrorMessage(t *testing.T) {
	t.Skipf("this function is used to observe error logs")
	Parse(expectIdentifier())
	Parse(expectSemiColonMessage())
	Parse(expectStringErrorMessage())
}


func TestProgram(t *testing.T) {
	input := `
	fetch user: 
	header -> "foobar";
	foo -> "";
	bar -> "";
	endfet
	`
	utils.SetSource(input)
	tokens := lexer.Tokenize(input)
	program := Parse(tokens)

	if len(program.Statements) != 1 {
		t.Errorf("Program.Statements len is not 1 got=%d\n", len(program.Statements))
	}
	
	fs, ok  := program.Statements[0].(*ast.FetchStmt)
	if !ok {
		t.Errorf("Statements[0] is not ast.FetchStmt, got=%T\n", program.Statements[0])
	}

	if !checkFetchStatmentName(fs, "user") {
		t.Errorf("FetchStmt identifier is not user, got=%s\n", fs.FetchIdentifier.Value)
	}

	testFetchStatmentBody(t, fs, "header", "foo", "bar")
}

func testFetchStatmentBody(t *testing.T, fs *ast.FetchStmt, identifiers ...string) {
	if len(fs.Body) != len(identifiers) {
		t.Errorf("FetchStmt body is not equal to expected identifiers, body(%d) != idents(%d)", len(fs.Body), len(identifiers))
	}

	for idx, stmt := range fs.Body {
		curgoAssign, ok := stmt.(*ast.CurgoAssignStatment)
		if !ok {
			t.Errorf("Expect ast.curgoAssignment, got=%T\n",stmt)
		}
		if curgoAssign.Arg.Value != identifiers[idx] {
			t.Errorf("curgoAssignment failed expect %s, got=%s", identifiers[idx], curgoAssign.Arg.Value)
		}
	}
}

func checkFetchStatmentName(fs *ast.FetchStmt, name string)  bool {
	return fs.FetchIdentifier.Value == name
}

func expectStringErrorMessage() []lexer.Token {
	input := `
	fetch addUser: 
	header -> "Content-Type:application/json";
	foo -> ;
	endfet
	`
	utils.SetSource(input)
	tokens := lexer.Tokenize(input)
	return tokens
}

func expectSemiColonMessage() []lexer.Token {
	input := `
	fetch addUser: 
	foo -> ""
	`
	utils.SetSource(input)
	tokens := lexer.Tokenize(input)
	return tokens
}

func expectIdentifier() []lexer.Token {
	input := `
	fetch addUser: 
	foo -> "";
	`
	utils.SetSource(input)
	tokens := lexer.Tokenize(input)
	return tokens
}
