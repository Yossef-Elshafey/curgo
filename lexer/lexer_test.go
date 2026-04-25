package lexer

import (
	token "curgo/types/tokens"
	"fmt"
	"testing"
)

func TestLineCounting(t *testing.T) {
	source := `
let url = "http://localhost:8000";

fetch user(id):
	  host          ->  url;
	  header        ->  "Content-Type:application/json";
	  header        ->  "Accpet:json";
	  method        ->  "POST";
endfet

user("123");
	`
	l := New(source)
	tok := l.NextToken()
	if tok.Kind != token.LET && tok.Line != 1 {
		t.Errorf("Kind and Line is incorrect ")
	}
	tok = l.NextToken()
	tok = l.NextToken()
	tok = l.NextToken()
	tok = l.NextToken()

	if tok.Kind != token.SEMICOLON && tok.Line != 1 {
		t.Errorf("Kind and Line is incorrect")
	}
	tok = l.NextToken()

	if tok.Kind != token.FETCH && tok.Line != 3 {
		t.Errorf("Kind and Line is incorrect")
	}

	tok = l.NextToken()
	tok = l.NextToken()
	tok = l.NextToken()
	tok = l.NextToken()
	tok = l.NextToken()
	tok = l.NextToken()

	if tok.Kind != token.IDENTIFIER && tok.Line != 4 {
		t.Errorf("Kind and Line is incorrect, got= %+v\n", tok)
	}
}
