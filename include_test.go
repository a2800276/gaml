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

const html_1 = `<!DOCTYPE html>
<html>
 <h1>
  Test Header
 </h1>
 <footer>
  This is a footer line.
 </footer>
</html>
`
const include_sample2 = `
!!!
%html
	%h1
		Test Header
		> included_sample
`

const html_2 = `<!DOCTYPE html>
<html>
 <h1>
  Test Header
  <footer>
   This is a footer line.
  </footer>
 </h1>
</html>
`
const include_sample3 = `
!!!
%html
	%h1
		> included_sample
		Test Header
`

const html_3 = `<!DOCTYPE html>
<html>
 <h1>
  <footer>
   This is a footer line.
  </footer>
  Test Header
 </h1>
</html>
`

func testInclude(t *testing.T, in string, expected string) {
	p := NewParserString(in)
	var loader DummyLoader
	p.IncludeLoader = loader

	if r, err := NewRenderer(p); err != nil {
		t.Error(err)
	} else {
		test_compare(t, r.ToHtmlString(), expected)
	}

}
func TestInclude(t *testing.T) {
	testInclude(t, include_sample, html_1)
	testInclude(t, include_sample2, html_2)
	testInclude(t, include_sample3, html_3)
}

/// TEST with filesystem loader!
