package term

import (
	"fmt"
	"strings"

	"github.com/mattn/go-runewidth"
	"github.com/rs/zerolog/log"

	"github.com/ab36245/go-ansi"
)

func New(width, height int) *Screen {
	s := &Screen{
		height:       height,
		width:        width,
		buffer:       NewBuffer(width, height),
		styles:       NewStyles(),
		row:          0,
		col:          0,
		scrollTop:    0,
		scrollBottom: height - 1,
	}
	return s
}

type Screen struct {
	height       int
	width        int
	buffer       *Buffer
	styles       *Styles
	row          int
	col          int
	scrollTop    int
	scrollBottom int
}

func (s *Screen) ApplySeq(seq ansi.Sequence) {
	switch seq := seq.(type) {
	case ansi.Csi:
		s.ApplyCsi(seq)
	case ansi.Ctrl:
		s.ApplyCtrl(seq)
	case ansi.Text:
		s.ApplyText(seq)
	}
}

func (s *Screen) ApplyCsi(csi ansi.Csi) {
	log.Debug().Stringer("csi", csi).Msg("Applying csi")
	params := csi.Params()
	switch csi.Action() {
	case 'A':
		// Cursor Up (CUU)
		n := params.Get(1)
		s.GoTo(s.row-n, s.col)

	case 'B':
		// Cursor Down (CUD)
		n := params.Get(1)
		s.GoTo(s.row+n, s.col)

	case 'C':
		// Cursor Forward (CUF)
		n := params.Get(1)
		s.GoTo(s.row, s.col-n)

	case 'D':
		// Cursor Backward (CUB)
		n := params.Get(1)
		s.GoTo(s.row, s.col+n)

	case 'G':
		// Cursor Character Absolute (CHA)
		n := params.Get(1)
		s.col = n

	case 'H':
		// Cursor Position (CUP)
		r := params.Get(1)
		c := params.Get(1)
		s.GoTo(r-1, c-1)

	case 'J':
		// Erase in Display (ED)
		n := params.Get(0)
		switch n {
		case 0:
			s.buffer.ClearRange(s.row, s.col, s.height-1, s.width-1)
		case 1:
			s.buffer.ClearRange(0, 0, s.row, s.col)
		case 2:
			s.buffer.ClearAll()
		}

	case 'K':
		// Erase in Line (EL)
		n := params.Get(0)
		switch n {
		case 0:
			s.buffer.ClearRange(s.row, s.col, s.row, s.width-1)
		case 1:
			s.buffer.ClearRange(s.row, 0, s.row, s.col)
		case 2:
			s.buffer.ClearRange(s.row, 0, s.row, s.width-1)
		}

	case 'L':
		// Insert Line (IL)
		// TODO

	case 'M':
		// Delete Line (DL)
		n := params.Get(1)
		room := s.scrollBottom - s.row + 1
		if n > room {
			n = room
		}
		startRow := s.row + n
		startCol := 0
		endRow := s.scrollBottom
		endCol := s.width - 1
		if startRow <= endRow {
			toRow := s.row
			toCol := 0
			s.buffer.CopyRange(startRow, startCol, endRow, endCol, toRow, toCol)
		}
		startRow = endRow - n + 1
		s.buffer.ClearRange(startRow, startCol, endRow, endCol)

	case 'd':
		// Line Position Absolute (VPA)
		n := params.Get(1)
		s.GoTo(n-1, s.col)

	case 'm':
		// Character Attributes (SGR)
		s.styles.ApplySgr(csi)

	case 'r':
		t := params.Get(1)
		b := params.Get(0)
		if b == 0 {
			b = s.height
		}
		s.scrollTop = t - 1
		s.scrollBottom = b - 1

	default:
		log.Warn().Stringer("csi", csi).Msg("unhandled csi")
	}
}

func (s *Screen) ApplyCtrl(ctrl ansi.Ctrl) {
	log.Debug().
		Stringer("ctrl", ctrl).
		Msg("ApplyCtrl")
	switch ctrl {
	case '\n':
		// s.GoTo(s.row+1, 0)
		s.GoTo(s.row, s.width)
	case '\r':
		s.GoTo(s.row, 0)
	}
}

func (s *Screen) ApplyRune(r rune) {
	if s.col >= s.width {
		s.row++
		s.col = 0
	}
	width := runewidth.RuneWidth(r)
	if s.col+width > s.width {
		s.row++
		s.col = 0
	}
	s.buffer.AssignCell(s.row, s.col, r, s.styles.current)
	for i := 1; i < width; i++ {
		s.buffer.AssignCell(s.row, s.col, 0, s.styles.current)
	}
	s.col += width
}

func (s *Screen) ApplyText(text ansi.Text) {
	log.Debug().
		Str("text", string(text)).
		Int("len", len(text)).
		Int("row", s.row).
		Int("col", s.col).
		Msg("ApplyText")
	for _, r := range text {
		s.ApplyRune(r)
	}
}

func (s *Screen) GoTo(row, col int) {
	if row < 0 {
		row = 0
	} else if row >= s.height {
		row = s.height - 1
	}
	s.row = row

	if col < 0 {
		col = 0
	} else if col >= s.width {
		col = s.width - 1
	}
	s.col = col
}

func (s *Screen) Dump() string {
	lines := []string{}

	add := func(mesg string, args ...any) {
		lines = append(lines, fmt.Sprintf(mesg, args...))
	}

	{
		head := fmt.Sprintf("%T", *s)
		head += fmt.Sprintf(" (%d x %d)", s.height, s.width)
		head += fmt.Sprintf(" pos (%d, %d)", s.row, s.col)
		if s.scrollTop != 0 || s.scrollBottom != s.height-1 {
			head += fmt.Sprintf(" window %d-%d", s.scrollTop, s.scrollBottom)
		}
		add(head)
	}

	{
		add("  styles")
		for i, s := range s.styles.byIndex {
			add("    %2d %s", i, s)
		}
	}

	{
		add("  buffer")
		i := 0
		for r := range s.height {
			str := fmt.Sprintf("    %2d", r)
			for range s.width {
				c := s.buffer.cells[i]
				i++
				str += fmt.Sprintf(" %2c[%02d]", c.value, c.style)
			}
			add(str)
		}
	}

	return strings.Join(lines, "\n")

}

func (s *Screen) View() []string {
	return s.buffer.View(s.styles)
}
