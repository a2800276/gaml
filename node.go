package gaml

import (
  "io"
  "strings"
)

type node struct {
  parent * node
  children []*node
  name string
  attributes map[string][]string
  text string
}

func newNode(parent * node)*node {
  n:= new(node)
  n.parent = parent
  if parent != nil {
    parent.children = append(parent.children, n)
  }
  return n
}

func (n * node) AddAttribute (name string, value string) {
  if n.attributes == nil {
    n.attributes = make(map[string][]string)
  }
  n.attributes[name] = append(n.attributes[name], value)
}


func (n * node) Render(writer io.Writer) {
  n.render(writer, 0)
}
func (n * node) render(w io.Writer, indent int) {
  if n.name != "" {
    n.renderTag(w, indent)
  } else {
    n.renderText(w, indent)
  }
}

func (n * node) render_debug(indent int) {
  for i:=0;i!=indent;i++ {
    print("-")
  }
  println("+")
  for _,child := range(n.children) {
    child.render_debug(indent + 1)
  }
}
func (n * node) renderTag(w io.Writer, indent int) {
  indentfunc := func () {
    for i:=0; i!=indent; i++ {
      io.WriteString(w, " ")
    }
  }
  indentfunc()
  io.WriteString(w, "<")
  io.WriteString(w, n.name)
  n.renderAttributes(w)
  io.WriteString(w, ">\n")

  for _,child := range(n.children) {
    child.render(w, indent+1)
  }

  indentfunc()
  io.WriteString(w, "</")
  io.WriteString(w, n.name)
  io.WriteString(w, ">\n")
}
func (n * node) renderAttributes(w io.Writer) {
  for key, values := range(n.attributes) {
    io.WriteString(w, " ")
    io.WriteString(w, key)
    io.WriteString(w, "='")
    io.WriteString(w, strings.Join(values, " "))
    io.WriteString(w, "'")
  }
}
func (n * node) renderText(w io.Writer, indent int) {
  for i:=0; i!=indent; i++ {
    io.WriteString(w, " ")
  }
  io.WriteString(w, n.text)
  io.WriteString(w, "\n")

}
