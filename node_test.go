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
  buf := bytes.NewBufferString("%html\n %head\n %body\n  %h1\n   %p\n   %p\n  %h2\n   %p")
  parser := NewParser(buf)
  if err := parser.Parse(); err != nil {
    t.Error(err)
  }
  for _, node := range(parser.rootNodes) {
    node.Render(&bufout)
  }
  println(bufout.String())
}
