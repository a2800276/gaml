package main

import (
	"gaml"
)

func main() {
	line := `
%html
	%head
		%body
			%h1
				Hello World`

	html,_ := gaml.GamlToHtml(line)
	println(html)
}
