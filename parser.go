package gaml

import (
  "bufio"
  "io"
  "bytes"
  "strings"
)
import (
  "fmt"
)

// %p
//   %h1


type Parser struct {
  scanner *bufio.Scanner
  line string
  line_no int
  indent int
  indentType iType
  indentSpaces int
}

type iType int
const (
  UNKNOWN iType = iota
  TAB
  SPACE
)

func NewParser(reader io.Reader)(parser * Parser) {
  parser = new(Parser)
  parser.scanner = bufio.NewScanner(reader)
  return
}

func (p * Parser) Parse()(err error) {
  for ;p.scanner.Scan(); {
    p.line_no++
    p.line = p.scanner.Text()
    if err = p.handleLine(); err != nil {
      return
    }
    fmt.Printf("(%d): %s\n", p.line_no, p.line)
  }
  if err = p.scanner.Err(); err != nil {
    return
  }
  return
}

func (p * Parser) handleLine()(err error) {
  if err = p.handleIndent(); err != nil {
    return
  }
  return 
}

func (p * Parser) handleIndent() (err error){
  if p.line[0] != ' ' && p.line[0] != '\t' {
    p.indent = 0
    return
  }
  var ws bytes.Buffer
  // collect initial indent
  for _, r := range(p.line) {
    if r != ' ' && r != '\t' {
      break
    }
    ws.WriteRune(r)
  }

  wsString := ws.String()

  if p.indentType == UNKNOWN {
    return p.initIndent(wsString)
  } else {
    return p._handleIndent(wsString)
  }
}

func (p * Parser) initIndent(ws string)error {
  switch {
    case ws[0] == ' ':
      p.indentType = SPACE
      p.indentSpaces = len(ws)
    case ws[0] == '\t':
      if len(ws) > 1 {
        return p.Err("initial indent > 1") 
      }
      p.indentType = TAB
  }
  return p._handleIndent(ws)
}

func (p * Parser) _handleIndent(ws string)error {
  switch p.indentType {
    case SPACE:
      if strings.IndexRune(ws, '\t') != -1 {
        return p.Err("cannot mix tabs with spaces")
      }
      if (len(ws) % p.indentSpaces) != 0 {
        return p.Err("incoherent number of space, not a multiple of intial indent!")
      }
      p.indent = len(ws) / p.indentSpaces
    case TAB:
      if strings.IndexRune(ws, ' ') != -1 {
        return p.Err("cannot mix spaces with tabs")
      }
      p.indent = len(ws)
    case UNKNOWN:
      panic("cannot happen!")
  }
  return nil
}

func (p * Parser) Err(msg string) error {
  return fmt.Errorf("%s line(%d):%s", msg, p.line_no, p.line)
}


