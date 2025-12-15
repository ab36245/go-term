package test

import (
	"strings"
	"testing"

	"github.com/ab36245/go-ansi"
	"github.com/ab36245/go-term"
)

func TestScreenEmpty(t *testing.T) {
	screen := screenMake(10, 5)
	builder := strings.Builder{}
	a := screenInput(screen, builder.String())

	e := `
	|term.Screen (5 x 10) pos (0, 0)
	|  styles
	|     0 flags 0x0000 fg (none) bg (none)
	|  buffer
	|     0   [00]   [00]   [00]   [00]   [00]   [00]   [00]   [00]   [00]   [00]
	|     1   [00]   [00]   [00]   [00]   [00]   [00]   [00]   [00]   [00]   [00]
	|     2   [00]   [00]   [00]   [00]   [00]   [00]   [00]   [00]   [00]   [00]
	|     3   [00]   [00]   [00]   [00]   [00]   [00]   [00]   [00]   [00]   [00]
	|     4   [00]   [00]   [00]   [00]   [00]   [00]   [00]   [00]   [00]   [00]
	`

	screenCheck(t, a, e)
}

func TestScreenText(t *testing.T) {
	screen := screenMake(10, 5)
	builder := strings.Builder{}
	builder.WriteString("now is the ")
	builder.WriteString("winter")
	builder.WriteString(" of our discontent")
	a := screenInput(screen, builder.String())

	e := `
	|term.Screen (5 x 10) pos (3, 5)
	|  styles
	|     0 flags 0x0000 fg (none) bg (none)
	|  buffer
	|     0  n[00]  o[00]  w[00]   [00]  i[00]  s[00]   [00]  t[00]  h[00]  e[00]
	|     1   [00]  w[00]  i[00]  n[00]  t[00]  e[00]  r[00]   [00]  o[00]  f[00]
	|     2   [00]  o[00]  u[00]  r[00]   [00]  d[00]  i[00]  s[00]  c[00]  o[00]
	|     3  n[00]  t[00]  e[00]  n[00]  t[00]   [00]   [00]   [00]   [00]   [00]
	|     4   [00]   [00]   [00]   [00]   [00]   [00]   [00]   [00]   [00]   [00]
	`

	screenCheck(t, a, e)
}

func TestScreenSgr(t *testing.T) {
	screen := screenMake(10, 5)
	builder := strings.Builder{}
	builder.WriteString("now is the ")
	builder.WriteString("\x1b[42mwinter\x1b[0m")
	builder.WriteString(" of our discontent")
	a := screenInput(screen, builder.String())

	e := `
	|term.Screen (5 x 10) pos (3, 5)
	|  styles
	|     0 flags 0x0000 fg (none) bg (none)
	|     1 flags 0x0000 fg (none) bg 4 0x02
	|  buffer
	|     0  n[00]  o[00]  w[00]   [00]  i[00]  s[00]   [00]  t[00]  h[00]  e[00]
	|     1   [00]  w[01]  i[01]  n[01]  t[01]  e[01]  r[01]   [00]  o[00]  f[00]
	|     2   [00]  o[00]  u[00]  r[00]   [00]  d[00]  i[00]  s[00]  c[00]  o[00]
	|     3  n[00]  t[00]  e[00]  n[00]  t[00]   [00]   [00]   [00]   [00]   [00]
	|     4   [00]   [00]   [00]   [00]   [00]   [00]   [00]   [00]   [00]   [00]
	`

	screenCheck(t, a, e)
}

func TestScreenCursor(t *testing.T) {
	screen := screenMake(10, 5)
	builder := strings.Builder{}
	builder.WriteString("now is the ")
	builder.WriteString("winter")
	builder.WriteString(" of our discontent")
	builder.WriteString("\x1b[1;2He")
	builder.WriteString("\x1b[3;6Hmal")
	a := screenInput(screen, builder.String())

	e := `
	|term.Screen (5 x 10) pos (2, 8)
	|  styles
	|     0 flags 0x0000 fg (none) bg (none)
	|  buffer
	|     0  n[00]  e[00]  w[00]   [00]  i[00]  s[00]   [00]  t[00]  h[00]  e[00]
	|     1   [00]  w[00]  i[00]  n[00]  t[00]  e[00]  r[00]   [00]  o[00]  f[00]
	|     2   [00]  o[00]  u[00]  r[00]   [00]  m[00]  a[00]  l[00]  c[00]  o[00]
	|     3  n[00]  t[00]  e[00]  n[00]  t[00]   [00]   [00]   [00]   [00]   [00]
	|     4   [00]   [00]   [00]   [00]   [00]   [00]   [00]   [00]   [00]   [00]
	`

	screenCheck(t, a, e)
}

func TestScreenOther(t *testing.T) {
	run := func(name string, p string, e string) {
		t.Run(name, func(t *testing.T) {
			screen := screenMake(10, 5)
			builder := strings.Builder{}
			builder.WriteString("now is the ")
			builder.WriteString("winter")
			builder.WriteString(" of our discontent")
			builder.WriteString("\x1b[2;5H")
			builder.WriteString("\x1b[" + p + "J")
			a := screenInput(screen, builder.String())
			screenCheck(t, a, e)
		})
	}

	run("above", "1", `
		|term.Screen (5 x 10) pos (1, 4)
		|  styles
		|     0 flags 0x0000 fg (none) bg (none)
		|  buffer
		|     0  _[00]   [00]   [00]   [00]   [00]   [00]   [00]   [00]   [00]   [00]
		|     1   [00]   [00]   [00]   [00]   [00]  e[00]  r[00]   [00]  o[00]  f[00]
		|     2   [00]  o[00]  u[00]  r[00]   [00]  d[00]  i[00]  s[00]  c[00]  o[00]
		|     3  n[00]  t[00]  e[00]  n[00]  t[00]   [00]   [00]   [00]   [00]   [00]
		|     4   [00]   [00]   [00]   [00]   [00]   [00]   [00]   [00]   [00]   [00]
	`)

	run("below", "", `
		|term.Screen (5 x 10) pos (1, 4)
		|  styles
		|     0 flags 0x0000 fg (none) bg (none)
		|  buffer
		|     0  n[00]  o[00]  w[00]   [00]  i[00]  s[00]   [00]  t[00]  h[00]  e[00]
		|     1   [00]  w[00]  i[00]  n[00]   [00]   [00]   [00]   [00]   [00]   [00]
		|     2   [00]   [00]   [00]   [00]   [00]   [00]   [00]   [00]   [00]   [00]
		|     3   [00]   [00]   [00]   [00]   [00]   [00]   [00]   [00]   [00]   [00]
		|     4   [00]   [00]   [00]   [00]   [00]   [00]   [00]   [00]   [00]   [00]
	`)
}

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
