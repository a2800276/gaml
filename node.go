package gaml


type node struct {
  parent * node
  children []*node
}

func newNode(parent * node)*node {
  n:= new(node)
  n.parent = parent
  if parent != nil {
    parent.children = append(parent.children, n)
  }
  return n
}

func (n * node) Render(indent int) {
  for i:=0;i!=indent;i++ {
    print("-")
  }
  println("+")
  for _,child := range(n.children) {
    child.Render(indent + 1)
  }
}
