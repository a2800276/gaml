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
)

type node struct {
	parent     *node               // parent, root nodes have parent == nil
	children   []*node             // child nodes
	name       string              // name of tags if this is a tag
	attributes map[string][]string // attributes if tag
	text       string              // text if this is a text node.

	// SPECIAL TREAT !!!
	// if `node` represents a DOCTYPE node (!!!), `name` == `text` == "!!!"
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
	if n.attributes == nil {
		n.attributes = make(map[string][]string)
	}
	n.attributes[name] = append(n.attributes[name], value)
}

// Write an html representation of this node to the specified `Writer`
// currently, there is no way to influence how the node will be
// rendered. Take it or leave it!
func (n *node) Render(writer io.Writer) {
	n.render(writer, 0)
}

func (n *node) render(w io.Writer, indent int) {
	if n.name == n.text && n.text == "!!!" {
		n.renderDocType(w)
	} else if n.name != "" {
		n.renderTag(w, indent)
	} else {
		n.renderText(w, indent)
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

	for _, child := range n.children {
		child.render(w, indent+1)
	}

	indentfunc()
	io.WriteString(w, "</")
	io.WriteString(w, n.name)
	io.WriteString(w, ">\n") // what to do about the trailing \n !?
}

func (n *node) renderAttributes(w io.Writer) {
	// this one is pretty straightforward, may need some escaping at some point.
	// currently my "security model" is that gaml templates come from a trusted
	// source (namely myself) and will be sanitized.

	for key, values := range n.attributes {
		io.WriteString(w, " ")
		io.WriteString(w, key)
		io.WriteString(w, "='")
		io.WriteString(w, strings.Join(values, " "))
		io.WriteString(w, "'")
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
