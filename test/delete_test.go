package test

import (
	"strings"
	"testing"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func TestDelete(t *testing.T) {
	log.Logger = log.Logger.Level(zerolog.WarnLevel)

	run := func(name string, p string, e string) {
		t.Run(name, func(t *testing.T) {
			screen := screenMake(10, 5)
			builder := strings.Builder{}
			builder.WriteString("0 abcdefgh")
			builder.WriteString("1 lmnopqrs")
			builder.WriteString("2 01234567")
			builder.WriteString("3 ABCDEFGH")
			builder.WriteString("4 LMNOPQRS")
			builder.WriteString("\x1b[2;5H")
			builder.WriteString("\x1b[" + p + "M")
			a := screenInput(screen, builder.String())
			screenCheck(t, a, e)
		})
	}

	run("one", "", `
		|term.Screen (5 x 10) pos (1, 4)
        |  styles
        |     0 flags 0x0000 fg (none) bg (none)
        |  buffer
        |     0  0[00]   [00]  a[00]  b[00]  c[00]  d[00]  e[00]  f[00]  g[00]  h[00]
        |     1  2[00]   [00]  0[00]  1[00]  2[00]  3[00]  4[00]  5[00]  6[00]  7[00]
        |     2  3[00]   [00]  A[00]  B[00]  C[00]  D[00]  E[00]  F[00]  G[00]  H[00]
        |     3  4[00]   [00]  L[00]  M[00]  N[00]  O[00]  P[00]  Q[00]  R[00]  S[00]
        |     4   [00]   [00]   [00]   [00]   [00]   [00]   [00]   [00]   [00]   [00]
	`)

	run("three", "3", `
		|term.Screen (5 x 10) pos (1, 4)
        |  styles
        |     0 flags 0x0000 fg (none) bg (none)
        |  buffer
        |     0  0[00]   [00]  a[00]  b[00]  c[00]  d[00]  e[00]  f[00]  g[00]  h[00]
        |     1  4[00]   [00]  L[00]  M[00]  N[00]  O[00]  P[00]  Q[00]  R[00]  S[00]
        |     2   [00]   [00]   [00]   [00]   [00]   [00]   [00]   [00]   [00]   [00]
        |     3   [00]   [00]   [00]   [00]   [00]   [00]   [00]   [00]   [00]   [00]
        |     4   [00]   [00]   [00]   [00]   [00]   [00]   [00]   [00]   [00]   [00]
	`)

	run("nine", "9", `
		|term.Screen (5 x 10) pos (1, 4)
        |  styles
        |     0 flags 0x0000 fg (none) bg (none)
        |  buffer
        |     0  0[00]   [00]  a[00]  b[00]  c[00]  d[00]  e[00]  f[00]  g[00]  h[00]
        |     1   [00]   [00]   [00]   [00]   [00]   [00]   [00]   [00]   [00]   [00]
        |     2   [00]   [00]   [00]   [00]   [00]   [00]   [00]   [00]   [00]   [00]
        |     3   [00]   [00]   [00]   [00]   [00]   [00]   [00]   [00]   [00]   [00]
        |     4   [00]   [00]   [00]   [00]   [00]   [00]   [00]   [00]   [00]   [00]
	`)
}
