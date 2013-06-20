package gaml

import (
	"bytes"
	"io"
	"strings"
)

type Renderer struct {
	nodes  []*node
	writer io.Writer
}

func GamlToHtml(gaml string) (html string, err error) {
	var renderer Renderer
	var buffer bytes.Buffer
	if renderer, err = NewRenderer(NewParser(bytes.NewBufferString(gaml)), &buffer); err != nil {
		return
	}
	renderer.ToHtml()
	return buffer.String(), nil
}

func NewRenderer(p *Parser, writer io.Writer) (r Renderer, err error) {
	if err = p.Parse(); err != nil {
		return
	}
	return Renderer{p.rootNode.children, writer}, nil
}

func (r *Renderer) ToHtml() {
	for _, node := range r.nodes {
		node.Render(r.writer)
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
