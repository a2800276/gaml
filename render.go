package gaml

import (
	"bytes"
	"io"
	"strings"
)

type Renderer struct {
	Indent IndentFunc
}

type renderer struct {
	writer io.Writer
	opts   *Renderer
}

// add custom indentation if you want.
// take the level of indentation and write the indent to the writer.
// default is to not indent. If using the `GamlToHtml`
// utility method, the returned string will be indented
// with a space per indent.
type IndentFunc func(indentLevel int, writer io.Writer)

var DefaultIndentFunc = func(i int, w io.Writer) { return }
var IndentSpace = func(indent int, w io.Writer) {
	for i := 0; i != indent; i++ {
		io.WriteString(w, " ")
	}
}

func GamlToHtml(gaml string) (html string, err error) {
	var buffer bytes.Buffer

	renderer := NewRenderer()

	renderer.Indent = IndentSpace
	parser := NewParser(bytes.NewBufferString(gaml))

	if root, err2 := parser.Parse(); err2 != nil {
		return "", err2
	} else {
		renderer.ToHtml(root, &buffer)
	}
	return buffer.String(), nil
}

func NewRenderer() *Renderer {
	return &Renderer{Indent: DefaultIndentFunc}
}

func (r *Renderer) ToHtml(root *node, writer io.Writer) {
	rr := renderer{writer, r}
	for _, node := range root.children {
		rr.Render(node)
	}
}

// Write an html representation of this node to the specified `Writer`
// currently, there is no way to influence how the node will be
// rendered. Take it or leave it!
func (r *renderer) Render(n *node) {
	r.render(n, 0)
}

func (r *renderer) render(n *node, indent int) {
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

func (r *renderer) renderDocType(n *node) {
	// this is in it's own method so all the doctypes
	// can be collected here one fine day when different
	// rendering options are supported.
	io.WriteString(r.writer, "<!DOCTYPE html>\n")
}

func (r *renderer) renderTag(n *node, indent int) {
	r.opts.Indent(indent, r.writer)
	io.WriteString(r.writer, "<")
	io.WriteString(r.writer, n.name)

	r.renderAttributes(n)

	io.WriteString(r.writer, ">\n")

	if r.isVoid(n) {
		return
	}
	r.renderChildren(n, indent+1)

	r.opts.Indent(indent, r.writer)
	io.WriteString(r.writer, "</")
	io.WriteString(r.writer, n.name)
	io.WriteString(r.writer, ">\n") // what to do about the trailing \n !?
}

func (r *renderer) renderChildren(n *node, indent int) {
	for _, child := range n.children {
		r.render(child, indent)
	}
}

func (r *renderer) isVoid(n *node) bool {
	for _, name := range voidElements {
		if n.name == name {
			return true
		}
	}
	return false
}

func (r *renderer) renderAttributes(n *node) {
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

func (r *renderer) renderText(n *node, indent int) {
	// ditto: will probably want some options for escaping here.
	for i := 0; i != indent; i++ {
		io.WriteString(r.writer, " ")
	}
	io.WriteString(r.writer, n.text)
	io.WriteString(r.writer, "\n")

}
