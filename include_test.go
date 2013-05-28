package gaml

import (
	"fmt"
	"testing"
)

type DummyLoader struct{}

func (d DummyLoader) Load(id interface{}) (p *Parser, err error) {
	// check string and equal to expected
	if str, ok := id.(string); !ok {
		return nil, fmt.Errorf("!? not a string")
	} else {
		if str != "included_sample" {
			return nil, fmt.Errorf("!? wrong string %s", str)
		}
		p = NewParserString(included_sample)
		return
	}

}

const included_sample = `
%footer
  This is a footer line.
`
const include_sample = `
!!!
%html
	%h1
		Test Header
	> included_sample
`
const include_sample2 = `
!!!
%html
	%h1
		> included_sample
`

func TestInclude(t *testing.T) {
	p := NewParserString(include_sample)
	var loader DummyLoader
	p.IncludeLoader = loader

	if r, err := NewRenderer(p); err != nil {
		t.Error(err)
	} else {
		print(r.ToHtmlString())
	}
}

/// TEST with filesystem loader!
