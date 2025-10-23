package ast

import (
	"curgo/lexer"
	"testing"
)

func TestStringify(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&LetStatment{
				Token: lexer.Token{Type: lexer.LET, Value: "let"},
				Name: &Identifier{
					Token: lexer.Token{Type: lexer.IDENTIFIER, Value: "myVar"},
					Value: "myVar",
				},
				Value: &Identifier{
					Token: lexer.Token{Type: lexer.IDENTIFIER, Value: "MyVar2"},
					Value: "MyVar2",
				},
			},
		},
	}

	if program.Stringify() != "let myVar = MyVar2;" {
		t.Errorf("Program.String() issued. got=%q", program.Stringify())
	}
}
