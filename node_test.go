package gaml

import (
	"bytes"
	"testing"
)

func test_simple(t *testing.T, in string, expected string) {

	var bufout bytes.Buffer
	buf := bytes.NewBufferString(in)
	parser := NewParser(buf)
	if err := parser.Parse(); err != nil {
		t.Error(err)
	}
	for _, node := range parser.rootNodes {
		node.Render(&bufout)
	}
	if bufout.String() != expected {
		t.Errorf("expected: %s, got: %s", expected, bufout.String())
	}
}

func TestClass(t *testing.T) {
	test_simple(t, ".bla", "<div class='bla'>\n</div>\n")
	test_simple(t, ".bla.bla.bla", "<div class='bla bla bla'>\n</div>\n")
	test_simple(t, ".bing.bang.bum", "<div class='bing bang bum'>\n</div>\n")
}

func TestId(t *testing.T) {
	test_simple(t, "#bla", "<div id='bla'>\n</div>\n")
	test_simple(t, ".bing#bang.bum", "<div class='bing bum' id='bang'>\n</div>\n")
}

func TestAttributes(t *testing.T) {
	test_simple(t, "%a#bla.blub(ding='dong' ping='pong pung') hello world!",
		"<a id='bla' class='blub' ding='dong' ping='pong pung'>\n</a>\n")
}

func TestDoctype(t *testing.T) {
	test_simple(t, "!!!\n%html\n %body\n  %h1 Hello World!",
		"<!DOCTYPE html>\n<html>\n <body>\n  <h1>\n    Hello World!\n  </h1>\n </body>\n</html>\n")
}

const blank_line = `
%html
	%head

	%body
`
const blank_expected = `<html>
 <head>
 </head>
 <body>
 </body>
</html>
`

func TestBlank(t *testing.T) {
	test_simple(t, blank_line, blank_expected)
}

//func TestNode(t * testing.T) {
//  buf := bytes.NewBufferString("%p\n %p\n  %p\n   %p\n   %p\n  %p\n%p\n %p")
//  parser := NewParser(buf)
//  if err := parser.Parse(); err != nil {
//    t.Error(err)
//  }
//  for _, node := range(parser.rootNodes) {
//    t.Log("%p\n %p\n  %p\n   %p\n   %p\n  %p\n%p\n %p")
//    node.render_debug(0)
//  }
//}

//func TestVaried2(t * testing.T) {
//  var bufout bytes.Buffer
//  buf := bytes.NewBufferString("%html\n %head\n  %title bla\n %body\n  %h1\n   heading!\n  %p\n  %p\n  %h2\n   %p")
//  parser := NewParser(buf)
//  if err := parser.Parse(); err != nil {
//    t.Error(err)
//  }
//  for _, node := range(parser.rootNodes) {
//    node.Render(&bufout)
//  }
//  println(bufout.String())
//}
