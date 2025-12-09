package term

import (
	"github.com/mattn/go-runewidth"

	"github.com/ab36245/go-ansi"
)

func New(width, height int) *Screen {
	return &Screen{
		height: height,
		width:  width,
		cells:  make([][]Cell, height),
		styles: newStyles(),
		top:    0,
		row:    0,
		col:    0,
	}
}

type Screen struct {
	height int
	width  int
	cells  [][]Cell
	styles *Styles
	top    int
	row    int
	col    int
}

func (s *Screen) ApplyCsi(csi ansi.Csi) {
	params := csi.Params()
	switch csi.Action() {
	case 'A':
		n := params.Get(1)
		s.GoTo(s.row-n, s.col)

	case 'B':
		n := params.Get(1)
		s.GoTo(s.row+n, s.col)

	case 'C':
		n := params.Get(1)
		s.GoTo(s.row, s.col-n)

	case 'D':
		n := params.Get(1)
		s.GoTo(s.row, s.col+n)

	case 'H':
		r := params.Get(1)
		c := params.Get(1)
		s.GoTo(r-1, c-1)

	case 'm':
		s.styles.ApplyCsi(csi)
	}
}

func (s *Screen) ApplyCtrl(ctrl ansi.Ctrl) {
	switch ctrl {
	case '\n':
	case '\r':
		s.col = 0
	}
}

func (s *Screen) ApplyText(text ansi.Text) {
	for _, r := range text {
		s.ApplyRune(r)
	}
}

func (s *Screen) ApplyRune(r rune) {
	for len(s.cells) <= s.row {
		s.cells = append(s.cells, nil)
	}
	for len(s.cells[s.row]) <= s.col {
		s.cells[s.row] = append(s.cells[s.row], Cell{' ', 0})
	}
	line := s.cells[s.row]
	span := 0
	for i := 0; i < s.col && i < len(line); i++ {
		span += runewidth.RuneWidth(line[i].value)
	}

	width := runewidth.RuneWidth(r)
	if s.col+width >= s.width {
		s.Newline()
	}

	s.cells[s.row][s.col] = Cell{r, s.styles.current}
	s.col++
	for i := 1; i < width; i++ {
		s.cells[s.row][s.col] = Cell{0, s.styles.current}
		s.col++
	}
}

func (s *Screen) GoTo(row, col int) {
	if row < s.top {
		row = s.top
	} else if row >= s.top+s.height {
		row = s.top + s.height - 1
	}
	s.MoveTo(row, col)
}

func (s *Screen) MoveTo(row, col int) {
	min := row - s.height + 1
	max := row
	if s.top < min {
		s.top = min
	} else if s.top > max {
		s.top = row
	}
	s.row = row

	if col < 0 {
		col = 0
	} else if col >= s.width {
		col = s.width - 1
	}
	s.col = col
}

func (s *Screen) Newline() {
	s.col = 0
	last := s.top + s.height - 1
	if s.row < last {
		s.row++
	} else {
		s.ScrollDown(1)
		s.row = last
	}
}

func (s *Screen) ScrollDown(n int) {
	if n >= s.height {
		n = s.height
	}
	s.cells = s.cells[n:]
	for range n {
		s.cells = append(s.cells, nil)
	}
}

// func newRow(cols int) []Cell {
// 	row := make([]Cell, cols)
// 	for c := range cols {
// 		row[c] = Cell{' ', 0}
// 	}
// 	return row
// }
