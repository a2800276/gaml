package gaml

import (
	"bufio"
	"bytes"
	"io"
	"strings"
)
import (
	"fmt"
)

// %p
//   %h1

type Parser struct {
	scanner       *bufio.Scanner
	line          string    // content of current line sans line ending
	gamlline      *gamlline // line with no comment or surrounding ws
	lineNo        int       // current line number
	indent        int       // current indention level
	prevIndent    int       // previous indent
	indentType    iType     // using tabs or space, determined by first line, mixing is not allowed
	indentSpaces  int       // how many space == one indention level, determined by usage on first indented line
	rootNode      *node     // "blank" root node of document that everything is attached to
	currentNode   *node     // keeps track of the current position while parsing
	done          bool      // done parsing?
	err           error     // cache error which may have occured during parsing
	IncludeLoader Loader
}

type iType int // use tabs or space for indention
const (
	DETERMINE iType = iota
	TAB
	SPACE
)

func NewParser(reader io.Reader) (parser *Parser) {
	parser = new(Parser)
	parser.scanner = bufio.NewScanner(reader)
	parser.rootNode = newRoot()
	parser.currentNode = parser.rootNode
	return
}

func NewParserString(gaml string) (parser *Parser) {
	return NewParser(bytes.NewBufferString(gaml))
}

func (p *Parser) Parse() (rootNode *node, err error) {
	if p.done {
		return p.rootNode, nil
	} else {
		p.done = true
	}
	for p.scanner.Scan() {
		p.lineNo++
		p.line = p.scanner.Text()
		if err = p.handleLine(); err != nil {
			p.err = err
			return
		}
		//fmt.Printf("(%d): %s\n", p.lineNo, p.line)
	}
	if err = p.scanner.Err(); err != nil {
		p.err = err
		return
	}
	return p.rootNode, err
}

func (p *Parser) parseInclude(name string) (err error) {
	var includedNode *node
	if includedNode, err = p.IncludeLoader.Load(strings.TrimSpace(name)); err != nil {
		return
	}

	// currentNode will be an INC node representing the
	// its children are the results of executing the include parser.
	p.currentNode.children = includedNode.children

	return
}

func (p *Parser) handleLine() (err error) {

	p.stripLine()

	if !p.gamlline.Empty() {
		if err = p.handleIndent(); err != nil {
			return
		}
		p.setCurrentNode()
		if err = p.gamlline.processIntoCurrentNode(p); err != nil {
			return err
		}
	}
	return
}

func (p *Parser) setCurrentNode() error {
	// %p      <-1
	//   %p    <-2
	//   %p    <-3
	//     %p  <-2
	// %p      <-1
	//   %p    <-2
	//     %p  <-2
	//   %p    <-4
	switch {
	case p.indent == 0: // case #1
		p.currentNode = newNode(p.rootNode)
	case p.indent > p.prevIndent: // case #2
		if p.indent-p.prevIndent > 1 {
			return p.Err("indention level increase by more than one")
		}
		p.currentNode = newNode(p.currentNode)

	// case p.indent == p.prevIndent:   // case #3
	//  p.currentNode = newNode(p.currentNode.parent)

	case p.indent <= p.prevIndent: // case #3 & #4
		parent := p.currentNode.parent
		for up := p.prevIndent - p.indent; up != 0; up-- {
			parent = parent.parent
		}
		p.currentNode = newNode(parent)
	}
	return nil
}

// remove comments as well as leading and trailing ws
// from p.line and assign to p.strippedLine
func (p *Parser) stripLine() {
	stripped := p.line
	if commentStart := strings.Index(p.line, "//"); commentStart != -1 {
		stripped = p.line[:commentStart]
	}
	//p.strippedLine = (gamlline)(strings.TrimSpace(stripped))
	p.gamlline = GamlLineFromString(strings.TrimSpace(stripped))
}

func (p *Parser) handleIndent() (err error) {
	if p.line == "" || (p.line[0] != ' ' && p.line[0] != '\t') {
		p.indent = 0
		return
	}
	var ws bytes.Buffer
	// collect initial indent
	for _, r := range p.line {
		if r != ' ' && r != '\t' {
			break
		}
		ws.WriteRune(r)
	}

	wsString := ws.String()

	if p.indentType == DETERMINE {
		return p.initIndent(wsString)
	} else {
		return p._handleIndent(wsString)
	}
}

func (p *Parser) initIndent(ws string) error {
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

func (p *Parser) _handleIndent(ws string) error {
	p.prevIndent = p.indent
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
	case DETERMINE:
		panic("cannot happen!")
	}
	return nil
}

func (p *Parser) Err(msg string) error {
	return fmt.Errorf("%s line(%d):%s", msg, p.lineNo, p.line)
}
