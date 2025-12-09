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
				{Value: "", Type: EOF},
			}},
		{input: "fn () {}",
			expected: []Token{
				{Value: "fn", Type: FN},
				{Value: "(", Type: OPEN_PAREN},
				{Value: ")", Type: CLOSE_PAREN},
				{Value: "{", Type: OPEN_CURLY},
				{Value: "}", Type: CLOSE_CURLY},
				{Value: "", Type: EOF},
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
			if lex[i].Type != ts.expected[i].Type {
				t.Errorf("Lexer error: lex:%+v, expected: %+v", lex[i], ts.expected[i])
			}
		}
	}
}
