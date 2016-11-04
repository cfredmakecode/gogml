package main

import (
	"fmt"
	"io"
	"strings"
)

type token struct {
	name    string
	value   string
	tokType tokenType
}

func (t *token) String() string {
	return fmt.Sprintf("token{type: %v, name: %#v, value: %#v}",
		t.tokType, t.name, t.value)
}

type lexerState int

const (
	LEXER_NONE lexerState = iota
	LEXER_IDENTIFIER
	LEXER_NUMBER
)

type lexer struct {
	stream *strings.Reader
	tokens []token
	state  lexerState
}

func newLexer(s string) *lexer {
	l := &lexer{}
	l.stream = strings.NewReader(s)
	return l
}

func (l *lexer) nextToken() (token, error) {
	var tok token
	r, _, err := l.stream.ReadRune()
	tok.tokType = getTokenType(r)
	tok.value = string(r)
	tok.name = tok.tokType.String()
	return tok, err
}

func (l *lexer) lex() error {
	for {
		tok, err := l.nextToken()
		l.tokens = append(l.tokens, tok)
		if err != nil {
			if err != io.EOF {
				return fmt.Errorf("failed to get token")
			}
			break
		}
	}
	l.coalsceTokens()
	return nil
}

// post-process the very basic per character tokens into
// identifiers, or multiple digit numbers (not yet floats, only consecutive digits here)
// this is a pretty crappy way of doing this
func (l *lexer) coalsceTokens() {
	processed := make([]token, 0, len(l.tokens)-1)
	inProgress := l.tokens[0]
	for _, t := range l.tokens[1 : len(l.tokens)-1] {
		switch l.state {
		case LEXER_NONE:
			processed = append(processed, inProgress)
			inProgress = t
		case LEXER_IDENTIFIER:
			if t.tokType == TOKEN_ALPHA || t.tokType == TOKEN_UNDERSCORE {
				inProgress.value += t.value
			} else {
				// call it an identifier just in time
				inProgress.tokType = TOKEN_IDENTIFIER
				inProgress.name = inProgress.tokType.String()
				// fmt.Printf("%#v\n", inProgress)
				processed = append(processed, inProgress)
				inProgress = t
			}
		case LEXER_NUMBER:
			if t.tokType == TOKEN_DIGIT {
				inProgress.value += t.value
			} else {
				inProgress.tokType = TOKEN_NUMBER
				inProgress.name = inProgress.tokType.String()
				// fmt.Printf("%#v\n", inProgress)
				processed = append(processed, inProgress)
				inProgress = t
			}
		}
		switch inProgress.tokType {
		case TOKEN_ALPHA, TOKEN_UNDERSCORE:
			l.state = LEXER_IDENTIFIER
		case TOKEN_DIGIT:
			l.state = LEXER_NUMBER
		default:
			l.state = LEXER_NONE
		}
	}
	l.tokens = processed
}
