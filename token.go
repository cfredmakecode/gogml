package main

import "unicode"

//go:generate stringer -type=tokenType
type tokenType int

const (
	TOKEN_UNKNOWN tokenType = iota
	TOKEN_WHITESPACE
	TOKEN_NEWLINE
	TOKEN_ALPHA
	TOKEN_DIGIT
	TOKEN_IDENTIFIER
	TOKEN_NUMBER
	// let's just go in order of whats on the keyboard?
	TOKEN_TILDE
	TOKEN_BACKTICK
	TOKEN_EXCLAMATION
	TOKEN_AT
	TOKEN_HASH
	TOKEN_DOLLAR
	TOKEN_PERCENT
	TOKEN_CARET
	TOKEN_AMP
	TOKEN_ASTERISK
	TOKEN_OPENPAREN
	TOKEN_CLOSEPAREN
	TOKEN_UNDERSCORE
	TOKEN_DASH
	TOKEN_PLUS
	TOKEN_EQUALS
	TOKEN_OPENCURLY
	TOKEN_CLOSECURLY
	TOKEN_OPENBRACKET
	TOKEN_CLOSEBRACKET
	TOKEN_PIPE
	TOKEN_BACKSLASH
	TOKEN_COLON
	TOKEN_SEMICOLON
	TOKEN_DOUBLEQUOTE
	TOKEN_QUOTE
	TOKEN_LT
	TOKEN_GT
	TOKEN_COMMA
	TOKEN_PERIOD
	TOKEN_QUESTION
	TOKEN_FORWARDSLASH
)

func getTokenType(r rune) tokenType {
	switch {
	case unicode.IsSpace(r):
		if r == '\n' {
			return TOKEN_NEWLINE
		} else {
			return TOKEN_WHITESPACE
		}
	case unicode.IsLetter(r):
		return TOKEN_ALPHA
	case unicode.IsDigit(r):
		return TOKEN_DIGIT
	}

	switch r {
	case '~':
		return TOKEN_TILDE
	case '`':
		return TOKEN_BACKTICK
	case '!':
		return TOKEN_EXCLAMATION
	case '@':
		return TOKEN_AT
	case '#':
		return TOKEN_HASH
	case '$':
		return TOKEN_DOLLAR
	case '%':
		return TOKEN_PERCENT
	case '^':
		return TOKEN_CARET
	case '&':
		return TOKEN_AMP
	case '*':
		return TOKEN_ASTERISK
	case '(':
		return TOKEN_OPENPAREN
	case ')':
		return TOKEN_CLOSEPAREN
	case '_':
		return TOKEN_UNDERSCORE
	case '-':
		return TOKEN_DASH
	case '+':
		return TOKEN_PLUS
	case '=':
		return TOKEN_EQUALS
	case '{':
		return TOKEN_OPENCURLY
	case '}':
		return TOKEN_CLOSECURLY
	case '[':
		return TOKEN_OPENBRACKET
	case ']':
		return TOKEN_CLOSEBRACKET
	case '|':
		return TOKEN_PIPE
	case '\\':
		return TOKEN_BACKSLASH
	case ':':
		return TOKEN_COLON
	case ';':
		return TOKEN_SEMICOLON
	case '"':
		return TOKEN_DOUBLEQUOTE
	case '\'':
		return TOKEN_QUOTE
	case '<':
		return TOKEN_LT
	case '>':
		return TOKEN_GT
	case ',':
		return TOKEN_COMMA
	case '.':
		return TOKEN_PERIOD
	case '?':
		return TOKEN_QUESTION
	case '/':
		return TOKEN_FORWARDSLASH
	default:
		return TOKEN_UNKNOWN
	}
}
