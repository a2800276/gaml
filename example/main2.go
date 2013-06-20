package main

import (
	"gaml"
	"bytes"
	"os"
)

// this configures the output options and can be reused
var renderer = gaml.NewRenderer()

func main() {
	line := `
%html
	%head
		%body
			%h1
				Hello World`

  reader := bytes.NewBufferString(line)

	// Construct a parser from an `io.Reader`, there is
	// also a convenienve method: NewParserString(string) 
	// we could have used ...

	parser := gaml.NewParser(reader)

	// the root node returned by `Parse` is an abstract 
	// represenation of the gaml. It too can be reused if
	// the underlying gaml does not change.

	root, _ := parser.Parse()

	// finally, render the abstract gaml
	renderer.ToHtml(root, os.Stdout)
}
