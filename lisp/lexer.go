package lisp

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

const (
	eof = -1

	tokenLeftParen tokenType = iota
	tokenRightParen
	tokenEOF
	tokenError
	tokenBool
	tokenNil
	tokenSymbol
	tokenString
	tokenInt
	tokenFloat
	tokenKeyword
)

type tokenType int

type pos int

type token struct {
	typ   tokenType
	pos   pos
	value string
	line  pos
}

type stateFn func(*lexer) stateFn

type lexer struct {
	state  stateFn
	tokens chan token
	input  string
	pos,
	start,
	width pos
	parenDepth int
	line       pos
}

func newLexer(input string) *lexer {
	l := &lexer{
		input:  input,
		tokens: make(chan token),
		line:   1,
	}
	return l
}

func (l *lexer) run() {
	for l.state = lexText; l.state != nil; {
		l.state = l.state(l)
	}
	close(l.tokens)
}

func (l *lexer) next() rune {
	if int(l.pos) >= len(l.input) {
		l.width = 0
		return eof
	}
	r, w := utf8.DecodeRuneInString(l.input[l.pos:])
	l.width = pos(w)
	l.pos += l.width
	if r == '\n' {
		l.line++
	}
	return r
}

func (l *lexer) errorf(format string, args ...any) stateFn {
	l.tokens <- token{tokenError, l.start, fmt.Sprintf(format, args...), l.line}
	return nil
}

func (l *lexer) accept(valid string) bool {
	if strings.ContainsRune(valid, l.next()) {
		return true
	}
	l.backup()
	return false
}

func (l *lexer) acceptRun(valid string) {
	for strings.ContainsRune(valid, l.next()) {
	}
	l.backup()
}

func (l *lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

func (l *lexer) backup() {
	l.pos -= l.width
	if l.width == 1 && l.input[l.pos] == '\n' {
		l.line--
	}
}

func (l *lexer) emit(t tokenType) {
	l.tokens <- token{t, l.start, l.input[l.start:l.pos], l.line}
	l.start = l.pos
}

func (l *lexer) ignore() {
	l.start = l.pos
}

func (l *lexer) drain() {
	for range l.tokens {
	}
}

func lexText(l *lexer) stateFn {
	for r := l.next(); isSpace(r) || r == '\n'; l.next() {
		r = l.peek()
	}
	l.backup()
	l.ignore()

	switch r := l.next(); {
	case r == eof:
		l.emit(tokenEOF)
		if l.parenDepth != 0 {
			return l.errorf("unclosed left paren")
		}
		return nil
	case r == ';':
		return lexComment
	case r == ',':
		l.ignore()
		return lexText
	case r == '(':
		return lexLeftParen
	case r == ')':
		return lexRightParen
	case r == '-' || ('0' <= r && r <= '9'):
		return lexNumber
	case r == ':':
		return lexKeyword
	case isAlphaNumeric(r):
		return lexSymbol
	case r == '"':
		return lexString
	default:
		l.errorf("unhandled token %q", r)
		return nil
	}
}

func lexString(l *lexer) stateFn {
	for r := l.next(); r != '"'; r = l.next() {
		if r == '\\' {
			r = l.next()
		}
		if r == eof {
			return l.errorf("unclosed quoted string")
		}
	}
	l.emit(tokenString)
	return lexText
}

func lexKeyword(l *lexer) stateFn {
	for r := l.next(); isAlphaNumeric(r); r = l.next() {
	}
	l.backup()
	l.emit(tokenKeyword)
	return lexText
}

func lexComment(l *lexer) stateFn {
	i := strings.Index(l.input[l.pos:], "\n")
	if i < 0 {
		return nil
	}
	l.pos += pos(i)
	l.ignore()
	return lexText
}

func lexSymbol(l *lexer) stateFn {
	for r := l.next(); isAlphaNumeric(r); r = l.next() {
	}
	l.backup()
	switch l.input[l.start:l.pos] {
	case "true", "false":
		l.emit(tokenBool)
		return lexText
	case "nil":
		l.emit(tokenNil)
		return lexText
	default:
		l.emit(tokenSymbol)
		return lexText
	}
}

func lexNumber(l *lexer) stateFn {
	const digits = "0123456789"

	l.accept("-")
	l.acceptRun(digits)
	if l.accept(".") {
		l.acceptRun(digits)
	}
	if l.pos-l.start == 1 && !unicode.IsDigit(rune(l.input[l.start])) {
		return lexSymbol
	}
	if strings.ContainsRune(l.input[l.start:l.pos], '.') {
		l.emit(tokenFloat)
	} else {
		l.emit(tokenInt)
	}
	return lexText
}

func lexLeftParen(l *lexer) stateFn {
	l.emit(tokenLeftParen)
	l.parenDepth++
	return lexText
}

func lexRightParen(l *lexer) stateFn {
	l.parenDepth--
	if l.parenDepth < 0 {
		return l.errorf("unexpected right paren")
	}
	l.emit(tokenRightParen)
	return lexText
}

func isSpace(r rune) bool {
	return r == ' ' || r == '\t'
}

func isAlphaNumeric(r rune) bool {
	return r == '>' || r == '<' || r == '=' || r == '-' ||
		r == '+' || r == '*' || r == '&' || r == '_' ||
		r == '@' || r == '^' || r == '~' || r == ':' ||
		r == '.' || r == '%' || r == '/' || r == '!' ||
		r == '?' || r == '#' || unicode.IsLetter(r) || unicode.IsDigit(r)
}
