# About

Fairly haml-ish like html templating for GO.

Makes it easier to type up html by hand by avoiding having to type
superfluous angled bracked, closing tags, and some other redundancies.
Fans of Haml claim that it's not just a shortcut, but beautiful (like Haikus)
but they have a warped sense of aesthetics.

Basically this:

    %html
      %head
      %body
        %h1 Hello World

translates to the corresponding HTML.

Attributes are handled like this:

    %a(href="http://example.com")

`id` and `class` attributes can be abbreviated further with `#` and `.` 
respectively:

    %a.className#idName

And if they are only attached to `div`s for formatting, you can even leave 
out the tag:

    .bla

becomes

    <div class='bla'></div>

In case you are interested, read more about haml [here](http://haml.info/) but beware that
most of the "advanced" features are not implemented by design.

This means the library lacking a number of haml features, which, depending on your point of
view is either a good thing or a catatrophe:

## no ruby attribute syntax:

    %a {bla => "durp"}

### Rationale

I don't see the point. I'm perfectly happy with "html" style attributes:

    %a(href='whatever')

This only makes sense in a ruby environment. Ruby code in those templates won't be
portable anyway.

## no variable interpolation

### Rationale

* I don't need it, I only want to generate static html
* It's too hard for me
* Go comes with a great templating engine which can easily be combined with gaml

## Less formatting options

### Rationale

I'm lazy, currently I'm just implementing the stuff I need. This will change as
I need more features or people contribute stuff.


## LICENSE

MIT, see LICENSE file.

