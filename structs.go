package main

import (
	"fmt"
	"io"
	"log"
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

type lexer struct {
	stream *strings.Reader
}

func newLexer(s string) *lexer {
	l := &lexer{}
	l.stream = strings.NewReader(s)
	return l
}

func (l *lexer) nextToken() (token, error) {
	var tok token
	var doneOnce bool
	for {
		r, _, err := l.stream.ReadRune()
		if err != nil {
			return tok, err
		}

		if doneOnce && tok.tokType != getTokenType(r) {
			l.stream.UnreadRune()
			break
		}
		tok.tokType = getTokenType(r)
		tok.value += string(r)
		tok.name = tok.tokType.String()
		doneOnce = true
	}
	return tok, nil
}

func pants(s string) {
	tokens := make([]token, 2048)
	l := newLexer(s)
	for {
		tok, err := l.nextToken()
		tokens = append(tokens, tok)
		if err != nil {
			if err != io.EOF {
				log.Fatalf("failed to get token")
			}
			// fmt.Printf("%#v\n", tok)
			break
		}
		// fmt.Printf("%#v\n", tok)
	}

	for _, t := range tokens {
		if t.tokType != TOKEN_WHITESPACE {
			fmt.Printf("%#v\n", t)
		}
	}
}
