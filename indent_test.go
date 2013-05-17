package gaml

import (
	"bytes"
	"testing"
)

func TestIndent(t *testing.T) {
	buf := bytes.NewBufferString("%p\n\t%p")
	parser := NewParser(buf)
	if err := parser.Parse(); err != nil {
		t.Error(err)
	}
	if parser.indent != 1 {
		t.Errorf("expected indent 1, got %d", parser.indent)
	}

	buf = bytes.NewBufferString("%p\n\t%p\n\t %p")
	parser = NewParser(buf)
	if err := parser.Parse(); err == nil {
		t.Error("expected err!")
	} else {
		if err.Error() != "cannot mix spaces with tabs line(3):\t %p" {
			t.Errorf("unexpected err! >%s<", err.Error())
		}
	}

	buf = bytes.NewBufferString("%p\n\t %p")
	parser = NewParser(buf)
	if err := parser.Parse(); err == nil {
		t.Error("expected err!")
	} else {
		if err.Error() != "initial indent > 1 line(2):\t %p" {
			t.Errorf("unexpected err! >%s<", err.Error())
		}
	}

	buf = bytes.NewBufferString("%p\n \t%p")
	parser = NewParser(buf)
	if err := parser.Parse(); err == nil {
		t.Error("expected err!")
	} else {
		if err.Error() != "cannot mix tabs with spaces line(2): \t%p" {
			t.Errorf("unexpected err! >%s<", err.Error())
		}
	}
	buf = bytes.NewBufferString("%p\n  %p\n   %p")
	parser = NewParser(buf)
	if err := parser.Parse(); err == nil {
		t.Error("expected err!")
	} else {
		if err.Error() != "incoherent number of space, not a multiple of intial indent! line(3):   %p" {
			t.Errorf("unexpected err! >%s<", err.Error())
		}
	}
	buf = bytes.NewBufferString("%p\n  %p\n    %p")
	parser = NewParser(buf)
	if err := parser.Parse(); err != nil {
		t.Errorf("unexpected err! >%s<", err.Error())
	}

	if parser.indent != 2 {
		t.Errorf("expected indent 2, got %d", parser.indent)
	}

	buf = bytes.NewBufferString("%p\n   %p\n      %p\n         %p")
	parser = NewParser(buf)
	if err := parser.Parse(); err != nil {
		t.Errorf("unexpected err! >%s<", err.Error())
	}

	if parser.indent != 3 {
		t.Errorf("expected indent 3, got %d", parser.indent)
	}
	buf = bytes.NewBufferString("%p\n   %p\n      %p\n         %p\n      %p")
	parser = NewParser(buf)
	if err := parser.Parse(); err != nil {
		t.Errorf("unexpected err! >%s<", err.Error())
	}

	if parser.indent != 2 {
		t.Errorf("expected indent 2, got %d", parser.indent)
	}
}
