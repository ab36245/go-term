package test

import (
	"fmt"
	"regexp"
	"strings"
	"testing"
	"unicode"
)

func report(t *testing.T, name string, a, e any) {
	value := func(v any) string {
		switch v := v.(type) {
		case string:
			return fmt.Sprintf("%q", v)
		default:
			return fmt.Sprintf("%v", v)
		}
	}

	s := "\n"
	s += "expected:"
	if name != "" {
		s += " " + name
	}
	s += " " + value(e)

	s += "\n"
	s += "actual  :"
	if name != "" {
		s += " " + name
	}
	s += " " + value(a)

	t.Fatal(s)
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
