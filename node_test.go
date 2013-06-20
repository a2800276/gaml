package gaml

import (
	"testing"
)

const f_furthest_child = `
%h1
	%h2
		%h3
			%h4
`

func TestFindFurthestChild(t *testing.T) {
	p := NewParserString(f_furthest_child)
	if _, err := p.Parse(); err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	node := p.rootNode.findFurthestChild()
	if node.name != "h4" {
		println(node.name)
		t.Fail()
	}
}
