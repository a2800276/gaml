package gaml

import (
  "bytes"
)

type gamlline string

func (g * gamlline) fillCurrNode(p* Parser)(err error){
  line := (string)(*g)
  node := p.currentNode
  var name bytes.Buffer
  switch line[0] {
    case '%':
      for _,r := range(line[1:]) {
        name.WriteRune(r)
      }
    default:
      return p.Err("unexpected char")
  }
  node.name = name.String()
  return nil
}
