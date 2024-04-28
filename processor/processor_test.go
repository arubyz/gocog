package processor

import (
	"bufio"
	"bytes"
	"io"
	"testing"
)

type CPTData struct {
	input  string
	output string
	prefix string
	first  bool
	err    error
}

func TestCogPlainText(t *testing.T) {
	p := New("foo", &Options{GenStart: "@GENERATE@"})

	tests := []CPTData{
		{"", "", "", true, NoCogCode},
		{"a\nb\nc", "", "", true, NoCogCode},
		{"a\nb\n@GENERATE@", "", "", true, io.ErrUnexpectedEOF},
		{"a\nb\n@GENERATE@\n", "a\nb\n@GENERATE@\n", "", true, nil},
		{"a\nb\n@GENERATE@  stuff\n and more stuff\n", "a\nb\n@GENERATE@  stuff\n", "", true, nil},
		{"a\nb\n// @GENERATE@\n", "a\nb\n// @GENERATE@\n", "// ", true, nil},
	}

	for i, test := range tests {
		_ = i
		in := bytes.NewBufferString(test.input)
		out := &bytes.Buffer{}

		r := bufio.NewReader(in)
		prefix, _, err := p.cogPlainText(r, out, test.first)

		if prefix != test.prefix {
			t.Errorf("CogPlainText Test %d: Expected prefix: '%s', Got prefix: '%s'", i, test.prefix, prefix)
		}

		if err != test.err {
			t.Errorf("CogPlainText Test %d: Expected error: '%v', Got error: '%v'", i, test.err, err)
		}

		output := out.String()
		if output != test.output {
			t.Errorf("CogPlainText Test %d: Expected output:\n'%s'\nGot output:\n'%s'", i, test.output, output)
		}
	}
}

type CTEData struct {
	input  string
	output string
	useEOF bool
	err    error
}

func TestCogToEnd(t *testing.T) {
	tests := []CTEData{
		{"", "", false, io.ErrUnexpectedEOF},
		{"", "", true, io.EOF},
		{"1\n2\n@END@", "@END@", false, io.EOF},
		{"1\n2\n@END@\n", "@END@\n", false, nil},
		{"1\n2", "", true, io.EOF},
		{"1\n2", "", false, io.ErrUnexpectedEOF},
		{"1\n2\n// @END@\n", "// @END@\n", false, nil},
	}

	opts := &Options{
		GenStart: "@GENERATE@",
		OutEnd:   "@END@",
	}
	p := New("foo", opts)

	for i, test := range tests {
		opts.UseEOF = test.useEOF

		in := bytes.NewBufferString(test.input)
		out := &bytes.Buffer{}

		r := bufio.NewReader(in)
		err := p.cogToEnd(r, out, "")

		if err != test.err {
			t.Errorf("CogToEnd Test %d: Expected error %v, got %v", i, test.err, err)
		}

		output := out.String()
		if output != test.output {
			t.Errorf("CogToEnd Test %d: Expected output:\n'%s'\nGot output:\n'%s'", i, test.output, output)
		}

	}

}
