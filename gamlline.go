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
	"strings"
)

type gamlline string

type gstate int

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

// let `parser` determine whether it can skip this line.
func (g gamlline) Empty() bool {
	return string(g) == ""
}

func (g *gamlline) fillCurrNode(p *Parser) (err error) {
	line := (string)(*g)
	node := p.currentNode
	var name bytes.Buffer
	switch line[0] {
	case '%':
		for _, r := range line[1:] {
			name.WriteRune(r)
		}
	default:
		return p.Err("unexpected char")
	}
	node.name = name.String()
	return nil
}

// this is where the magic happens...
func (g gamlline) ProcessIntoCurrentNode(p *Parser) (err error) {
	// utility/help : a string typedef means we can no longer create
	// slices using [:]. `line` is here to avoid having to cast all the
	// time.
	line := string(g)

	node := p.currentNode
	// node is the current node that we will be filling with content.
	// it's place in the hierarchy of nodes has been determined by
	// `Parser` using the indentation.

	if 0 == strings.Index(line, "!!!") {
		node.text = "!!!" // not nice!
		node.name = "!!!" // use text == name == "!!!" to signal doctype.
	}

	// this will contain the generic "value" we are collecting. This
	// could contain any number of things in the course of parsing a
	// line: the tag name, attribute name/values or text.
	var value bytes.Buffer

	fillInName := func() {
		node.name = value.String()
		value.Reset()
	}

	// some callback functions just to keep things confusing!
	add := func(attrN string) func() {
		return func() {
			if node.name == "" {
				node.name = "div"
			}
			node.AddAttribute(attrN, value.String())
			value.Reset()
		}
	}

	addClass := add("class")
	addId := add("id")

	// TEXT_NEW is a state that's reached for nodes that have a tag AND text
	// on the same line:
	// %h1 HEADING!
	// technically, the text node is a child of the tag node so we need
	// to swap things around ...
	textNew := func() {
		_node := newNode(node)
		node = _node
	}

	// state machine starts here.
	state          := INITIAL
	prevStateBrace := ERR


	var name string // remember name of name = attribute pairs

	for _, r := range line {
	REWIND:
		switch state {
		case INITIAL:
			state = initial(r, &value)
		case TAG_NAME:
			prevStateBrace = state
			state = tag(r, &value, fillInName, TAG_NAME)
		case CLASS:
			prevStateBrace = state
			state = tag(r, &value, addClass, CLASS)
		case ID:
			prevStateBrace = state
			state = tag(r, &value, addId, ID)
		case OPEN_BRACE:
			state = openBrace(r, &value, prevStateBrace)
			if state != PASS_LITERAL {
				goto REWIND
			}
		case PASS_LITERAL:
			state = passLiteral(r, &value)
		case CLOSE_BRACE:
			state = closeBrace(r, &value, prevStateBrace)
		case INCLUDE:
			// ignore for now ... ?
			node.text = "<!-- include not handled -->"
			return nil
		case TEXT:
			value.WriteRune(r)
		case TEXT_OR_ATTRIBUTES:
			// once the tag part of the node is through and we encounter
			// whitespace, it's not yet known whether attributes (a = 'b')
			// will follow or text.
			state = textOrAttribute(r, &value)
		case TEXT_NEW:
			textNew()
			value.WriteRune(r)
			state = TEXT
		case ATTRIBUTES:
			state = attributes(r, &value)
		case ATTRIBUTES_NAME:
			state, name = attributes_name(r, &value)
		case ATTRIBUTES_AFTER_NAME:
			if state, err = attributes_after_name(r, &value); err != nil {
				return err
			}
		case ATTRIBUTES_VALUES:
			state = attributes_values(r, &value, func() { node.AddAttribute(name, value.String()) })
		}
	}

	// stow away the value we have been collecting once we've
	// passed through the entire string. Go's bufio.Scanner throws out
	// the trailing \n, \r\n so the for loop above won't be called after
	// the last char, we need to clean up depending on which state we're in.
	switch state {
	case INITIAL:
		return p.Err("impossible state! (really: can't happen)")
	case TAG_NAME:
		fillInName()
	case CLASS:
		addClass()
	case ID:
		addId()
	case INCLUDE:
		// TODO
	case TEXT:
		node.text = value.String()
	case TEXT_OR_ATTRIBUTES, TEXT_NEW:
		textNew()
		node.text = value.String()
	case ATTRIBUTES, ATTRIBUTES_NAME, ATTRIBUTES_AFTER_NAME, ATTRIBUTES_VALUES:
		return p.Err("implausible state!")
	}
	return
}

// below are the state functions, more or less one per state.
// typically they return the subsequent state, but sometimes it was necessary to
// cheat.

func attributes_values(r rune, buf *bytes.Buffer, fillInValue func()) gstate {
	switch r {
	case '"', '\'':
		fillInValue()
		buf.Reset()
		return ATTRIBUTES
	default:
		buf.WriteRune(r)
		return ATTRIBUTES_VALUES
	}
}

func attributes_after_name(r rune, buf *bytes.Buffer) (gstate, error) {
	// this one is stupid.
	switch r {
	case ' ', '=': // <-- allows a ==    == = 'bla'
		return ATTRIBUTES_AFTER_NAME, nil
	case '\'', '"': // <-- allows a = 'Bla"
		return ATTRIBUTES_VALUES, nil
	default:
		return ERR, fmt.Errorf("unquoted attribute values")
	}
}

func attributes_name(r rune, buf *bytes.Buffer) (gstate, string) {
	switch r {
	case ' ', '=':
		name := buf.String()
		buf.Reset()
		return ATTRIBUTES_AFTER_NAME, name
	default:
		buf.WriteRune(r)
		return ATTRIBUTES_NAME, ""
	}
}
func attributes(r rune, buf *bytes.Buffer) gstate {
	switch r {
	case ' ':
		return ATTRIBUTES
	case ')':
		return TEXT
	default:
		buf.WriteRune(r)
		return ATTRIBUTES_NAME
	}
}

func textOrAttribute(r rune, buf *bytes.Buffer) gstate {
	switch r {
	case ' ':
		buf.WriteRune(r)
		return TEXT_OR_ATTRIBUTES
	case '(':
		buf.Reset()
		return ATTRIBUTES
	default: // default in all of these state functions is not correct, catch all sort of unicode crap people may throw at us ...
		buf.WriteRune(r)
		return TEXT_NEW
	}
}

func tag(r rune, buf *bytes.Buffer, fillInValue func(), state gstate) gstate {
	switch r {
	case '.':
		fillInValue()
		return CLASS
	case '#':
		fillInValue()
		return ID
	case ' ':
		fillInValue()
		buf.WriteRune(r)
		return TEXT_OR_ATTRIBUTES
	case '(':
		fillInValue()
		return ATTRIBUTES
	case '{':
		buf.WriteRune(r)
		return OPEN_BRACE
	default:
		buf.WriteRune(r)
		return state
	}
}

func openBrace(r rune, buf *bytes.Buffer, prev_state gstate) gstate {
	switch r {
		case '{': // second {, collect literally
			buf.WriteRune(r)
			return PASS_LITERAL
		default:
			return prev_state
	}
}
func passLiteral(r rune, buf *bytes.Buffer) gstate {
	buf.WriteRune(r)
	switch r {
		case '}':
			return CLOSE_BRACE
		default:
			return PASS_LITERAL
	}
}

func closeBrace(r rune, buf *bytes.Buffer, prev_state gstate) gstate {
	buf.WriteRune(r)
	// r was either a second }, so we're back out of the {{ }} block
	switch r {
		case '}':
			return prev_state
		default:
			return PASS_LITERAL
	}
}

func initial(r rune, buf *bytes.Buffer) gstate {
	switch r {
	case '%':
		return TAG_NAME
	case '.':
		return CLASS
	case '#':
		return ID
	case '>':
		return INCLUDE
	default:
		buf.WriteRune(r)
		return TEXT
	}
}
