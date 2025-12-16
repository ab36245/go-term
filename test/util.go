package test

import (
	"regexp"
	"strings"
	"testing"
	"unicode"

	"github.com/ab36245/go-ansi"
	"github.com/ab36245/go-term"
)

func screenCheck(t *testing.T, a, e string) {
	indent := "  "
	e = addBorder(stripBorder(e), indent)
	a = addBorder(a, indent)
	if a != e {
		s := "\n"
		s += "expected:\n"
		s += e
		s += "\n"
		s += "actual:\n"
		s += a
		t.Fatalf("%s", s)
	}
}

func screenInput(screen *term.Screen, input string) string {
	reader := strings.NewReader(input)
	parser := ansi.NewParser(reader, nil)
LOOP:
	for {
		seq := parser.Next()
		switch seq.(type) {
		case ansi.EOF, ansi.Err:
			break LOOP
		default:
			screen.ApplySeq(seq)
		}
	}
	return screen.Dump()
}

func screenMake(width, height int) *term.Screen {
	screen := term.New(width, height)
	// TODO
	return screen
}

const borderMark = "|"

var borderRegexp = regexp.MustCompile(`^\s*\|(.*)$`)

func addBorder(in, indent string) string {
	in = strings.TrimRightFunc(in, unicode.IsSpace)
	out := []string{}
	for line := range strings.SplitSeq(in, "\n") {
		out = append(out, indent+borderMark+line)
	}
	return strings.Join(out, "\n")
}

func stripBorder(in string) string {
	out := []string{}
	for s := range strings.SplitSeq(in, "\n") {
		matches := borderRegexp.FindStringSubmatch(s)
		if len(matches) == 2 {
			out = append(out, matches[1])
		}
	}
	return strings.Join(out, "\n")
}
