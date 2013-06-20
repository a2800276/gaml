package gaml

import (
	"io"
	"testing"
)

func test_with_options(t *testing.T, r *Renderer, gaml string, should string) {
	if html, err := GamlToHtmlWithRenderer(gaml, r); err != nil {
		t.Error(err)
	} else {
		test_compare(t, html, should)
	}

}

func TestIndentFunc(t *testing.T) {
	f := func(i int, w io.Writer) {
		for ; i != 0; i-- {
			io.WriteString(w, "\t")
		}
	}
	renderer := &Renderer{f, true}
	test_with_options(t, renderer, "%html\n  %body", "<html>\n\t<body>\n\t</body>\n</html>\n")

	renderer = &Renderer{DefaultIndentFunc, true}
	test_with_options(t, renderer, "%html\n  %body2", "<html>\n<body2>\n</body2>\n</html>\n")

	renderer = &Renderer{DefaultIndentFunc, false}
	test_with_options(t, renderer, "%html\n  %body3", "<html><body3></body3></html>")
	test_with_options(t, renderer, "%html\n %h1\n  Hello World!\n  yeah.\n %h2", "<html><h1>Hello World!\nyeah.\n</h1><h2></h2></html>")
}
