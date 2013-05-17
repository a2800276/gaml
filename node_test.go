package gaml

import (
  "testing"
  "bytes"
)

func TestNode(t * testing.T) {
  buf := bytes.NewBufferString("%p\n %p\n  %p\n   %p\n   %p\n  %p\n%p\n %p")
  parser := NewParser(buf)
  if err := parser.Parse(); err != nil {
    t.Error(err)
  }
  for _, node := range(parser.rootNodes) {
    t.Log("%p\n %p\n  %p\n   %p\n   %p\n  %p\n%p\n %p")
    node.render_debug(0)
  }
}

func TestRender(t * testing.T) {
  var bufout bytes.Buffer
  buf := bytes.NewBufferString("%p\n %p\n  %p\n   %p\n   %p\n  %p\n%p\n %p")
  parser := NewParser(buf)
  if err := parser.Parse(); err != nil {
    t.Error(err)
  }
  for _, node := range(parser.rootNodes) {
    node.Render(&bufout)
  }
  println(bufout.String())
}

func TestVaried(t * testing.T) {
  var bufout bytes.Buffer
  buf := bytes.NewBufferString("%html\n %head\n  %title bla\n %body\n  %h1\n   %p\n   %p\n  %h2\n   %p")
  parser := NewParser(buf)
  if err := parser.Parse(); err != nil {
    t.Error(err)
  }
  for _, node := range(parser.rootNodes) {
    node.Render(&bufout)
  }
  println(bufout.String())
}
func TestVaried2(t * testing.T) {
  var bufout bytes.Buffer
  buf := bytes.NewBufferString("%html\n %head\n  %title bla\n %body\n  %h1\n   heading!\n  %p\n  %p\n  %h2\n   %p")
  parser := NewParser(buf)
  if err := parser.Parse(); err != nil {
    t.Error(err)
  }
  for _, node := range(parser.rootNodes) {
    node.Render(&bufout)
  }
  println(bufout.String())
}

func TestClass(t * testing.T) {
  var bufout bytes.Buffer
  buf := bytes.NewBufferString(".bla")
  parser := NewParser(buf)
  if err := parser.Parse(); err != nil {
    t.Error(err)
  }
  for _, node := range(parser.rootNodes) {
    node.Render(&bufout)
  }
  expected := "<div class='bla'>\n</div>\n"
  if bufout.String() != expected  {
    t.Errorf("expected: %s, got: %s", expected, bufout.String())
  }

  bufout.Reset()
  buf = bytes.NewBufferString(".bla.bla.bla")
  parser = NewParser(buf)
  if err := parser.Parse(); err != nil {
    t.Error(err)
  }
  for _, node := range(parser.rootNodes) {
    node.Render(&bufout)
  }
  expected = "<div class='bla bla bla'>\n</div>\n"
  if bufout.String() != expected  {
    t.Errorf("expected: %s, got: %s", expected, bufout.String())
  }

  bufout.Reset()
  buf = bytes.NewBufferString(".bing.bang.bum")
  parser = NewParser(buf)
  if err := parser.Parse(); err != nil {
    t.Error(err)
  }
  for _, node := range(parser.rootNodes) {
    node.Render(&bufout)
  }
  expected = "<div class='bing bang bum'>\n</div>\n"
  if bufout.String() != expected  {
    t.Errorf("expected: %s, got: %s", expected, bufout.String())
  }
}

func TestId(t * testing.T) {
  var bufout bytes.Buffer
  buf := bytes.NewBufferString("#bla")
  parser := NewParser(buf)
  if err := parser.Parse(); err != nil {
    t.Error(err)
  }
  for _, node := range(parser.rootNodes) {
    node.Render(&bufout)
  }
  expected := "<div id='bla'>\n</div>\n"
  if bufout.String() != expected  {
    t.Errorf("expected: %s, got: %s", expected, bufout.String())
  }
}
