package lexer

import (
	token "curgo/types/tokens"
	"testing"
)

func checkToken(t *testing.T, tok token.Token, kind token.TokenKind, value string, line int) {
	t.Helper()
	if tok.Kind != kind {
		t.Errorf("wrong Kind: expected=%q, got=%q (val=%q, line=%d)", kind, tok.Kind, tok.Value, tok.Line)
	}
	if tok.Line != line {
		t.Errorf("wrong Line for %q: expected=%d, got=%d (val=%q)", kind, line, tok.Line, tok.Value)
	}
	if tok.Value != value {
		t.Errorf("wrong Value for %q: expected=%q, got=%q (line=%d)", kind, value, tok.Value, tok.Line)
	}
}

func checkEOF(t *testing.T, tok token.Token, line int) {
	t.Helper()
	if tok.Kind != token.EOF {
		t.Errorf("wrong Kind: expected=EOF, got=%q (val=%q, line=%d)", tok.Kind, tok.Value, tok.Line)
	}
	if tok.Line != line {
		t.Errorf("wrong Line for EOF: expected=%d, got=%d", line, tok.Line)
	}
}

func TestLineCounting(t *testing.T) {
	source := `let url = "http://localhost:8000";

fetch user(id):
	  host          ->  url;
	  header        ->  "Content-Type:application/json";
	  header        ->  "Accpet:json";
	  method        ->  "POST";
endfet

user("123");
`
	l := New(source)

	checkToken(t, l.NextToken(), token.LET, "let", 1)
	checkToken(t, l.NextToken(), token.IDENTIFIER, "url", 1)
	checkToken(t, l.NextToken(), token.ASSIGN, "=", 1)
	checkToken(t, l.NextToken(), token.STRING, "http://localhost:8000", 1)
	checkToken(t, l.NextToken(), token.SEMICOLON, ";", 1)

	checkToken(t, l.NextToken(), token.FETCH, "fetch", 3)
	checkToken(t, l.NextToken(), token.IDENTIFIER, "user", 3)
	checkToken(t, l.NextToken(), token.LPAREN, "(", 3)
	checkToken(t, l.NextToken(), token.IDENTIFIER, "id", 3)
	checkToken(t, l.NextToken(), token.RPAREN, ")", 3)
	checkToken(t, l.NextToken(), token.COLON, ":", 3)

	checkToken(t, l.NextToken(), token.IDENTIFIER, "host", 4)
	checkToken(t, l.NextToken(), token.TRANSPILEASSIGN, "->", 4)
	checkToken(t, l.NextToken(), token.IDENTIFIER, "url", 4)
	checkToken(t, l.NextToken(), token.SEMICOLON, ";", 4)

	checkToken(t, l.NextToken(), token.IDENTIFIER, "header", 5)
	checkToken(t, l.NextToken(), token.TRANSPILEASSIGN, "->", 5)
	checkToken(t, l.NextToken(), token.STRING, "Content-Type:application/json", 5)
	checkToken(t, l.NextToken(), token.SEMICOLON, ";", 5)

	checkToken(t, l.NextToken(), token.IDENTIFIER, "header", 6)
	checkToken(t, l.NextToken(), token.TRANSPILEASSIGN, "->", 6)
	checkToken(t, l.NextToken(), token.STRING, "Accpet:json", 6)
	checkToken(t, l.NextToken(), token.SEMICOLON, ";", 6)

	checkToken(t, l.NextToken(), token.IDENTIFIER, "method", 7)
	checkToken(t, l.NextToken(), token.TRANSPILEASSIGN, "->", 7)
	checkToken(t, l.NextToken(), token.STRING, "POST", 7)
	checkToken(t, l.NextToken(), token.SEMICOLON, ";", 7)

	checkToken(t, l.NextToken(), token.ENDFETCH, "endfet", 8)

	checkToken(t, l.NextToken(), token.IDENTIFIER, "user", 10)
	checkToken(t, l.NextToken(), token.LPAREN, "(", 10)
	checkToken(t, l.NextToken(), token.STRING, "123", 10)
	checkToken(t, l.NextToken(), token.RPAREN, ")", 10)
	checkToken(t, l.NextToken(), token.SEMICOLON, ";", 10)

	checkEOF(t, l.NextToken(), 11)
}

func TestLineLeadingNewline(t *testing.T) {
	source := "\nlet x = 5;"
	l := New(source)

	checkToken(t, l.NextToken(), token.LET, "let", 2)
	checkToken(t, l.NextToken(), token.IDENTIFIER, "x", 2)
	checkToken(t, l.NextToken(), token.ASSIGN, "=", 2)
	checkToken(t, l.NextToken(), token.NUMBER, "5", 2)
	checkToken(t, l.NextToken(), token.SEMICOLON, ";", 2)
	checkEOF(t, l.NextToken(), 2)
}

func TestLineEmptyInput(t *testing.T) {
	checkEOF(t, New("").NextToken(), 1)
}

func TestLineOnlyWhitespace(t *testing.T) {
	checkEOF(t, New("   \t  \n  \n  ").NextToken(), 3)
	checkEOF(t, New("\n\n\n").NextToken(), 4)
}

func TestLineNoTrailingNewline(t *testing.T) {
	source := "abc\n123"
	l := New(source)

	checkToken(t, l.NextToken(), token.IDENTIFIER, "abc", 1)
	checkToken(t, l.NextToken(), token.NUMBER, "123", 2)
	checkEOF(t, l.NextToken(), 2)
}

func TestLineMultipleBlankLines(t *testing.T) {
	source := "a\n\n\n\nb"
	l := New(source)

	checkToken(t, l.NextToken(), token.IDENTIFIER, "a", 1)
	checkToken(t, l.NextToken(), token.IDENTIFIER, "b", 5)
}

func TestLineComments(t *testing.T) {
	source := `let x = 1; // comment here
y = 2; // another
z = 3;
`
	l := New(source)

	checkToken(t, l.NextToken(), token.LET, "let", 1)
	checkToken(t, l.NextToken(), token.IDENTIFIER, "x", 1)
	checkToken(t, l.NextToken(), token.ASSIGN, "=", 1)
	checkToken(t, l.NextToken(), token.NUMBER, "1", 1)
	checkToken(t, l.NextToken(), token.SEMICOLON, ";", 1)

	checkToken(t, l.NextToken(), token.IDENTIFIER, "y", 2)
	checkToken(t, l.NextToken(), token.ASSIGN, "=", 2)
	checkToken(t, l.NextToken(), token.NUMBER, "2", 2)
	checkToken(t, l.NextToken(), token.SEMICOLON, ";", 2)

	checkToken(t, l.NextToken(), token.IDENTIFIER, "z", 3)
	checkToken(t, l.NextToken(), token.ASSIGN, "=", 3)
	checkToken(t, l.NextToken(), token.NUMBER, "3", 3)
	checkToken(t, l.NextToken(), token.SEMICOLON, ";", 3)

	checkEOF(t, l.NextToken(), 4)
}

func TestLineCommentAtEOF(t *testing.T) {
	source := "abc // trailing comment"
	l := New(source)

	checkToken(t, l.NextToken(), token.IDENTIFIER, "abc", 1)
	checkEOF(t, l.NextToken(), 1)
}

func TestLineStringTokens(t *testing.T) {
	source := `"hello"
"world"
"foo"`
	l := New(source)

	checkToken(t, l.NextToken(), token.STRING, "hello", 1)
	checkToken(t, l.NextToken(), token.STRING, "world", 2)
	checkToken(t, l.NextToken(), token.STRING, "foo", 3)
}

func TestLineCRLF(t *testing.T) {
	source := "a\r\nb\r\nc"
	l := New(source)

	checkToken(t, l.NextToken(), token.IDENTIFIER, "a", 1)
	checkToken(t, l.NextToken(), token.IDENTIFIER, "b", 2)
	checkToken(t, l.NextToken(), token.IDENTIFIER, "c", 3)
	checkEOF(t, l.NextToken(), 3)
}

func TestLineMixedOperators(t *testing.T) {
	source := "1 == 2\n3 != 4\n5 -> 6"
	l := New(source)

	checkToken(t, l.NextToken(), token.NUMBER, "1", 1)
	checkToken(t, l.NextToken(), token.EQ, "==", 1)
	checkToken(t, l.NextToken(), token.NUMBER, "2", 1)

	checkToken(t, l.NextToken(), token.NUMBER, "3", 2)
	checkToken(t, l.NextToken(), token.NOT_EQ, "!=", 2)
	checkToken(t, l.NextToken(), token.NUMBER, "4", 2)

	checkToken(t, l.NextToken(), token.NUMBER, "5", 3)
	checkToken(t, l.NextToken(), token.TRANSPILEASSIGN, "->", 3)
	checkToken(t, l.NextToken(), token.NUMBER, "6", 3)
}

func TestLineIllegalToken(t *testing.T) {
	source := "a\n@\nb"
	l := New(source)

	checkToken(t, l.NextToken(), token.IDENTIFIER, "a", 1)
	checkToken(t, l.NextToken(), token.ILLEGAL, "@", 2)
	checkToken(t, l.NextToken(), token.IDENTIFIER, "b", 3)
}

func TestLineCRLFMixed(t *testing.T) {
	tests := []struct {
		name     string
		source   string
		expected []token.Token
	}{
		{
			name:   "crlf then lf",
			source: "a\r\nb\nc",
			expected: []token.Token{
				{Kind: token.IDENTIFIER, Value: "a", Line: 1},
				{Kind: token.IDENTIFIER, Value: "b", Line: 2},
				{Kind: token.IDENTIFIER, Value: "c", Line: 3},
			},
		},
		{
			name:   "cr only",
			source: "a\rb\rc",
			expected: []token.Token{
				{Kind: token.IDENTIFIER, Value: "a", Line: 1},
				{Kind: token.IDENTIFIER, Value: "b", Line: 2},
				{Kind: token.IDENTIFIER, Value: "c", Line: 3},
			},
		},
		{
			name:   "lf then crlf",
			source: "a\nb\r\nc",
			expected: []token.Token{
				{Kind: token.IDENTIFIER, Value: "a", Line: 1},
				{Kind: token.IDENTIFIER, Value: "b", Line: 2},
				{Kind: token.IDENTIFIER, Value: "c", Line: 3},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := New(tt.source)
			for i, exp := range tt.expected {
				tok := l.NextToken()
				if tok.Kind != exp.Kind || tok.Value != exp.Value || tok.Line != exp.Line {
					t.Errorf("token[%d]: expected {%s %q line=%d}, got {%s %q line=%d}",
						i, exp.Kind, exp.Value, exp.Line, tok.Kind, tok.Value, tok.Line)
				}
			}
		})
	}
}

func TestLineNumbersAdjacent(t *testing.T) {
	source := `1
2

3`
	l := New(source)

	checkToken(t, l.NextToken(), token.NUMBER, "1", 1)
	checkToken(t, l.NextToken(), token.NUMBER, "2", 2)
	checkToken(t, l.NextToken(), token.NUMBER, "3", 4)
}

func TestLineBacktickString(t *testing.T) {
	source := "a\n`raw\nstring`\nb"
	l := New(source)

	checkToken(t, l.NextToken(), token.IDENTIFIER, "a", 1)
	checkToken(t, l.NextToken(), token.STRING, "raw\nstring", 2)
	checkToken(t, l.NextToken(), token.IDENTIFIER, "b", 4)
}
