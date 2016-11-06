package main

//go:generate C:\code\go\bin\stringer.exe -type=itemType

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

const eof = rune(0)

type stateFn func(*lexer) stateFn
type item struct {
	typ itemType
	val string
}

func (i *item) String() string {
	switch i.typ {
	case itemEOF:
		return "EOF"
	case itemError:
		return i.val
	}
	return fmt.Sprintf("%v:%#v", i.typ.String(), i.val)
}

type lexer struct {
	name  string
	input string
	start int
	pos   int
	width int
	items chan item
}

func lex(name, input string) (*lexer, chan item) {
	l := &lexer{
		name:  name,
		input: input,
		items: make(chan item),
	}
	go l.run()
	return l, l.items
}

func (l *lexer) emit(t itemType) {
	l.items <- item{t, l.input[l.start:l.pos]}
	l.start = l.pos
}

func (l *lexer) ignore() {
	l.start = l.pos
}
func (l *lexer) backup() {
	l.pos -= l.width
}
func (l *lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}
func (l *lexer) accept(validChars string) bool {
	if strings.IndexRune(validChars, l.next()) >= 0 {
		return true
	}
	l.backup()
	return false
}

// allow us to consume until we see something not in validChar
func (l *lexer) acceptRunOfRunes(validChars string) {
	for strings.IndexRune(validChars, l.next()) >= 0 {
	}
	l.backup()
}

func (l *lexer) next() (r rune) {
	if l.pos >= len(l.input) {
		l.width = 0
		return eof
	}
	r, l.width = utf8.DecodeRuneInString(l.input[l.pos:])
	l.pos += l.width
	return r
}

func lexWhitespace(l *lexer) stateFn {
	for {
		cur := l.peek()
		if cur == '\n' {
			l.emit(itemWhitespace)
			l.next()
			l.emit(itemNewline)
			return lexText
		}
		if !unicode.IsSpace(cur) {
			break
		}
		l.next()
	}
	l.emit(itemWhitespace)
	return lexText
}

// a default starting state
func lexText(l *lexer) stateFn {
	// check multiple character items before anything
	//
	//
	// todo(caf): generically check all ops and handle first
	// check keywords only after getting an identifier - or compare later in parsing?
	// for tokenizing step that'll be sufficient, we'll handle any nontrivial grouping / arg parsing in ... parsing
	//
	//
	if strings.HasPrefix(l.input[l.pos:], "/*") {
		if l.pos > l.start {
			l.emit(itemText)
		}
		return lexStartMultilineComment
	}
	if strings.HasPrefix(l.input[l.pos:], "///") {
		if l.pos > l.start {
			l.emit(itemText)
		}
		return lexStartDocstringLineComment
	}
	if strings.HasPrefix(l.input[l.pos:], "//") {
		if l.pos > l.start {
			l.emit(itemText)
		}
		return lexStartLineComment
	}
	//
	//
	//
	//

	cur := l.next()
	if cur == eof {
		if l.pos > l.start {
			l.emit(itemText)
		}
		l.emit(itemEOF)
		return nil
	}

	// match >1 char operators
	// note(caf): THE SEQUENCE OF CHECKS IS IMPORTANT, and stopping at multichar_end too
	// we need to check longest rune sequences to shortest so we don't identify a substring instead
	for i := op_beginning + 1; i < op_multichar_end; i++ {
		if strings.HasPrefix(l.input[l.start:], keywordStr[i]) {
			l.pos += len(keywordStr[i])
			l.emit(i)
			return lexText
		}
	}

	switch cur {
	case '-':
		// TODO(CAF): we interpret sequences like
		//  x+7-y-9
		// as
		//  x + 7 -y -9
		// we might punt the issue off to the parser (detect no operator between two number/identifiers && precending minus sign, strip sign and create sub op node)
		//
		// could be a number, and could be a float with no leading zero
		// todo(caf): maybe just have lexNumber decide this anyway..
		if unicode.IsNumber(l.peek()) || l.peek() == '.' {
			return lexNumber
		}
		l.emit(opSub)
		return lexText
	case '+':
		if unicode.IsNumber(l.peek()) || l.peek() == '.' {
			return lexNumber
		}
		l.emit(opAdd)
	case ';':
		l.emit(itemSemicolon)
	case '(':
		l.emit(itemOpenParen)
	case ')':
		l.emit(itemCloseParen)
	case '{':
		l.emit(itemOpenCurly)
	case '}':
		l.emit(itemCloseCurly)
	case '\n':
		l.emit(itemNewline)
	case '\r': // always completely strip carriage returns, we can put em back later based on configuration
		l.ignore()
		return lexText
	case '"':
		return lexString('"', l)
	case '\'':
		return lexString('\'', l)
	}

	if unicode.IsSpace(cur) {
		return lexWhitespace
	}
	if unicode.IsNumber(cur) || (cur == '.' && unicode.IsNumber(l.peek())) {
		return lexNumber
	}
	if unicode.IsLetter(cur) {
		return lexIdentifier
	}

	if l.pos > l.start {
		l.emit(itemText)
	}
	return lexText
}

func lexStartLineComment(l *lexer) stateFn {
	l.pos += len("//") // skip it.
	l.emit(itemStartLineComment)
	return lexLineComment
}

func lexLineComment(l *lexer) stateFn {
	if strings.HasPrefix(l.input[l.pos:], "\r\n") || strings.HasPrefix(l.input[l.pos:], "\n") {
		if l.pos > l.start {
			l.emit(itemLineComment)
		}
		return lexText
	}
	cur := l.next()
	if cur == eof {
		return lexText
	}
	return lexLineComment
}

func lexStartDocstringLineComment(l *lexer) stateFn {
	l.pos += len("///") // skip it.
	l.emit(itemStartDocstringLineComment)
	return lexDocstringLineComment
}

func lexDocstringLineComment(l *lexer) stateFn {
	if strings.HasPrefix(l.input[l.pos:], "\r") || strings.HasPrefix(l.input[l.pos:], "\n") {
		if l.pos > l.start {
			l.emit(itemDocstringLineComment)
		}
		return lexText
	}
	cur := l.next()
	if cur == eof {
		return lexText
	}
	return lexDocstringLineComment
}

func lexStartMultilineComment(l *lexer) stateFn {
	l.pos += len("/*") // skip it.
	l.emit(itemStartMultilineComment)
	return lexMultilineComment
}
func lexEndMultilineComment(l *lexer) stateFn {
	l.pos += len("*/") // skip it.
	l.emit(itemEndMultilineComment)
	return lexText
}

func lexMultilineComment(l *lexer) stateFn {
	if strings.HasPrefix(l.input[l.pos:], "*/") {
		if l.pos > l.start {
			l.emit(itemMultilineComment)
		}
		return lexEndMultilineComment
	}
	cur := l.next()
	if cur == eof {
		// todo decide if we care.. probably no?
		return lexText // l.errorf("did not close multiline comment before EOF")
	}
	return lexMultilineComment
}

// func lexOpenBracket(l *lexer) stateFn {
// 	l.pos += len("[")
// 	l.emit(itemOpenBracket)
// 	return lexInsideBrackets
// }
// func lexCloseBracket(l *lexer) stateFn {
// 	l.pos += len("]")
// 	l.emit(itemCloseBracket)
// 	return lexText
// }

// func lexInsideBrackets(l *lexer) stateFn {
// 	for {
// 		if strings.HasPrefix(l.input[l.pos:], "]") {
// 			return lexCloseBracket
// 		}
// 		switch r := l.next(); {
// 		case r == eof || r == '\n':
// 			return l.errorf("unclosed brackets")
// 		case unicode.IsSpace(r):
// 			l.ignore()
// 		case r == '|':
// 			l.emit(itemPipe)
// 		case r == '+' || r == '-' || '0' <= r && r <= '9':
// 			l.backup()
// 			return lexNumber
// 		case isAlphaNumeric(r):
// 			l.backup()
// 			return lexIdentifier
// 		default:
// 			return l.errorf("unhandled symbol inside brackets")
// 		}
// 	}
// }

func lexString(quoteCharacter rune, l *lexer) stateFn {
	if quoteCharacter == '\'' {
		l.emit(itemSingleQuote)
	} else {
		l.emit(itemDoubleQuote)
	}
	for {
		cur := l.next()
		if cur == eof {
			return l.errorf("unterminated string")
		}
		if cur == '\\' {
			if l.peek() == quoteCharacter { // ignore escaped versions
				l.pos += 2
				continue
			} else {
				l.emit(itemString)
				l.next()
				break
			}
		}
		if cur == quoteCharacter {
			l.backup()
			l.emit(itemString)
			l.next()
			break
		}
	}
	if quoteCharacter == '\'' {
		l.emit(itemSingleQuote)
	} else {
		l.emit(itemDoubleQuote)
	}
	return lexText
}

func lexIdentifier(l *lexer) stateFn {
	// implicitly disallows preceding underscore - not legal in GML anyway
	l.acceptRunOfRunes("abcdefghijklmnopqrtstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_0123456789")
	maybeKeyword := lookup(l.input[l.start:l.pos])
	l.emit(maybeKeyword)
	return lexText
}

func lexNumber(l *lexer) stateFn {
	// could be leading sign
	// note(caf): we could also be a float with no leading zero..
	l.accept("+-.")
	// hex?
	allowDigits := "0123456789"
	if l.accept("0") && l.accept("xX") {
		allowDigits = "01234567890abcdefABCDEF"
	}
	l.acceptRunOfRunes(allowDigits)
	if l.accept(".") { // naively this allows multiple periods in the case of a starting period
		l.acceptRunOfRunes(allowDigits)
	}
	// todo(caf) any other odd formats GML supports?
	l.emit(itemNumber)
	return lexText
}

func (l *lexer) errorf(format string, args ...interface{}) stateFn {
	l.items <- item{
		itemError,
		fmt.Sprintf(format, args...),
	}
	return nil
}

func isAlphaNumeric(r rune) bool {
	return unicode.IsNumber(r) || unicode.IsLetter(r)
}

func (l *lexer) run() {
	for state := lexText; state != nil; {
		state = state(l)
	}
	close(l.items)
}
