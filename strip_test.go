package gaml

import (
  "testing"
)

func TestStrip(t *testing.T) {
  parser := new(Parser)
  parser.line = "  bla  "
  parser.stripLine()
  if parser.line != "  bla  " {
    t.Error("stripLine changed line")
  }
  if parser.strippedLine != "bla" {
    t.Errorf("expected: bla got: %s", parser.strippedLine)
  }

  parser.line = "  bla  // bla bla"
  parser.stripLine()
  if parser.line != "  bla  // bla bla" {
    t.Error("stripLine changed line")
  }
  if parser.strippedLine != "bla" {
    t.Errorf("expected: bla got: %s", parser.strippedLine)
  }
}
