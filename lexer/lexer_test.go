package lexer

import (
	"testing"
)

func TestLexer(t *testing.T) {
	compareValueToKind("foreach", FOREACH, t)
	compareValueToKind("import", IMPORT, t)
	compareValueToKind("true", TRUE, t)
	compareValueToKind("", EOF, t)
	compareValueToKind("foo", IDENTIFIER, t)
}

func compareValueToKind(input string, tk TokenKind, t *testing.T) {
	lex := Tokenize(string(input))
	token := lex[0]
	if token.Type != tk {
		t.Errorf("Type mismatch want= %s, got= %s", TokenKindString(tk), TokenKindString(token.Type))
	}
}

func TestExpectedTokens(t *testing.T) {
	tests := []struct {
		input    string
		expected []Token
	}{
		{input: "let x = 10;",
			expected: []Token{
				{Value: "let", Type: LET},
				{Value: "x", Type: IDENTIFIER},
				{Value: "=", Type: ASSIGNMENT},
				{Value: "10", Type: NUMBER},
				{Value: ";", Type: SEMI_COLON},
				{Value: "EOF", Type: EOF},
			}},
		{input: "fetch () {}",
			expected: []Token{
				{Value: "fetch", Type: FETCH},
				{Value: "(", Type: OPEN_PAREN},
				{Value: ")", Type: CLOSE_PAREN},
				{Value: "{", Type: OPEN_CURLY},
				{Value: "}", Type: CLOSE_CURLY},
				{Value: "EOF", Type: EOF},
			},
		},

		{input: "data if else return ==",
			expected: []Token{
				{Value: "data", Type: DATA},
				{Value: "if", Type: IF},
				{Value: "else", Type: ELSE},
				{Value: "return", Type: RETURN},
				{Value: "==", Type: EQUALS},
				{Value: "EOF", Type: EOF},
			},
		},
	}

	for _, ts := range tests {
		lex := Tokenize(ts.input)
		if len(lex) != len(ts.expected) {
			t.Errorf("Lexer error: lex got %d tokens, expected has %d tokens", len(lex), len(ts.expected))
			t.FailNow()
		}

		for i := 0; i < len(lex); i++ {
			if lex[i].Type != ts.expected[i].Type || lex[i].Value != ts.expected[i].Value {
				t.Errorf("Lexer error: lex:%+v, expected: %+v", lex[i], ts.expected[i])
			}
		}
	}
}
