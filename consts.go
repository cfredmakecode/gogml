package main

type itemType int

const (
	itemError itemType = iota
	itemEOF
	itemNumber
	itemText
	itemStartMultilineComment
	itemMultilineComment
	itemEndMultilineComment
	itemStartDocstringLineComment
	itemDocstringLineComment
	itemStartLineComment
	itemLineComment
	itemIdentifier
	itemWhitespace
	itemBuiltinFunction
	itemOpenBracket
	itemCloseBracket
	itemSemicolon
	itemPipe
	itemOpenParen
	itemCloseParen
	itemSingleQuote
	itemDoubleQuote
	itemString
	itemOpenCurly
	itemCloseCurly
	itemNewline

	keyword_beginning
	// keywords
	keywordIf
	keywordElse
	keywordRepeat
	keywordWith
	keywordWhile
	keywordDo
	keywordUntil
	keywordSwitch
	keywordCase
	keywordDefault
	keywordBreak
	keywordContinue
	keywordExit
	keywordVar
	keywordDiv
	keywordMod
	keywordOther
	keywordSelf
	keywordAll
	keywordNoone
	keywordEnum
	keywordUndefined
	keywordTrue
	keywordFalse
	keywordAnd
	keywordOr
	keywordXor
	keywordNot

	// built-in functions
	keywordString
	keywordReal
	keywordArray
	keywordBoolean
	keywordPtr

	keyword_end

	op_beginning

	compEq
	compNEq
	compLtEq
	compGtEq

	opOrEq
	opXorEq
	opAddEq
	opSubEq
	opInc
	opDec

	op_multichar_end // we'll need to iterate over all the operators that are multiple characters

	opEq
	opAdd
	opSub
	opMul
	opDiv
	opMod
	opXor

	op_end
)

var keywordStr = [...]string{
	// keywords
	keywordIf:        "if",
	keywordElse:      "else",
	keywordRepeat:    "repeat",
	keywordWith:      "with",
	keywordWhile:     "while",
	keywordDo:        "do",
	keywordUntil:     "until",
	keywordSwitch:    "switch",
	keywordCase:      "case",
	keywordDefault:   "default",
	keywordBreak:     "break",
	keywordContinue:  "continue",
	keywordExit:      "exit",
	keywordVar:       "var",
	keywordDiv:       "div",
	keywordMod:       "mod",
	keywordOther:     "other",
	keywordSelf:      "self",
	keywordAll:       "all",
	keywordNoone:     "noone",
	keywordEnum:      "enum",
	keywordUndefined: "undefined",
	keywordTrue:      "true",
	keywordFalse:     "false",
	keywordAnd:       "and",
	keywordOr:        "or",
	keywordXor:       "xor",
	keywordNot:       "not",

	// built-ins
	keywordString:  "string",
	keywordReal:    "real",
	keywordArray:   "array",
	keywordBoolean: "boolean",
	keywordPtr:     "ptr",

	// ops
	compEq:   "==",
	compNEq:  "!=",
	compLtEq: "<=",
	compGtEq: ">=",

	opOrEq:  "|=",
	opXorEq: "^=",
	opAddEq: "+=",
	opSubEq: "-=",
	opInc:   "++",
	opDec:   "--",

	opEq:  "=",
	opAdd: "+",
	opSub: "-",
	opMul: "*",
	opDiv: "/",
	opMod: "%",
	opXor: "^",
}

var keywords map[string]itemType
var ops map[string]itemType

func init() {
	keywords = make(map[string]itemType)
	ops = make(map[string]itemType)
	for i := keyword_beginning + 1; i < keyword_end; i++ {
		keywords[keywordStr[i]] = i
	}
	for i := op_beginning + 1; i < op_end; i++ {
		ops[keywordStr[i]] = i
	}
}

func lookup(identifier string) itemType {
	if kw, ok := keywords[identifier]; ok {
		return kw
	}
	return itemIdentifier
}
