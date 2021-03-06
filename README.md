# About

Fairly haml-ish like html templating for Go.

Makes it easier to type up html by hand by avoiding having to type
superfluous angled bracked, closing tags, and some other redundancies.
Fans of Haml claim that it's not just a shortcut, but beautiful (like
Haikus) but they have a warped sense of aesthetics.

I'm going to assume you have a vague notion of what haml is so I'll keep
it short.  In case you are interested, read more about haml
[here](http://haml.info/) but beware that most of the "advanced"
features are not implemented by design.

This means the library lacking a number of haml features, which,
depending on your point of view is either a good thing or a catatrophe.

## G/HAML in 10 seconds.

Basically this:

    %html
      %head
      %body
        %h1 Hello World

translates to the corresponding HTML. You may indent using tabs or
spaces, but be consistant. The first indented line determines 
whether you are using tabs or spaces, and in the latter case how many spaces consitute one level of
indention.

Attributes are handled like this:

    %a(href='http://example.com')

Note the single quotes around the attribute values. No choice of single or
double. (I might switch to double quotes or makes this configurable, though)

`id` and `class` attributes can be abbreviated further with `#` and `.` 
respectively:

    %a.className#idName

And if they are only attached to `div`s for formatting, you can even leave 
out the tag:

    .bla

becomes:

    <div class='bla'></div>

In consequence, an id or class shortcut-attribute can't contain dot or
hash in it's name. In general, most errors are currently handled via
Garbage-in-garbage-out and escaping special characters is not supported. With
one exception:

### Go Templates

Gaml doesn't support variable interpolation. Mainly because I'm to lazy
(... too stupid ...). I figure if you need variable interpolation, you
should use the Gaml output as input to Go templates. Go templates make
judivious use of the dot (.) which, of course conflicts with the g/haml
class shortcut. Therefore, the exception to the "no-escaping" rule is:
everything in go template double braces (`{{ .go_template_stuff }}`) is
passed through and not considered to be a g/haml dot, hash or whatever.


## Usage

The simplest case: convert a string to html:
```go
	package main
	import ("gaml")
	
	func main() {
		line := `
	%html
	  %head
	  %body
	    %h1
	      Hello World`
		
		html, _ := gaml.GamlToHtml(line)
		println(html)
	}
```	
	
yields:

    $ go run main/main.go
    <html>
     <head>
      <body>
       <h1>
        Hello World
       </h1>
      </body>
     </head>
    </html>	  

Though typically, you will want to render directly to an `io.Writer`:

```go
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
```

this outputs:

    $ go run main/main2.go
    <html><head><body><h1>Hello World
    </h1></body></head></html>

Note the default Renderer does not indent the output and doesn't insert newlines 
between tags, though it does after text elements.

### HTTP

The package also includes an `http.Handler` implementation for integration with 
the `net.http` package. See the `NewGamlHandler` docs in `http.go`

## Additional Functionality (compared to Haml)

### Includes

gaml is able to include other fragments using `> fileToIncluder`
notation, e.g.

    %html
      %body
        %h1
          > childOfH1.gaml

or

    %html
      %body
        %h1
        > childOfBodySiblingOfH1.gaml

The examples above will insert the fragments named `childOfH1.gaml`, resp
`childOfBodySiblingOfH1.gaml` into the resulting html at the position 
suggested by their names. 

To use includes, the Parser needs to be assigned a Loader so it knows
how to retrieve the includes. If the Parser is created using the
FileSystemLoader, loading includes are handled by the same
FileSystemLoader by default. (This needs to be explained more clearly:)

## Most glaring HAML incompatibilities

### no ruby attribute syntax:

    %a {bla => "durp"}

#### Rationale

I don't see the point. I'm perfectly happy with "html" style attributes:

    %a(href='whatever')

This only makes sense in a ruby environment. Ruby code in those templates won't be
portable anyway.

## Single quoted attribute values onlu

This is legal:

    %a(href='http://example.com')

this is not:

    %a(href="http://example.com")

#### Rationale

Simplicity of implementation. Also, gaml main raison d'etre is to ease
typing, doublequotes require you pressing "shift" :)

### no variable interpolation

#### Rationale

* I don't need it, I only want to generate static html
* It's too hard for me
* Go comes with a great templating engine which can easily be combined with gaml
  see the `template.go` example in the example subdirectory.

### Fewer/no formatting options

#### Rationale

I'm lazy, currently I'm just implementing the stuff I need. This will change as
I need more features or people contribute stuff.


## LICENSE

MIT, see LICENSE file.

