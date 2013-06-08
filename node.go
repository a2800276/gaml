package gaml

// The contents of this file are meant for "internal use only"
// end users of gaml shouldn't have to deal with any of this.
// This file contains a Node datatype that keeps track of
// each xml/html node in the hierchary.
// Proper Elements/Tags and Text nodes are represented by the
// same type.
// Nodes are also responsible for rendering themselves. So any
// adjustment to the rendering would be done in here.
// At the moment, rendering only knows about tags that close, are
// followed by a newline and are indented one space per level.
// DOCTYPE nodes are only ever rendered to html5 (<!DOCTYPE html>)

import (
	"io"
	"strings"
	"fmt"
)

type nodeType int

const (
	UNKNOWN nodeType = iota
	DOCTYPE
	TAG
	TXT
	ROOT
	INC
)

//http://www.w3.org/TR/html5/syntax.html#void-elements
var voidElements = []string{"area", "base", "br", "col", "command", "embed", "hr",
	"img", "input", "keygen", "link", "meta", "param",
	"source", "track", "wbr"}

type node struct {
	parent     *node               // parent, root nodes have parent == nil
	children   []*node             // child nodes
	name       string              // name of tags if this is a tag
	attributes map[string][]string // attributes if tag
	text       string              // text if this is a text node.
	nodeType   nodeType

	// SPECIAL TREAT !!!
	// if `node` represents a DOCTYPE node (!!!), `name` == `text` == "!!!"
}

func (n * node) String () string {
	return fmt.Sprintf("name: >%s< text >%s< parent:>%s<", n.name, n.text, n.parent);
}

// creates a new node with the specified parent.
// Specify `nil` for a root node. This method automatically
// adds the newly created node as a child of the parent.
func newNode(parent *node) *node {
	n := new(node)
	n.parent = parent
	if parent != nil {
		parent.children = append(parent.children, n)
	}
	return n
}

// Appends a new value to the list of attributes of this node.
func (n *node) AddAttribute(name string, value string) {
	n.addAttribute(name, value)
}

func (n *node) AddBooleanAttribute(name string) {
	n.addAttribute(name, nil)
}

func (n *node) addAttribute(name string, value interface{}) {
	if n.attributes == nil {
		n.attributes = make(map[string][]string)
	}
	if value == nil {
		n.attributes[name] = nil
		return
	}
	n.attributes[name] = append(n.attributes[name], value.(string))
}


func (n * node) findFurthestChild()(*node) {
	if 0 == len(n.children) {
		return n
	} else {
		return n.children[len(n.children)-1]
	}
}

func (n *node) addChild(child *node) {
	switch {
		case n.nodeType == TAG || n.nodeType == ROOT:
			child.parent = n
			n.children = append(n.children, child)
		case n.nodeType == TXT:
			n.parent.addChild(child)
		case n.nodeType == INC:
			n.findFurthestChild().addChild(child)
		default:
			// nutin.
	}
}

// Write an html representation of this node to the specified `Writer`
// currently, there is no way to influence how the node will be
// rendered. Take it or leave it!
func (n *node) Render(writer io.Writer) {
	n.render(writer, 0)
}

func (n *node) render(w io.Writer, indent int) {
	switch {
		case n.nodeType == DOCTYPE:
			n.renderDocType(w)
	  case n.name == "" && n.text == "":
			// blank node (root, include)
			n.renderChildren(w, indent)
		case n.name != "":
			n.renderTag(w, indent)
		default:
			n.renderText(w,indent)
	}
}

func (n *node) renderDocType(w io.Writer) {
	// this is in it's own method so all the doctypes
	// can be collected here one fine day when different
	// rendering options are supported.
	io.WriteString(w, "<!DOCTYPE html>\n")
}

func (n *node) renderTag(w io.Writer, indent int) {
	indentfunc := func() {
		for i := 0; i != indent; i++ {
			io.WriteString(w, " ")
		}
	}
	indentfunc()
	io.WriteString(w, "<")
	io.WriteString(w, n.name)

	n.renderAttributes(w)

	io.WriteString(w, ">\n")

	if n.isVoid() {
		return
	}
	n.renderChildren(w, indent+1)

	indentfunc()
	io.WriteString(w, "</")
	io.WriteString(w, n.name)
	io.WriteString(w, ">\n") // what to do about the trailing \n !?
}

func (n *node) renderChildren(w io.Writer, indent int) {
	for _, child := range n.children {
		child.render(w, indent)
	}
}

func (n *node) isVoid() bool {
	for _, name := range voidElements {
		if n.name == name {
			return true
		}
	}
	return false
}

func (n *node) renderAttributes(w io.Writer) {
	// this one is pretty straightforward, may need some escaping at some point.
	// currently my "security model" is that gaml templates come from a trusted
	// source (namely myself) and will be sanitized.

	for key, values := range n.attributes {
		io.WriteString(w, " ")
		io.WriteString(w, key)
		if values != nil {
			io.WriteString(w, "='")
			io.WriteString(w, strings.Join(values, " "))
			io.WriteString(w, "'")
		}
	}
}

func (n *node) renderText(w io.Writer, indent int) {
	// ditto: will probably want some options for escaping here.
	for i := 0; i != indent; i++ {
		io.WriteString(w, " ")
	}
	io.WriteString(w, n.text)
	io.WriteString(w, "\n")

}

// internal use only!
func (n *node) render_debug(indent int) {
	for i := 0; i != indent; i++ {
		print("-")
	}
	println("+")
	for _, child := range n.children {
		child.render_debug(indent + 1)
	}
}
