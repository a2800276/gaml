package gaml

import (
	"strings"
	"testing"
)

func test_compare(t *testing.T, is string, expected string) {
	if is != expected {
		t.Errorf("expected: %s, got: %s", expected, is)
	}
}

func test_simple(t *testing.T, in string, expected string) {
	if out, err := GamlToHtml(in); err != nil {
		t.Error(err)
	} else {
		test_compare(t, out, expected)
	}
}

func test_err(t *testing.T, in string, expected_err string) {
	if _, err := GamlToHtml(in); err != nil {
		test_compare(t, err.Error(), expected_err)
	} else {
		t.Error("expected fail")
	}
}

func TestClass(t *testing.T) {
	test_simple(t, ".bla", "<div class='bla'>\n</div>\n")
	test_simple(t, ".bla.bla.bla", "<div class='bla bla bla'>\n</div>\n")
	test_simple(t, ".bing.bang.bum", "<div class='bing bang bum'>\n</div>\n")
}

func TestId(t *testing.T) {
	test_simple(t, "#bla", "<div id='bla'>\n</div>\n")
	test_simple(t, ".bing#bang.bum", "<div class='bing bum' id='bang'>\n</div>\n")
	test_simple(t, "%img#test_id.test_class", "<img id='test_id' class='test_class'>\n")
	test_simple(t, "%img#test_image.img_class#id2", "<img id='test_image id2' class='img_class'>\n")
	// not sure how to go about non-sensical input.
	// as of now, I'll go with garbage-in-garbage-out
	// but it might be friendly to at least warn about
	// quotes in tags, ids and class...
	test_simple(t, "%img#{src='test.png'}", "<img id='{src='test' class='png'}'>\n")
}

func TestAttributeNoValue(t *testing.T) {
	test_simple(t, "%html(ng-app)", "<html ng-app>\n</html>\n")
	test_simple(t, "%html(ng-app ding='dong')", "<html ng-app ding='dong'>\n</html>\n")
	test_simple(t, "%html(ding='dong' ng-app)", "<html ding='dong' ng-app>\n</html>\n")
}
func TestAttributes(t *testing.T) {
	test_simple(t, "%a#bla.blub(ding='dong' ping='pong pung') hello world!",
		"<a id='bla' class='blub' ding='dong' ping='pong pung'>\n  hello world!\n</a>\n")
}

func TestDoctype(t *testing.T) {
	test_simple(t, "!!!\n%html\n %body\n  %h1 Hello World!",
		"<!DOCTYPE html>\n<html>\n <body>\n  <h1>\n    Hello World!\n  </h1>\n </body>\n</html>\n")
}

func TestSpecialInQuotes(t *testing.T) {
	test_simple(t, "%bla(class='bla.bla')", "<bla class='bla.bla'>\n</bla>\n")
	test_simple(t, "%bla(class='bla#bla')", "<bla class='bla#bla'>\n</bla>\n")
	test_simple(t, "%li#{{.First}}_{{.Last}} {{.First}}", "<li id='{{.First}}_{{.Last}}'>\n  {{.First}}\n</li>\n")
	test_simple(t, "%li#{.First}} {{.First}}", "<li id='{' class='First}}'>\n  {{.First}}\n</li>\n")
	test_simple(t, "%li#{{.First} }} {{.First}}", "<li id='{{.First} }}'>\n  {{.First}}\n</li>\n")
	test_simple(t, "%{{%%.#}}", "<{{%%.#}}>\n</{{%%.#}}>\n")
	test_simple(t, "%a.{{.class_name}}", "<a class='{{.class_name}}'>\n</a>\n")
}

func TestSpecialInQuotesFail(t *testing.T) {
	if str, err := GamlToHtml("%a.{{.something or the other"); err == nil {
		t.Errorf("expected an error, did not get! instead: %s", str)
	} else if strings.Index(err.Error(), "implausible state!") != 0 {
		t.Errorf("unexpected error: %s", err.Error())
	}
}

func TestTextOnSameLine(t *testing.T) {
	test_simple(t, "%bla Bla", "<bla>\n  Bla\n</bla>\n")
	test_simple(t, "%bla(something) Bla", "<bla something>\n  Bla\n</bla>\n")
	test_simple(t, "%bla(something='something else') Bla", "<bla something='something else'>\n  Bla\n</bla>\n")
}

func TestBooleanAttribute(t *testing.T) {
	test_simple(t, "%bla(bool)", "<bla bool>\n</bla>\n")
	test_simple(t, "%bla(bool) Just Testin'", "<bla bool>\n  Just Testin'\n</bla>\n")
}
func TestVoidElement(t *testing.T) {
	test_simple(t, "%br", "<br>\n")
	test_simple(t, "%img#test_image.img_class(src='test.png')", "<img id='test_image' class='img_class' src='test.png'>\n")
}

func TestQuote(t *testing.T) {
	test_simple(t, "%body(onload='alert(\"hello1\")')", "<body onload='alert(\"hello1\")'>\n</body>\n")
	test_simple(t, "%body(onload='alert(\\'hello2\\')')", "<body onload='alert('hello2')'>\n</body>\n")
	test_simple(t, "%tag(attr='some\\tthing')", "<tag attr='some\tthing'>\n</tag>\n")
	test_simple(t, "%tag(attr='some\\\\thing')", "<tag attr='some\\thing'>\n</tag>\n")
	// support arbitrary escapes
	test_simple(t, "%tag(attr='some\\Zthing')", "<tag attr='someZthing'>\n</tag>\n")
	test_err(t, "%body(onload=\"alert('hello')\")", "attribute values must be in single quote line(1):%body(onload=\"alert('hello')\")")
}

const blank_line = `
%html
	%head

	%body
`
const blank_expected = `<html>
 <head>
 </head>
 <body>
 </body>
</html>
`

func TestBlank(t *testing.T) {
	test_simple(t, blank_line, blank_expected)
}
