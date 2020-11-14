package lexer

import "monkey/token"

type Lexer struct {
	input        string
	position     int  // curtent position in input (point to current character).
	readPosition int  // current reading position in input (after current ch).
	ch           byte //current ch under examination.
}

// New creates a new Lexer.
func New(input string) *Lexer {
	lexer := &Lexer{input: input} // if this doesn't behave as it should change the address to the struct.
	lexer.readChar()
	return lexer
}

// readChar gives us the next character and advance the
// position of the input string.
func (l *Lexer) readChar() {
	// if we reached the end of the input
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}

// NextToken returns the next token read.
func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhiteSpace()

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			// save l.ch in order not to lose the current character before proceeding to the next one.
			char := l.ch
			l.readChar()
			tok = token.Token{Type: token.EQ, Literal: string(char) + string(l.ch)}
		} else {
			tok = newToken(token.ASSIGN, l.ch)
		}
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case '-':
		tok = newToken(token.MINUS, l.ch)
	case '*':
		tok = newToken(token.ASTERISK, l.ch)
	case '<':
		tok = newToken(token.LT, l.ch)
	case '>':
		tok = newToken(token.GT, l.ch)
	case '!':
		if l.peekChar() == '=' {
			char := l.ch
			l.readChar()
			tok = token.Token{Type: token.NOT_EQ, Literal: string(char) + string(l.ch)}
		} else {
			tok = newToken(token.BANG, l.ch)
		}
	case '/':
		tok = newToken(token.SLASH, l.ch)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		// we check for identifiers because we didn't find a keyword.
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			// we return here in order not to readChar outside the switch
			// because we "exhaust" our positions in readIdentifier().
			return tok
		} else if isDigit(l.ch) {
			tok.Type = token.INT
			tok.Literal = l.readNumber()
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}
	l.readChar()
	return tok
}

// initialise new tokens
func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

// reads the name of an identifier until it reaches a non letter character.
func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

// skipWhiteSpace skips all the whitespaces/tabs/newlines etc because monkey
// doesn't take in mind all these.
func (l *Lexer) skipWhiteSpace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

// readNumber reads the number and returns it.
func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

// peekChar peeks into the next character of our input
// we use it as a helper to decipher != and == operators.
func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

// check if letter. This little func declares what our compiler will accept as identifiers
// e.g. it accepts foo_bar name. If we want to add more (!,?...) this is the place to do so.
func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

// isDigit checks if the ch byte is a digit.
func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
