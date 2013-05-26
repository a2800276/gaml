# About

Fairly haml-ish like html templating for GO.

Makes it easier to type up html by hand by avoiding having to type
superfluous angled bracked, closing tags, and some other redundancies.
Fans of Haml claim that it's not just a shortcut, but beautiful (like
Haikus) but they have a warped sense of aesthetics.

I'm going to assume you have a vague notion of what haml is so I'll keep
it short.  In case you are interested, read more about haml
[here](http://haml.info/) but beware that most of the "advanced"
features are not implemented by design.

This means the library lacking a number of haml features, which,
depending on your point of view is either a good thing or a catatrophe:

## G/HAML in 10 seconds.

Basically this:

    %html
      %head
      %body
        %h1 Hello World

translates to the corresponding HTML. You may indent using tabs or
spaces, but be consistant. The first indented line determines a.)
whether you are using tabs, resp. how many spaces consitute one level of
indention.

Attributes are handled like this:

    %a(href="http://example.com")

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
Garbage-in-garbage-out and there is no escaping special characters. With
one exception:

### Go Templates

Gaml doesn't support variable interpolation. Mainly because I'm to lazy
(... too stupid ...). I figure if you need variable interpolation, you
can use the Gaml output as input to Go templates. Go templates make
judivious use of the dot (.) which, of course conflicts with the g/haml
class shortcut. Therefor, the exception to the "no-escaping" rule is:
everything in go template double braces (`{{ .go_template_stuff }}`) is
passed through and not considered to be a g/haml dot, hash or whatever.


## Most glaring HAML incompatibilities

### no ruby attribute syntax:

    %a {bla => "durp"}

### Rationale

I don't see the point. I'm perfectly happy with "html" style attributes:

    %a(href='whatever')

This only makes sense in a ruby environment. Ruby code in those templates won't be
portable anyway.

### no variable interpolation

### Rationale

* I don't need it, I only want to generate static html
* It's too hard for me
* Go comes with a great templating engine which can easily be combined with gaml
  see the `template.go` example in the example subdirectory.

### Fewer/no formatting options

### Rationale

I'm lazy, currently I'm just implementing the stuff I need. This will change as
I need more features or people contribute stuff.


## LICENSE

MIT, see LICENSE file.
