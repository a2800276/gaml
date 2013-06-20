package gaml

import (
	"bytes"
	"fmt"
	"testing"
)

type DummyLoader struct{}

func (d DummyLoader) Load(id interface{}) (p *Parser, err error) {
	// check string and equal to expected
	if str, ok := id.(string); !ok {
		return nil, fmt.Errorf("!? not a string")
	} else {
		if str != "included_sample" && str != "included_sample2" {
			return nil, fmt.Errorf("!? wrong string %s", str)
		}
		switch str {
		case "included_sample":
			p = NewParserString(included_sample)
		case "included_sample2":
			p = NewParserString(included_sample2)
		}
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

const include_sample4 = `
> included_sample
%whatever
`

const html_4 = `<footer>
 This is a footer line.
</footer>
<whatever>
</whatever>
`

const include_sample5 = `
> included_sample
	%whatever
`

const html_5 = `<footer>
 This is a footer line.
 <whatever>
 </whatever>
</footer>
`
const included_sample2 = `
%h1
	%h2
		%h3
			%h4
`
const include_sample6 = `
> included_sample2
%whatever
`

const html_6 = `<h1>
 <h2>
  <h3>
   <h4>
   </h4>
  </h3>
 </h2>
</h1>
<whatever>
</whatever>
`
const include_sample7 = `
> included_sample2
	%whatever
`
const html_7 = `<h1>
 <h2>
  <h3>
   <h4>
    <whatever>
    </whatever>
   </h4>
  </h3>
 </h2>
</h1>
`
const include_sample8 = `
%html
	> included_sample2
	%whatever
`
const html8 = `<html>
 <h1>
  <h2>
   <h3>
    <h4>
    </h4>
   </h3>
  </h2>
 </h1>
 <whatever>
 </whatever>
</html>
`

func testInclude(t *testing.T, in string, expected string) {
	p := NewParserString(in)
	var loader DummyLoader
	p.IncludeLoader = loader
	var output bytes.Buffer

	if r, err := NewRenderer(p, &output); err != nil {
		t.Error(err)
	} else {
		r.ToHtml()
		test_compare(t, output.String(), expected)
	}

}
func TestInclude(t *testing.T) {
	testInclude(t, include_sample, html_1)
	testInclude(t, include_sample2, html_2)
	testInclude(t, include_sample3, html_3)
}

func TestFirstLineInclude(t *testing.T) {
	testInclude(t, include_sample4, html_4)
	testInclude(t, include_sample5, html_5)
	testInclude(t, include_sample6, html_6)
	testInclude(t, include_sample7, html_7)
	testInclude(t, include_sample8, html8)

}

func TestTest(t *testing.T) {
	testInclude(t, include_sample5, html_5)
}

/// TEST with filesystem loader!
