package test

import (
	"strings"
	"testing"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func TestEraseDisplay(t *testing.T) {
	log.Logger = log.Logger.Level(zerolog.WarnLevel)

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

	run("above", "1", `
		|term.Screen (5 x 10) pos (1, 4)
		|  styles
		|     0 flags 0x0000 fg (none) bg (none)
		|  buffer
		|     0   [00]   [00]   [00]   [00]   [00]   [00]   [00]   [00]   [00]   [00]
		|     1   [00]   [00]   [00]   [00]   [00]  e[00]  r[00]   [00]  o[00]  f[00]
		|     2   [00]  o[00]  u[00]  r[00]   [00]  d[00]  i[00]  s[00]  c[00]  o[00]
		|     3  n[00]  t[00]  e[00]  n[00]  t[00]   [00]   [00]   [00]   [00]   [00]
		|     4   [00]   [00]   [00]   [00]   [00]   [00]   [00]   [00]   [00]   [00]
	`)

	run("all", "2", `
		|term.Screen (5 x 10) pos (1, 4)
		|  styles
		|     0 flags 0x0000 fg (none) bg (none)
		|  buffer
		|     0   [00]   [00]   [00]   [00]   [00]   [00]   [00]   [00]   [00]   [00]
		|     1   [00]   [00]   [00]   [00]   [00]   [00]   [00]   [00]   [00]   [00]
		|     2   [00]   [00]   [00]   [00]   [00]   [00]   [00]   [00]   [00]   [00]
		|     3   [00]   [00]   [00]   [00]   [00]   [00]   [00]   [00]   [00]   [00]
		|     4   [00]   [00]   [00]   [00]   [00]   [00]   [00]   [00]   [00]   [00]
	`)
}

func TestEraseLine(t *testing.T) {
	log.Logger = log.Logger.Level(zerolog.WarnLevel)

	run := func(name string, p string, e string) {
		t.Run(name, func(t *testing.T) {
			screen := screenMake(10, 5)
			builder := strings.Builder{}
			builder.WriteString("now is the ")
			builder.WriteString("winter")
			builder.WriteString(" of our discontent")
			builder.WriteString("\x1b[2;5H")
			builder.WriteString("\x1b[" + p + "K")
			a := screenInput(screen, builder.String())
			screenCheck(t, a, e)
		})
	}

	run("right", "", `
		|term.Screen (5 x 10) pos (1, 4)
		|  styles
		|     0 flags 0x0000 fg (none) bg (none)
		|  buffer
		|     0  n[00]  o[00]  w[00]   [00]  i[00]  s[00]   [00]  t[00]  h[00]  e[00]
		|     1   [00]  w[00]  i[00]  n[00]   [00]   [00]   [00]   [00]   [00]   [00]
		|     2   [00]  o[00]  u[00]  r[00]   [00]  d[00]  i[00]  s[00]  c[00]  o[00]
		|     3  n[00]  t[00]  e[00]  n[00]  t[00]   [00]   [00]   [00]   [00]   [00]
		|     4   [00]   [00]   [00]   [00]   [00]   [00]   [00]   [00]   [00]   [00]
	`)

	run("left", "1", `
		|term.Screen (5 x 10) pos (1, 4)
		|  styles
		|     0 flags 0x0000 fg (none) bg (none)
		|  buffer
		|     0  n[00]  o[00]  w[00]   [00]  i[00]  s[00]   [00]  t[00]  h[00]  e[00]
		|     1   [00]   [00]   [00]   [00]   [00]  e[00]  r[00]   [00]  o[00]  f[00]
		|     2   [00]  o[00]  u[00]  r[00]   [00]  d[00]  i[00]  s[00]  c[00]  o[00]
		|     3  n[00]  t[00]  e[00]  n[00]  t[00]   [00]   [00]   [00]   [00]   [00]
		|     4   [00]   [00]   [00]   [00]   [00]   [00]   [00]   [00]   [00]   [00]
	`)

	run("all", "2", `
        |term.Screen (5 x 10) pos (1, 4)
        |  styles
        |     0 flags 0x0000 fg (none) bg (none)
        |  buffer
        |     0  n[00]  o[00]  w[00]   [00]  i[00]  s[00]   [00]  t[00]  h[00]  e[00]
        |     1   [00]   [00]   [00]   [00]   [00]   [00]   [00]   [00]   [00]   [00]
        |     2   [00]  o[00]  u[00]  r[00]   [00]  d[00]  i[00]  s[00]  c[00]  o[00]
        |     3  n[00]  t[00]  e[00]  n[00]  t[00]   [00]   [00]   [00]   [00]   [00]
        |     4   [00]   [00]   [00]   [00]   [00]   [00]   [00]   [00]   [00]   [00]
	`)
}
