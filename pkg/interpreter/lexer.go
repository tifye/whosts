package interpreter

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

type itemType int

const (
	itemError itemType = iota
	itemIP
	itemHost
	itemComment
	itemNewline
	itemEOF
)

const (
	commentStart        rune = '#'
	eof                      = -1
	spaceChars               = " \t\r\n"
	spaceCharsNoNewLine      = " \t"
	newLine                  = '\n'
)

type item struct {
	typ itemType
	val string
}

func (i item) String() string {
	switch i.typ {
	case itemEOF:
		return "EOF"
	case itemError:
		return i.val
	}
	if len(i.val) > 10 {
		return fmt.Sprintf("%.10q...", i.val)
	}
	return fmt.Sprintf("%q", i.val)
}

type lexer struct {
	input string
	start int
	pos   int
	items chan item
}

func lex(input string) (*lexer, <-chan item) {
	l := &lexer{
		input: input,
		items: make(chan item),
	}
	go l.run()
	return l, l.items
}

func (l *lexer) run() {
	for state := lexLine; state != nil; {
		state = state(l)
	}
	close(l.items)
}

func (l *lexer) emit(t itemType) stateFn {
	return l.emitItem(item{
		typ: t,
		val: l.input[l.start:l.pos],
	})
}

func (l *lexer) emitItem(i item) stateFn {
	l.items <- i
	l.start = l.pos
	return nil
}

func (l *lexer) next() rune {
	if l.pos >= len(l.input) {
		return eof
	}
	r, width := utf8.DecodeRuneInString(l.input[l.pos:])
	l.pos += width
	return r
}

func (l *lexer) accept(valid string) bool {
	if strings.ContainsRune(valid, l.next()) {
		return true
	}
	l.backup()
	return false
}

func (l *lexer) acceptRun(valid string) {
	r := l.next()
	for ; strings.ContainsRune(valid, r); r = l.next() {
	}
	if r != eof {
		l.backup()
	}
}

func (l *lexer) backup() {
	if l.pos <= 0 {
		return
	}

	_, width := utf8.DecodeLastRuneInString(l.input[:l.pos])
	l.pos -= width
}

func (l *lexer) peak() rune {
	r := l.next()
	if r != eof {
		l.backup()
	}
	return r
}

func (l *lexer) ignore() {
	l.start = l.pos
}

func (l *lexer) errorf(format string, args ...any) stateFn {
	i := item{
		typ: itemError,
		val: fmt.Sprintf(format, args...),
	}
	return l.emitItem(i)
}

func (l *lexer) scanOctet() bool {
	digits := "0123456789"

	if l.accept("0") {
		return !unicode.IsDigit(l.peak())
	}

	l.acceptRun(digits)
	return true
}

type stateFn func(*lexer) stateFn

func lexLine(l *lexer) stateFn {
	l.acceptRun(spaceCharsNoNewLine)
	l.ignore()

	switch l.peak() {
	case commentStart:
		return lexComment
	case newLine:
		return lexNewline
	case eof:
		return l.emit(itemEOF)
	default:
		return lexIP
	}
}

func lexIP(l *lexer) stateFn {
	for range 3 {
		if !l.scanOctet() {
			return l.errorf("invalid IP address, expected octet at %d", l.pos)
		}
		if !l.accept(".") {
			return l.errorf("invalid IP address, exptected '.' at %d", l.pos)
		}
	}
	if !l.scanOctet() {
		return l.errorf("invalid IP address, expected octet at %d", l.pos)
	}

	l.emit(itemIP)

	if !isSpace(l.peak()) {
		return l.errorf("expected space between IP and Host at %d", l.pos)
	}
	l.acceptRun(spaceChars)
	l.ignore()

	return lexHost
}

func lexHost(l *lexer) stateFn {
	for {
		r := l.next()
		if !isHostAllowed(r) {
			if r != eof {
				l.backup()
			}
			break
		}
	}

	if l.pos <= l.start {
		return l.errorf("missing Host definition at %d", l.pos)
	}

	l.emit(itemHost)

	if r := l.peak(); !isSpace(r) && r != eof {
		return l.errorf("unexpected character at %d, got: %q", l.pos, r)
	}

	l.acceptRun(spaceCharsNoNewLine)
	l.ignore()

	switch l.peak() {
	case commentStart:
		return lexComment
	case newLine:
		return lexNewline
	case eof:
		return l.emit(itemEOF)
	default:
		return l.errorf("unexpected character at %d, got: %q", l.pos, l.peak())
	}
}

func lexComment(l *lexer) stateFn {
	if !l.accept("#") {
		return l.errorf("expected '#' but got: %q", l.peak())
	}

	x := strings.Index(l.input[l.pos:], "\n")
	if x < 0 {
		l.pos = len(l.input)
		if l.pos > l.start {
			l.emit(itemComment)
		}
		return l.emit(itemEOF)
	}

	l.pos += x
	l.emit(itemComment)

	return lexNewline
}

func isSpace(r rune) bool {
	return r == ' ' || r == '\t' || r == '\n' || r == '\r'
}

func isHostAllowed(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '.' || r == '_' || r == '-' || r == '*'
}

func lexNewline(l *lexer) stateFn {
	if !l.accept("\n") {
		return l.errorf("expected new line but got: %q", l.peak())
	}
	l.emit(itemNewline)
	return lexLine
}
