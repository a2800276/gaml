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
	"fmt"

//	"io"
//	"strings"
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
}

func (n *node) String() string {
	return fmt.Sprintf("name: >%s< text >%s< parent:>%s<", n.name, n.text, n.parent)
}

// creates a new node with the specified parent.
// Specify `nil` for a root node. This method automatically
// adds the newly created node as a child of the parent.
func newNode(parent *node) *node {
	n := new(node)
	if parent != nil {
		parent.addChild(n)
	}
	return n
}

func newRoot() *node {
	return &node{nodeType: ROOT}
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

func (n *node) findFurthestChild() *node {
	// if the line following an include node is indented, the
	// final node of the include tree will serve as it's parent
	// this func is to locate that node.
	//
	// i.e.
	//     > head.gaml
	//       %h1 hello
	//
	// head.gaml contains:
	//    !!!
	//    html
	//      head
	//      body
	//
	// in the resulting html, <h1>hello</h1> should be a child of body.
	if 0 == len(n.children) {
		return n
	} else {
		return n.children[len(n.children)-1].findFurthestChild()
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
