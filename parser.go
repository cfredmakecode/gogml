package main

type nodeType int

const (
	NODE_UNKNOWN nodeType = iota
	NODE_PROGRAM
)

type astNode struct {
	name     string
	value    interface{}
	thisType nodeType
	children []astNode
}

type parser struct {
	// toks   []token
	offset int
}

func (p *parser) parseAST() error {
	return nil
}
