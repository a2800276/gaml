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
	"fmt"
	"strconv"
	"strings"
)

type stateFunc func() stateFunc

type gamlline struct {
	line      string
	r         rune
	node      *node
	parser    *Parser
	err       string
	value     bytes.Buffer
	attr_name string
	stateFunc stateFunc
}

func GamlLineFromString(s string) *gamlline {
	g := gamlline{line: s}
	g.stateFunc = g.initial
	return &g
}

// let `parser` determine whether it can skip this line.
func (g *gamlline) Empty() bool {
	return g.line == ""
}

const END = '\x00'

// this is where the magic happens...
func (g *gamlline) processIntoCurrentNode(p *Parser) (err error) {
	// node is the current node that we will be filling with content.
	// it's place in the hierarchy of nodes has been determined by
	// `Parser` using the indentation.
	g.node = p.currentNode
	g.parser = p

	if 0 == strings.Index(g.line, "!!!") {
		g.node.nodeType = DOCTYPE
		return
	}

	for _, r := range g.line {
		g.r = r
		if g.stateFunc = g.stateFunc(); g.stateFunc == nil {
			return p.Err(g.err)
		}
	}

	g.r = END
	final_state := g.stateFunc()
	if final_state == nil {
		return p.Err(fmt.Sprintf("implausible state! (%s)", g.state_func_to_string()))
	}
	return
}

func (g *gamlline) collectRune() {
	g.value.WriteRune(g.r)
}

func (g *gamlline) ok() stateFunc {
	return g.ok
}

func (g *gamlline) initial() stateFunc {
	switch g.r {
	case END:
		return nil
	case '%':
		g.node.nodeType = TAG
		return g.tagName
	case '.':
		g.node.nodeType = TAG
		return g.class
	case '#':
		g.node.nodeType = TAG
		return g.id
	case '>':
		g.node.nodeType = INC
		return g.include
	default:
		g.node.nodeType = TXT
		g.collectRune()
		return g.text
	}
}

// handle three similiar cases:
// %tag
// .class
// #id
// they are different in that:
// 1.) once the "end" of their state is reached, they either need to
//     fill in the nodes tag, class or id. This is handled by the `f`
//     argument
// 2.) The state they remain in, if they haven't reached the end of
//     their respective tag, class or id. This is passed as the `s`
//     argument.
//
// This is complicated, not elegant AND I'm not explaining it well.
// refer to the `tagName`, `class` and `id` functions below to see how
// this is called.

func (g *gamlline) baseTag(f func(), s stateFunc) stateFunc {
	switch g.r {
	case END:
		f()
		return g.ok
	case '.':
		f()
		return g.class
	case '#':
		f()
		return g.id
	case ' ':
		f()
		g.collectRune()
		return g.textOrAttribute
	case '(':
		f()
		return g.attributes
	case '{':
		g.collectRune()
		return g.openBrace(s)
	default:
		g.collectRune()
		return s
	}
}
func (g *gamlline) tagName() stateFunc {
	return g.baseTag(g.fillInName, g.tagName)
}
func (g *gamlline) class() stateFunc {
	return g.baseTag(g.fillInDivClass, g.class)
}
func (g *gamlline) id() stateFunc {
	return g.baseTag(g.fillInDivId, g.id)
}

func (g *gamlline) include() stateFunc {
	switch g.r {
	case END:
		g.parser.parseInclude(g.value.String())
		return g.include
	default:
		g.collectRune()
		return g.include
	}
}
func (g *gamlline) text() stateFunc {
	switch g.r {
	case END:
		g.node.text = g.value.String()
		return g.ok
	default:
		g.collectRune()
		return g.text
	}
}
func (g *gamlline) textOrAttribute() stateFunc {
	switch g.r {
	case END:
		newTextNode(g)
		g.node.text = g.value.String()
		return g.ok
	case ' ':
		g.collectRune()
		return g.textOrAttribute
	case '(':
		g.value.Reset()
		return g.attributes
	default: // default in all of these state functions is not correct, catch all sort of unicode crap people may throw at us ...
		g.collectRune()
		return g.textNew
	}
}
func (g *gamlline) textNew() stateFunc {
	newTextNode(g)
	switch g.r {
	case END:
		g.node.text = g.value.String()
		return g.ok
	default:
		g.collectRune()
		return g.text
	}
}

func newTextNode(g *gamlline) {
	node := newNode(g.node)
	node.nodeType = TXT
	g.node = node
}

func (g *gamlline) attributes() stateFunc {
	switch g.r {
	case END:
		return nil
	case ' ':
		return g.attributes
	case ')':
		return g.textNew
	default:
		g.collectRune()
		return g.attributesName
	}
}

func (g *gamlline) attributesName() stateFunc {
	switch g.r {
	case END:
		return nil
	case ' ', '=':
		g.attr_name = g.value.String()
		g.value.Reset()
		return g.attributesAfterName
	case ')':
		g.node.AddBooleanAttribute(g.value.String())
		g.value.Reset()
		return g.textNew
	default:
		g.collectRune()
		return g.attributesName
	}
}

func (g *gamlline) attributesAfterName() stateFunc {
	// this one is stupid.
	switch g.r {
	case END:
		g.node.AddBooleanAttribute(g.attr_name)
		return g.ok
	case ' ', '=': // <-- allows a ==    == = 'bla'
		return g.attributesAfterName
	//case '\'', '"': // <-- allows a = 'Bla"
	case '\'': // <-- allows only a = 'Bla'
		return g.attributesValues
	// valueless attribute, start of next attr or )
	case '"':
		g.err = "attribute values must be in single quote"
		return nil
	default:
		g.node.AddBooleanAttribute(g.attr_name)
		g.value.Reset()
		return g.attributes()
	}
}

func (g *gamlline) attributesValues() stateFunc {
	switch g.r {
	//case '"', '\'':
	case END:
		return nil
	case '\\':
		return g.valueEscape
	case '\'':
		g.node.AddAttribute(g.attr_name, g.value.String())
		g.value.Reset()
		return g.attributes
	default:
		g.collectRune()
		return g.attributesValues
	}
}

func (g *gamlline) valueEscape() stateFunc {
	switch g.r {
	case 'a':
		g.r = '\a' // bell
	case 'b':
		g.r = '\b' // backspace
	case 'f':
		g.r = '\f' // form feed
	case 'n':
		g.r = '\n' // newline
	case 'r':
		g.r = '\r' // cr
	case 't':
		g.r = '\t' // tab
	case 'v':
		g.r = '\v' // vertical tab
		// default is to literally pass through the rune after the backslash
		// no matter what the value.
	}
	g.collectRune()
	return g.attributesValues
}

func (g *gamlline) openBrace(s stateFunc) stateFunc {
	return func() stateFunc {
		switch g.r {
		case END:
			return nil
		case '{': // second {, collect literally
			g.collectRune()
			return g.passLiteral(s)
		default:
			return s()
		}
	}
}

func (g *gamlline) passLiteral(s stateFunc) stateFunc {
	return func() stateFunc {
		g.collectRune()
		switch g.r {
		case END:
			return nil
		case '}':
			return g.closeBrace(s)
		default:
			return g.passLiteral(s)
		}
	}
}
func (g *gamlline) closeBrace(s stateFunc) stateFunc {
	return func() stateFunc {
		g.collectRune()
		switch g.r {
		case END:
			return nil
		case '}':
			return s
		default:
			return g.passLiteral(s)
		}
	}
}

func (g *gamlline) fillInName() {
	g.node.name = g.value.String()
	g.value.Reset()
}

func (g *gamlline) fillInAttribute(attribute_name string) {
	if g.node.name == "" {
		g.node.name = "div"
	}
	g.node.AddAttribute(attribute_name, g.value.String())
	g.value.Reset()

}

func (g *gamlline) fillInDivClass() {
	g.fillInAttribute("class")
}

func (g *gamlline) fillInDivId() {
	g.fillInAttribute("id")
}

func (g *gamlline) String() string {
	str := "line: " + g.line
	str += " - r: " + strconv.QuoteRune(g.r) + " - state: "
	// a bit ridiculous
	return str + g.state_func_to_string()
}

func (g *gamlline) state_func_to_string() string {
	switch fmt.Sprintf("%p", g.stateFunc) {
	case fmt.Sprintf("%p", g.initial):
		return "initial"
	case fmt.Sprintf("%p", g.tagName):
		return "tagName"
	case fmt.Sprintf("%p", g.class):
		return "class"
	case fmt.Sprintf("%p", g.id):
		return "id"
	case fmt.Sprintf("%p", g.include):
		return "include"
	case fmt.Sprintf("%p", g.text):
		return "text"
	case fmt.Sprintf("%p", g.textOrAttribute):
		return "textOrAttribute"
	case fmt.Sprintf("%p", g.textNew):
		return "textNew"
	case fmt.Sprintf("%p", g.attributes):
		return "attributes"
	case fmt.Sprintf("%p", g.attributesName):
		return "attributesName"
	case fmt.Sprintf("%p", g.attributesAfterName):
		return "attributesAfterName"
	case fmt.Sprintf("%p", g.attributesValues):
		return "attributesValues"
	case fmt.Sprintf("%p", g.valueEscape):
		return "valueEscape"
	case fmt.Sprintf("%p", g.ok):
		return "ok"
	default:
		return fmt.Sprintf("%p", g.stateFunc)
	}
}
