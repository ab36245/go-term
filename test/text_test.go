package test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func TestText(t *testing.T) {
	log.Logger = log.Logger.Level(zerolog.WarnLevel)

	run := func(name string, e string) {
		t.Run(name, func(t *testing.T) {
			screen := screenMake(21, 10)
			builder := strings.Builder{}
			c := rune(0x30a1)
			for r := range screen.Height() {
				s := fmt.Sprintf("%02d ", r)
				builder.WriteString(s)
				for range (screen.Width() - len(s)) / 2 {
					builder.WriteRune(c)
					c++
				}
			}
			screenInput(screen, builder.String())
			av := ""
			for _, r := range screen.View() {
				av += "|"
				av += r
				av += "|"
				av += "\n"
			}
			if av != e {
				s := "\n"
				s += "expected:\n"
				s += e
				s += "actual:\n"
				s += av
				t.Fatal(s)
			}
		})
	}

	run("one", `
		
	`)
}
