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
		r.Render(node)
	}
}

// Write an html representation of this node to the specified `Writer`
// currently, there is no way to influence how the node will be
// rendered. Take it or leave it!
func (r *Renderer) Render(n *node) {
	r.render(n, 0)
}

func (r *Renderer) render(n *node, indent int) {
	switch {
	case n.nodeType == DOCTYPE:
		r.renderDocType(n)
	case n.name == "" && n.text == "":
		// blank node (root, include)
		r.renderChildren(n, indent)
	case n.name != "":
		r.renderTag(n, indent)
	default:
		r.renderText(n, indent)
	}
}

func (r *Renderer) renderDocType(n *node) {
	// this is in it's own method so all the doctypes
	// can be collected here one fine day when different
	// rendering options are supported.
	io.WriteString(r.writer, "<!DOCTYPE html>\n")
}

func (r *Renderer) renderTag(n *node, indent int) {
	indentfunc := func() {
		for i := 0; i != indent; i++ {
			io.WriteString(r.writer, " ")
		}
	}
	indentfunc()
	io.WriteString(r.writer, "<")
	io.WriteString(r.writer, n.name)

	r.renderAttributes(n)

	io.WriteString(r.writer, ">\n")

	if r.isVoid(n) {
		return
	}
	r.renderChildren(n, indent+1)

	indentfunc()
	io.WriteString(r.writer, "</")
	io.WriteString(r.writer, n.name)
	io.WriteString(r.writer, ">\n") // what to do about the trailing \n !?
}

func (r *Renderer) renderChildren(n *node, indent int) {
	for _, child := range n.children {
		r.render(child, indent)
	}
}

func (r *Renderer) isVoid(n *node) bool {
	for _, name := range voidElements {
		if n.name == name {
			return true
		}
	}
	return false
}

func (r *Renderer) renderAttributes(n *node) {
	// this one is pretty straightforward, may need some escaping at some point.
	// currently my "security model" is that gaml templates come from a trusted
	// source (namely myself) and will be sanitized.

	for key, values := range n.attributes {
		io.WriteString(r.writer, " ")
		io.WriteString(r.writer, key)
		if values != nil {
			io.WriteString(r.writer, "='")
			io.WriteString(r.writer, strings.Join(values, " "))
			io.WriteString(r.writer, "'")
		}
	}
}

func (r *Renderer) renderText(n *node, indent int) {
	// ditto: will probably want some options for escaping here.
	for i := 0; i != indent; i++ {
		io.WriteString(r.writer, " ")
	}
	io.WriteString(r.writer, n.text)
	io.WriteString(r.writer, "\n")

}
