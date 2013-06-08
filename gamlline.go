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

// let `parser` determine whether it can skip this line.
func (g gamlline) Empty() bool {
	return string(g) == ""
}

// this is where the magic happens...
func (g gamlline) processIntoCurrentNode(p *Parser) (err error) {
	// utility/help : a string typedef means we can no longer create
	// slices using [:]. `line` is here to avoid having to cast all the
	// time.
	line := string(g)

	node := p.currentNode
	// node is the current node that we will be filling with content.
	// it's place in the hierarchy of nodes has been determined by
	// `Parser` using the indentation.

	if 0 == strings.Index(line, "!!!") {
		node.nodeType = DOCTYPE
		return
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
		_node.nodeType = TXT
		node = _node
	}

	// state machine starts here.
	state := INITIAL
	prevStateBrace := ERR

	var name string // remember name of name = attribute pairs

	for _, r := range line {
	REWIND:
		switch state {
		case INITIAL:
			state = initial(r, node, &value)
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
			//node.text = "<!-- include not handled -->"
			//return nil
			value.WriteRune(r)
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
			if state = attributes_after_name(r, &value); (state != ATTRIBUTES_AFTER_NAME) && (state != ATTRIBUTES_VALUES) {
				node.AddBooleanAttribute(name)
				goto REWIND
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
		err = p.parseInclude(value.String())
	case TEXT:
		node.text = value.String()
	case TEXT_OR_ATTRIBUTES, TEXT_NEW:
		textNew()
		node.text = value.String()
	case ATTRIBUTES_AFTER_NAME:
		node.AddBooleanAttribute(name)
	default:
		return p.Err(fmt.Sprintf("implausible state! (%s)", state.String()))
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

func attributes_after_name(r rune, buf *bytes.Buffer) gstate {
	// this one is stupid.
	switch r {
	case ' ', '=': // <-- allows a ==    == = 'bla'
		return ATTRIBUTES_AFTER_NAME
	case '\'', '"': // <-- allows a = 'Bla"
		return ATTRIBUTES_VALUES
	// valueless attribute, start of next attr or )
	default:
		return ATTRIBUTES
	}
}

func attributes_name(r rune, buf *bytes.Buffer) (gstate, string) {
	switch r {
	case ' ', '=', ')':
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

func initial(r rune, node *node, buf *bytes.Buffer) gstate {
	switch r {
	case '%':
		node.nodeType = TAG
		return TAG_NAME
	case '.':
		node.nodeType = TAG
		return CLASS
	case '#':
		node.nodeType = TAG
		return ID
	case '>':
		node.nodeType = INC
		return INCLUDE
	default:
		node.nodeType = TXT
		buf.WriteRune(r)
		return TEXT
	}
}

// the final three functions handle go template style braces.
// since '.' tends to come up in go templates quite a lot
// and -on the other hand- has a special meaning in html/g(h)aml
// it require special handling.
// In other cases there is no escaping of ' " . # %
// A dot (.) simply can't be part of a tag name, id or class (*)
// BUT, consider the following:
//
//    %tag#{{.calculated_id}}
//
// This would result in
//
//    <tag id="{{" class="calculated_id}}"> ...
//
// So anything following an opening {{ up until }} is passed through.
// we do need to keep track of the state we were previously in,
// though (TAG, CLASS, ID) to return to it after handling {{}}.
//
// (*) dot (.) can be part of an id or class name if they aren't
// defined in shortcut notation, i.e.:
//
//    %tag(id=".")
//
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
