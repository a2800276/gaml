package gaml

import (
	"bytes"
	"io"
)

type Renderer struct {
	nodes []*node
}

func GamlToHtml(gaml string) (html string, err error) {
	var renderer Renderer
	if renderer, err = NewRenderer(NewParser(bytes.NewBufferString(gaml))); err != nil {
		return
	}
	return renderer.ToHtmlString(), nil
}

func NewRenderer(p *Parser) (r Renderer, err error) {
	if err = p.Parse(); err != nil {
		return
	}
	return Renderer{p.rootNode.children}, nil
}

func (r *Renderer) ToHtml(writer io.Writer) {
	for _, node := range r.nodes {
		node.Render(writer)
	}
}

func (r *Renderer) ToHtmlString() (html string) {
	var output bytes.Buffer
	r.ToHtml(&output)
	return output.String()
}
