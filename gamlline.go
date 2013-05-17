package gaml

import (
  "bytes"
  "fmt"
  "strings"
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
  TEXT_NEW // need to differentiate between a "pure" line of text and text that comes after a tag.
  ATTRIBUTES
  ATTRIBUTES_NAME
  ATTRIBUTES_AFTER_NAME
  ATTRIBUTES_VALUES
  ERR
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

  if 0 == strings.Index(line, "!!!") {
    node.text = "!!!" // not nice!
    node.name = "!!!" // use text == name == "!!!" to signal doctype.
  }


  var value bytes.Buffer

  fillInName := func() {
    node.name = value.String()
    value.Reset()
  }

  addClass := func() {
    if node.name == "" {
      node.name = "div"
    }
    node.AddAttribute("class", value.String())
    value.Reset()
  }

  addId := func() {
    if node.name == "" {
      node.name = "div"
    }
    node.AddAttribute("id", value.String())
    value.Reset()
  }

  textNew := func () {
    _node := newNode(node)
    node = _node
  }


  state := INITIAL
  var name string // remember name of name = attribute pairs

  for _, r := range(line) {
    switch state {
      case INITIAL:
        state = initial(r, &value)
      case TAG_NAME:
        state = tag(r, &value, fillInName, TAG_NAME)
      case CLASS:
        state = tag(r, &value, addClass, CLASS)
      case ID:
        state = tag(r, &value, addId, ID)
      case INCLUDE:
        // ignore for now ... ? 
        node.text = "<!-- include not handled -->"
        return nil
      case TEXT:
        value.WriteRune(r)
      case TEXT_OR_ATTRIBUTES:
        state = textOrAttribute(r, &value)
      case TEXT_NEW:
        textNew() 
        value.WriteRune(r)
        state = TEXT
      case ATTRIBUTES:
        state = attributes(r, &value)
      case ATTRIBUTES_NAME:
        state, name = attributes_name(r, &value)
      case ATTRIBUTES_AFTER_NAME:
        if state, err = attributes_after_name(r, &value); err != nil {
          return err
        }
      case ATTRIBUTES_VALUES:
        state = attributes_values(r, &value, func(){ node.AddAttribute(name, value.String())} )
    }
  }


  // stow away the value we have been collecting once we've
  // passed through the entire string.
  switch state {
    case INITIAL:
      panic("cannot happen")
    case TAG_NAME:
      fillInName()
    case CLASS:
      addClass()
    case ID:
      addId()
    case INCLUDE:
      // TODO
    case TEXT:
      node.text = value.String()
    case TEXT_OR_ATTRIBUTES, TEXT_NEW:
      textNew()
      node.text = value.String()
    case ATTRIBUTES, ATTRIBUTES_NAME, ATTRIBUTES_AFTER_NAME, ATTRIBUTES_VALUES:
      return p.Err("implausible state!")


  }
  return
}

func attributes_values (r rune, buf * bytes.Buffer, fillInValue func()) gstate {
  switch r {
    case '"', '\'':
      fillInValue()
      buf.Reset()
      return ATTRIBUTES
    default:
      buf.WriteRune(r)
      return ATTRIBUTES_VALUES
  }
}

func attributes_after_name (r rune, buf * bytes.Buffer)(gstate, error) {
  // this one is stupid.
  switch r {
    case ' ', '=':                             // <-- allows a ==    == = 'bla'
      return ATTRIBUTES_AFTER_NAME, nil
    case '\'', '"':                            // <-- allows a = 'Bla"
      return ATTRIBUTES_VALUES, nil
    default:
      return ERR, fmt.Errorf("unquoted attribute values")
  }
}

func attributes_name (r rune, buf * bytes.Buffer)(gstate, string) {
  switch r {
    case ' ', '=':
      name := buf.String() 
      buf.Reset()
      return ATTRIBUTES_AFTER_NAME, name
    default:
      buf.WriteRune(r)
      return ATTRIBUTES_NAME, ""
  }
}
func attributes(r rune, buf * bytes.Buffer) gstate {
  switch r {
    case ' ':
      return ATTRIBUTES
    case ')':
      return TEXT
    default:
      buf.WriteRune(r)
      return ATTRIBUTES_NAME
  }
}

func textOrAttribute(r rune, buf * bytes.Buffer) gstate {
  switch r {
    case ' ':
      buf.WriteRune(r)
      return TEXT_OR_ATTRIBUTES
    case '(':
      buf.Reset()
      return ATTRIBUTES
    default:  // default in all of these state functions is not correct, catch all sort of unicode crap people may throw at us ...
      buf.WriteRune(r)
      return TEXT_NEW
  }
}

func tag(r rune, buf * bytes.Buffer, fillInValue func (), state gstate )gstate {
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
      return state
  }
}

func initial(r rune, buf * bytes.Buffer)gstate {
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
      buf.WriteRune(r)
      return TEXT
  }
}
