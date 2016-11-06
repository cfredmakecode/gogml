package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
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
	items := make([]item, 1000)
	_, ch := lex("broken", string(f))
	for i := range ch {
		if i.typ != itemWhitespace {
			fmt.Printf("%v\n", i.String())
		}
		items = append(items, i)
		if i.typ == itemError {
			log.Fatalf("error; done")
		}
	}
	curIndent := 0
	var addIndent bool
	// poorly write back out indented version - we didn't really parse things yet so it's naive but sort of works
	for i := range items {
		if addIndent {
			fmt.Printf("%v", indent(curIndent))
			addIndent = false
		}
		//
		//
		// todo(caf): we don't pull out newlines in multline comments, so the output is wrong
		// need to setup some way of linking a single multiline comment together across multiple tokens
		// todo(caf): need to keep multi character operators like == together
		//
		//
		switch items[i].typ {
		case itemNewline:
			addIndent = true
			fmt.Printf("\n")
			continue
		case itemEOF:
			fmt.Printf("\n")
		case itemOpenCurly:
			curIndent++
		case itemCloseCurly:
			curIndent--
		}
		if len(strings.TrimSpace(items[i].val)) > 0 {
			fmt.Printf("%v ", strings.TrimSpace(items[i].val))
		}
	}
	if curIndent != 0 {
		// someone didn't close their curly brace somewhere. and has a case of mondays
		// todo(caf): detect unbalanced curlies not here
		log.Fatalf("unbalanced {}s")
	}
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
