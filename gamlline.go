package gaml

// not part of the public interface!
// end users should not have to deal with this code.

// basically, the gaml library works by: using the GO scanner to
// chop up the provided input into lines, the "parser" handles the indention,
// strips whitespace and comments, and the resulting stripped line with no comments
// is of type `gamlline`. It is a string that is capable of parsing itself into a
// `node`
// The logic to transform the string into a node is defined in a big, nasty, brutishly
// amateur statemachine which makes up the bulk of this file.
// the `gamlline` type is not public, so it can't be used outside this library (and would be
// fairly useless) Capitalized mathods indicate that they are meant to be used
// outside of this file (within the library), lower_case methods are meant as local helper functions.

import (
	"bytes"
	//"fmt"
	"strings"
)

type gstate int

type stateFunc func(*gamlline, rune)(stateFunc)

const (
	INITIAL gstate = iota
	TAG_NAME
	CLASS
	ID
	OPEN_BRACE
	PASS_LITERAL
	CLOSE_BRACE
	INCLUDE
	TEXT
	TEXT_OR_ATTRIBUTES
	TEXT_NEW // need to differentiate between a "pure" line of text and text that comes after a tag. (see below)
	ATTRIBUTES
	ATTRIBUTES_NAME
	ATTRIBUTES_AFTER_NAME
	ATTRIBUTES_VALUES
	ERR
)
type gamlline struct {
	line string
	node *node
	value bytes.Buffer
	attr_name string
	stateFunc stateFunc
	prevStateBrace gstate
}


func GamlLineFromString(s string)gamlline {
	g := gamlline{line: s}
	g.stateFunc = (*gamlline).initial
	g.prevStateBrace = ERR
	return g
}




// let `parser` determine whether it can skip this line.
func (g *gamlline) Empty() bool {
	return g.line == ""
}


// this is where the magic happens...
func (g *gamlline) processIntoCurrentNode(p *Parser) (err error) {
	// node is the current node that we will be filling with content.
	// it's place in the hierarchy of nodes has been determined by
	// `Parser` using the indentation.
	g.node = p.currentNode

	if 0 == strings.Index(g.line, "!!!") {
		g.node.nodeType = DOCTYPE
		return
	}

	for _, r := range g.line {
	println(g)
		g.stateFunc = g.stateFunc(g,r)
	}
	return
}

func (g *gamlline) initial (r rune) stateFunc {
	println(g)
	switch r {
  case '%':
		g.node.nodeType = TAG
		return (*gamlline).tagName
  case '.':
		g.node.nodeType = TAG
		return (*gamlline).class
  case '#':
		g.node.nodeType = TAG
		return (*gamlline).id
  case '>':
		g.node.nodeType = INC
		return (*gamlline).include
	default:
		g.node.nodeType = TXT
		g.value.WriteRune(r)
		return (*gamlline).text
	}
}

func (g *gamlline) tagName (r rune) stateFunc {
	switch r {
	case '.':
		g.fillInName()
		return (*gamlline).class
	case '#':
		g.fillInName()
		return (*gamlline).id
	case ' ':
		g.fillInName()
		g.value.WriteRune(r)
		return (*gamlline).textOrAttribute
	case '(':
		g.fillInName()
		return (*gamlline).attributes
	case '{':
		g.value.WriteRune(r)
		return (*gamlline).openBrace(g, (*gamlline).tagName)
	default:
		g.value.WriteRune(r)
		return (*gamlline).tagName
	}
}
func (g *gamlline) class (r rune) stateFunc {
	switch r {
	case '.':
		g.fillInDivClass()
		return (*gamlline).class
	case '#':
		g.fillInDivClass()
		return (*gamlline).id
	case ' ':
		g.fillInDivClass()
		g.value.WriteRune(r)
		return (*gamlline).textOrAttribute
	case '(':
		g.fillInDivClass()
		return (*gamlline).attributes
	case '{':
		g.value.WriteRune(r)
		return (*gamlline).openBrace(g, (*gamlline).class)
	default:
		g.value.WriteRune(r)
		return (*gamlline).class
	}
}
func (g *gamlline) id (r rune) stateFunc {
	switch r {
	case '.':
		g.fillInDivId()
		return (*gamlline).class
	case '#':
		g.fillInDivId()
		return (*gamlline).id
	case ' ':
		g.fillInDivId()
		g.value.WriteRune(r)
		return (*gamlline).textOrAttribute
	case '(':
		g.fillInDivId()
		return (*gamlline).attributes
	case '{':
		g.value.WriteRune(r)
		return (*gamlline).openBrace(g, (*gamlline).id)
	default:
		g.value.WriteRune(r)
		return (*gamlline).id
	}
}
func (g *gamlline) include (r rune) stateFunc {
	g.value.WriteRune(r)
	return (*gamlline).include
}
func (g *gamlline) text (r rune) stateFunc {
	g.value.WriteRune(r)
	return (*gamlline).text
}
func (g *gamlline) textOrAttribute (r rune) stateFunc {
	switch r {
	case ' ':
		g.value.WriteRune(r)
		return (*gamlline).textOrAttribute
	case '(':
		g.value.Reset()
		return (*gamlline).attributes
	default: // default in all of these state functions is not correct, catch all sort of unicode crap people may throw at us ...
		g.value.WriteRune(r)
		return (*gamlline).textNew
	}
}
func (g *gamlline) textNew (r rune) stateFunc {
	node := newNode(g.node)
	node.nodeType = TXT
	g.node = node
	g.value.WriteRune(r)
	return (*gamlline).text
}
func (g *gamlline) attributes (r rune) stateFunc {
	switch r {
	case ' ':
		return (*gamlline).attributes
	case ')':
		return (*gamlline).textNew
	default:
		g.value.WriteRune(r)
		return (*gamlline).attributesName
	}
}

func (g *gamlline) attributesName (r rune) stateFunc {
	switch r {
	case ' ', '=':
		g.attr_name = g.value.String()
		g.value.Reset()
		return (*gamlline).attributesAfterName
	case ')':
		g.node.AddBooleanAttribute(g.value.String())
		g.value.Reset()
		return (*gamlline).textNew
	default:
		g.value.WriteRune(r)
		return (*gamlline).attributesName
	}
}
func (g *gamlline) attributesAfterName (r rune) stateFunc {
	// this one is stupid.
	switch r {
	case ' ', '=': // <-- allows a ==    == = 'bla'
		return (*gamlline).attributesAfterName
	//case '\'', '"': // <-- allows a = 'Bla"
	case '\'': // <-- allows only a = 'Bla'
		return (*gamlline).attributesValues
	// valueless attribute, start of next attr or )
	default:
		g.node.AddBooleanAttribute(g.value.String())
		g.value.Reset()
		return (*gamlline).attributes(g,r)
	}
}

func (g *gamlline) attributesValues (r rune) stateFunc {
	switch r {
	//case '"', '\'':
	case '\'':
		g.node.AddAttribute(g.attr_name, g.value.String())
		g.value.Reset()
		return (*gamlline).attributes
	default:
		g.value.WriteRune(r)
		return (*gamlline).attributesValues
	}
}
func (g *gamlline) openBrace (s stateFunc) stateFunc {
	return func(g *gamlline, r rune)(stateFunc) {
		switch r {
		case '{': // second {, collect literally
			g.value.WriteRune(r)
			return (*gamlline).passLiteral(g,s)
		default:
			return s(g,r)
		}
	}
}


func (g *gamlline) passLiteral (s stateFunc) stateFunc {
	return func(g *gamlline, r rune)(stateFunc) {
		g.value.WriteRune(r)
		switch r {
		case '}':
			return (*gamlline).closeBrace(g,s)
		default:
			return (*gamlline).passLiteral(g,s)
		}
	}
}
func (g *gamlline) closeBrace (s stateFunc) stateFunc {
	return func(g *gamlline, r rune)(stateFunc) {
		g.value.WriteRune(r)
		switch r {
		case '}':
			return s
		default:
			return (*gamlline).passLiteral(g,s)
		}
	}
}



func (g *gamlline) fillInName() {
	g.node.name = g.value.String()
	g.value.Reset()
}

func (g *gamlline) fillInDivClass() {
	if g.node.name == "" {
		g.node.name = "div"
	}
	g.node.AddAttribute("class", g.value.String())
	g.value.Reset()
}
func (g *gamlline) fillInDivId() {
	if g.node.name == "" {
		g.node.name = "div"
	}
	g.node.AddAttribute("id", g.value.String())
	g.value.Reset()
}


func (s gstate) String() string {
	switch s {
	case INITIAL:
		return "INITIAL"
	case TAG_NAME:
		return "TAG_NAME"
	case CLASS:
		return "CLASS"
	case ID:
		return "ID"
	case OPEN_BRACE:
		return "OPEN_BRACE"
	case PASS_LITERAL:
		return "PASS_LITERAL"
	case CLOSE_BRACE:
		return "CLOSE_BRACE"
	case INCLUDE:
		return "INCLUDE"
	case TEXT:
		return "TEXT"
	case TEXT_OR_ATTRIBUTES:
		return "TEXT_OR_ATTRIBUTES"
	case TEXT_NEW:
		return "TEXT_NEW"
	case ATTRIBUTES:
		return "ATTRIBUTES"
	case ATTRIBUTES_NAME:
		return "ATTRIBUTES_NAME"
	case ATTRIBUTES_AFTER_NAME:
		return "ATTRIBUTES_AFTER_NAME"
	case ATTRIBUTES_VALUES:
		return "ATTRIBUTES_VALUES"
	case ERR:
		return "ERR"
	default:
		return "unknown state"
	}
}

