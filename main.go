package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"unicode"
)

const spacesPerIndentLevel = 4

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("missing filename.gml")
	}
	f, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Fatalf("failed to open %v", os.Args[1])
	}
	pants(string(f))
	//fmt.Fprintf(os.Stdout, fixit(f))
}

func pants(s string) {
	l := newLexer(s)
	err := l.lex()
	if err != nil {
		log.Fatalf("problem parsing tokens from text stream: %v", err.Error())
	}
	for _, t := range l.tokens {
		if t.tokType != TOKEN_WHITESPACE {
			fmt.Printf("%#v\n", t)
		}
	}
	p := &parser{l: l}
	_ = p
}

func fixit(r io.Reader) string {
	currentIndent := 0
	lastIndent := 0
	commentDepth := 0
	quoteDepth := 0
	_ = quoteDepth
	var out string
	s := bufio.NewScanner(r)
	for s.Scan() {
		// find our tab or space characters at the beginning of the line
		whitespace := 0
		lastIndent = currentIndent
		doneBeginning := false
		pastSingleLineComment := false
		for i, char := range s.Text() {
			if !doneBeginning && unicode.IsSpace(char) {
				whitespace++
				continue
			}
			doneBeginning = true
			if !pastSingleLineComment && char == '/' && len(s.Text())-1 > i && s.Text()[i+1] == '/' {
				continue
			}
			pastSingleLineComment = true
			if char == '/' && len(s.Text())-1 > i && s.Text()[i+1] == '*' {
				commentDepth++
			}
			if char == '*' && len(s.Text())-1 > i && s.Text()[i+1] == '/' {
				commentDepth--
			}
			// todo(caf): detect and ignore when the `{` or `}` are in a comment, or in a string
			if commentDepth == 0 && quoteDepth == 0 {
				if char == '{' {
					currentIndent++
				} else if char == '}' {
					currentIndent--
				}
			}
			// todo(caf): force if/else expressions to have gofmt style curly braces
			// including when they're on the same line
			// or start on the next line
			// or missing altogether (single statement conditional expression)
		}
		// out += fmt.Sprintf("//indent %v spaces %v prev %v commentdepth %v\n", currentIndent, whitespace, lastIndent, commentDepth)
		// don't touch the user's spacing inside comments or quotes
		if commentDepth != 0 {
			out += fmt.Sprintf("%v%v", indent(lastIndent), s.Text())
		} else {
			out += fmt.Sprintf("%v%v", indent(lastIndent), strings.TrimSpace(s.Text()))
		}
		// todo(caf): add missing semicolons
		// todo(caf): indent over once if the user opened a paren on the previous line end
		if s.Err() != nil {
			break
		}
		out += "\n"
	}
	return out
}

func indent(level int) string {
	var out string
	for _ = range make([]interface{}, level) {
		for _ = range make([]interface{}, spacesPerIndentLevel) {
			out += " "
		}
	}
	return out
}
