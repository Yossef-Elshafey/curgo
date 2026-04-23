package lexer

import (
	"curgo/types/tokens"
)

type Lexer struct {
	input         string
	position      int
	readPosition  int
	ch            byte
	line          int
}


func New(source string) *Lexer {
	l := &Lexer{input: source, line: 1}
	l.readChar()
	return l
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token
	l.skipWhitespace()
	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Kind: token.EQ, Value: literal, Line: l.line}
		} else {
			tok = newToken(token.ASSIGN, l.ch, l.line)
		}
	case '+':
		tok = newToken(token.PLUS, l.ch, l.line)
	case '-':
		if l.peekChar() == '>' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Kind: token.TRANSPILEASSIGN, Value: literal, Line: l.line}
		} else {
			tok = newToken(token.MINUS, l.ch, l.line)
		}
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Kind: token.NOT_EQ, Value: literal, Line: l.line}
		} else {
			tok = newToken(token.BANG, l.ch, l.line)
		}
	case '/':
		if l.peekChar() == '/' {
			l.readChar()
			l.readChar()
			for l.ch != '\n' && l.ch != 0 {
				l.readChar()
			}
			return l.NextToken()
		} else {
			tok = newToken(token.SLASH, l.ch, l.line)
		}
	case '*':
		tok = newToken(token.ASTERISK, l.ch, l.line)
	case '<':
		tok = newToken(token.LT, l.ch, l.line)
	case '>':
		tok = newToken(token.GT, l.ch, l.line)
	case ';':
		tok = newToken(token.SEMICOLON, l.ch, l.line)
	case ',':
		tok = newToken(token.COMMA, l.ch, l.line)
	case '{':
		tok = newToken(token.LBRACE, l.ch, l.line)
	case '}':
		tok = newToken(token.RBRACE, l.ch, l.line)
	case '(':
		tok = newToken(token.LPAREN, l.ch, l.line)
	case ')':
		tok = newToken(token.RPAREN, l.ch, l.line)
	case '.':
		tok = newToken(token.DOT, l.ch, l.line)
	case 0:
		tok.Value = ""
		tok.Kind = token.EOF
		tok.Line = l.line
	case '"':
		tok.Kind = token.STRING
		tok.Value = l.readString('"')
		tok.Line = l.line
	case '`':
		tok.Kind = token.STRING
		tok.Value = l.readString('`')
		tok.Line = l.line
	case '[':
		tok = newToken(token.LBRACKET, l.ch, l.line)
	case ']':
		tok = newToken(token.RBRACKET, l.ch, l.line)
	case ':':
		tok = newToken(token.COLON, l.ch, l.line)
	default:
		if isLetter(l.ch) {
			tok.Value = l.readIdentifier()
			tok.Kind = token.LookupIdent(tok.Value)
			tok.Line = l.line
			return tok
		} else if isDigit(l.ch) {
			tok.Kind = token.NUMBER
			tok.Value = l.readNumber()
			tok.Line = l.line
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.ch, l.line)
		}
	}
	l.readChar()
	return tok
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	if l.ch == '\n' {
		l.line += 1
	}
	l.position = l.readPosition
	l.readPosition += 1
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func newToken(tokenKind token.TokenKind, ch byte, line int) token.Token {
	return token.Token{Kind: tokenKind, Value: string(ch), Line: line }
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

func (l *Lexer) readString(stop byte) string {
	position := l.position + 1
	for {
		l.readChar()
		if l.ch == stop || l.ch == 0 {
			break
		}
		if l.ch == stop || l.ch == 0 {
			break
		}
	}
	return l.input[position:l.position]
}
