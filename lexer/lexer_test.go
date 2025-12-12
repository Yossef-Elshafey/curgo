package lexer

import (
	"curgo/types/tokens"
	"testing"
)

func TestLexer(t *testing.T) {
	compareValueToKind("foreach", tokens.FOREACH, t)
	compareValueToKind("import", tokens.IMPORT, t)
	compareValueToKind("true", tokens.TRUE, t)
	compareValueToKind("", tokens.EOF, t)
	compareValueToKind("foo", tokens.IDENTIFIER, t)
}

func compareValueToKind(input string, tk tokens.TokenKind, t *testing.T) {
	lex := Tokenize(string(input))
	token := lex[0]
	if token.Type != tk {
		t.Errorf("Type mismatch want= %s, got= %s", TokenKindString(tk), TokenKindString(token.Type))
	}
}

func TestExpectedTokens(t *testing.T) {
	tests := []struct {
		input    string
		expected []tokens.Token
	}{
		{input: "let x = 10;",
			expected: []tokens.Token{
				{Value: "let", Type: tokens.LET},
				{Value: "x", Type: tokens.IDENTIFIER},
				{Value: "=", Type: tokens.ASSIGNMENT},
				{Value: "10", Type: tokens.NUMBER},
				{Value: ";", Type: tokens.SEMI_COLON},
				{Value: "EOF", Type: tokens.EOF},
			}},
		{input: "fetch () {}",
			expected: []tokens.Token{
				{Value: "fetch", Type: tokens.FETCH},
				{Value: "(", Type: tokens.OPEN_PAREN},
				{Value: ")", Type: tokens.CLOSE_PAREN},
				{Value: "{", Type: tokens.OPEN_CURLY},
				{Value: "}", Type: tokens.CLOSE_CURLY},
				{Value: "EOF", Type: tokens.EOF},
			},
		},

		{input: "data if else return ==",
			expected: []tokens.Token{
				{Value: "data", Type: tokens.DATA},
				{Value: "if", Type: tokens.IF},
				{Value: "else", Type: tokens.ELSE},
				{Value: "return", Type: tokens.RETURN},
				{Value: "==", Type: tokens.EQUALS},
				{Value: "EOF", Type: tokens.EOF},
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
