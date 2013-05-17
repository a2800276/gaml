package gaml

import (
  "bytes"
)

type gamlline string

type gstate int

const (
  INITIAL gstate = iota
  TAG_NAME
  CLASS
  ID
  INCLUDE
  TEXT
  TEXT_OR_ATTRIBUTES
  ATTRIBUTES
)

func (g gamlline) empty()bool {
  return string(g) == ""
}
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


func (g gamlline) sm_curr_node(p* Parser)(err error) {
  line := string(g)
  node := p.currentNode
  var value bytes.Buffer

  fillInName := func() {
    node.name = value.String()
    value.Reset()
  }

  addClass := func() {}
  addId := func() {}

  state := INITIAL

  for _, r := range(line) {
    switch state {
      case INITIAL:
        state = initial(r)
      case TAG_NAME:
        state = tag(r, &value, fillInName)
      case CLASS:
        state = tag(r, &value, addClass)
      case ID:
        state = tag(r, &value, addId)
    }
  }

  switch state {
    case INITIAL:
      panic("cannot happen")
    case TAG_NAME:
      node.name = value.String()
  }
  return
}


func tag(r rune, buf * bytes.Buffer, fillInValue func () )gstate {
    switch r {
    case '.':
      fillInValue()
      return CLASS
    case '#':
      fillInValue()
      return ID
    case ' ':
      fillInValue()
      buf.WriteRune(r)
      return TEXT_OR_ATTRIBUTES
    case '(':
      fillInValue()
      return ATTRIBUTES
    default:
      buf.WriteRune(r)
      return TAG_NAME
  }
}

func initial(r rune)gstate {
  switch r {
    case '%':
      return TAG_NAME
    case '.':
      return CLASS
    case '#':
      return ID
    case '>':
      return INCLUDE
    default:
      return TEXT
  }
}
