package main

// todo(caf): need to put together a list of GML built-ins
// and verify several details:
//  can an identifier start with underscore?
//  can there be multline strings?
//
type nodeType int

const (
	NODE_UNKNOWN nodeType = iota
)

type astNode struct {
	name     string
	value    interface{}
	thisType nodeType
	children []*astNode
}

type parser struct {
	l *lexer
	a *astNode
}

func (p *parser) parseAST() error {
	return nil
}
