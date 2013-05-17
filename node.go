package gaml

import (
  "io"
)

type node struct {
  parent * node
  children []*node
  name string
}

func newNode(parent * node)*node {
  n:= new(node)
  n.parent = parent
  if parent != nil {
    parent.children = append(parent.children, n)
  }
  return n
}


func (n * node) Render(writer io.Writer) {
  n.render(writer, 0)
}
func (n * node) render(w io.Writer, indent int) {
  indentfunc := func () {
    for i:=0; i!=indent; i++ {
      io.WriteString(w, " ")
    }
  }
  indentfunc()
  io.WriteString(w, "<")
  io.WriteString(w, n.name)
  io.WriteString(w, ">\n")

  for _,child := range(n.children) {
    child.render(w, indent+1)
  }

  indentfunc()
  io.WriteString(w, "</")
  io.WriteString(w, n.name)
  io.WriteString(w, ">\n")
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
