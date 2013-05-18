package main

import (
	"gaml"
	"html/template"
	"fmt"
	"os"
)


const gaml_template_1 = `
!!!
%html
	%head
		%title {{.}}
	%body

`

type Person struct {
	First string
	Last  string
}

var People = []Person{Person{"Bob", "Marley"}, Person{"Peter", "Tosh"}, Person{"Bunny", "Wailer"}}

const gaml_template_2 = `
!!!
%html
	%head
	%body
		%ul
			{{ range . }}
			%li.{{"name"}} {{.First}} {{.Last}} 
			{{ end }}

`
// DOESN'T WORK!
const gaml_template_3 = `
%html
	%body
		%ul
			{{ range . }}
			<!-- currently doesn't work because of . in the id value !! -->
			%li#{{.First}}_{{.Last}} {{.First}} {{.Last}} 
			{{ end }}

`

func main () {
	html_t, err := gaml.GamlToHtml(gaml_template_1)
	if err != nil {
		fmt.Printf("error: %s", err.Error())
	}
	template,err := template.New("test_template").Parse(html_t)
	template.Execute(os.Stdout, "Hello World!")

	html_t, err = gaml.GamlToHtml(gaml_template_2)
	if err != nil {
		fmt.Printf("error: %s", err.Error())
	}
	println(html_t)
	template,err = template.New("test_template2").Parse(html_t)
	if err != nil {
		fmt.Printf("error: %s", err.Error())
	}
	template.Execute(os.Stdout, People)


}