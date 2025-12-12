package term

import (
	"github.com/mattn/go-runewidth"
	"github.com/rs/zerolog/log"

	"github.com/ab36245/go-ansi"
)

func New(width, height int) *Screen {
	s := &Screen{
		height:       height,
		width:        width,
		cells:        nil,
		styles:       newStyles(),
		top:          0,
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
	cells        [][]Cell
	styles       *Styles
	top          int
	row          int
	col          int
	scrollTop    int
	scrollBottom int
}

func (s *Screen) ApplyCsi(csi ansi.Csi) {
	log.Debug().Stringer("csi", csi).Msg("Applying csi")
	params := csi.Params()
	switch csi.Action() {
	case 'A':
		n := params.Get(1)
		s.GoToInView(s.row-n, s.col)

	case 'B':
		n := params.Get(1)
		s.GoToInView(s.row+n, s.col)

	case 'C':
		n := params.Get(1)
		s.GoToInView(s.row, s.col-n)

	case 'D':
		n := params.Get(1)
		s.GoToInView(s.row, s.col+n)

	case 'G':
		n := params.Get(1)
		s.col = n

	case 'H':
		r := params.Get(1)
		c := params.Get(1)
		s.GoToInView(r-1, c-1)

	case 'J':
		n := params.Get(0)
		switch n {
		case 0:
			s.EraseDisplayBelow()
		case 1:
			s.EraseDisplayAbove()
		case 2:
			s.EraseDisplayAll()
		}

	case 'K':
		n := params.Get(0)
		switch n {
		case 0:
			s.EraseLineRight()
		case 1:
			s.EraseLineLeft()
		case 2:
			s.EraseLineAll()
		}

	case 'L':
		// Insert
		// TODO

	case 'M':
		n := params.Get(1)
		s.DeleteRows(n)

	case 'd':
		n := params.Get(1)
		s.row = s.top + n - 1

	case 'm':
		s.styles.ApplyCsi(csi)

	case 'r':
		t := params.Get(1)
		b := params.Get(0)
		if b == 0 {
			b = s.height - 1
		}
		s.scrollTop = t
		s.scrollBottom = b

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
		s.GoTo(s.row+1, 0)
		// s.GoTo(s.row, s.width)
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
	s.Ensure(s.row, s.col+width-1)
	s.cells[s.row][s.col] = Cell{r, s.styles.current}
	for i := 1; i < width; i++ {
		s.cells[s.row][s.col+i] = Cell{0, s.styles.current}
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

func (s *Screen) DeleteRows(n int) {
	if s.row < s.scrollTop {
		return
	}
	if s.row > s.scrollBottom {
		return
	}
	n = min(n, s.scrollBottom-s.row+1)
	for r := s.row; r <= s.scrollBottom; r++ {
		if r >= len(s.cells) {
			break
		}
		if r+n <= s.scrollBottom {
			s.cells[r] = s.cells[r+n]
		} else {
			s.cells[r] = nil
		}
	}
}

func (s *Screen) EraseDisplayAbove() {
	log.Warn().Str("err", "not implemented").Msg("EraseDisplayAbove")
}

func (s *Screen) EraseDisplayAll() {
	s.cells = nil
}

func (s *Screen) EraseDisplayBelow() {
	// log.Warn().Str("err", "not implemented").Msg("EraseDisplayBelow")
	if s.row < len(s.cells) && s.cells[s.row] != nil {
		s.cells[s.row] = s.cells[s.row][:s.col]
	}
	for i := s.row + 1; i < s.top+s.height; i++ {
		if i < len(s.cells) {
			s.cells[i] = nil
		}
	}
}

func (s *Screen) EraseLineAll() {
	s.cells[s.row] = nil
}

func (s *Screen) EraseLineLeft() {
	for col := range s.col {
		if col < len(s.cells[s.row]) {
			s.cells[s.row][col] = Cell{' ', 0}
		}
	}
}

func (s *Screen) EraseLineRight() {
	s.cells[s.row] = s.cells[s.row][:s.col]
}

func (s *Screen) GoTo(row, col int) {
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
	s.Ensure(s.row, s.col)
}

func (s *Screen) GoToInView(row, col int) {
	if row < s.top {
		row = s.top
	} else if row >= s.top+s.height {
		row = s.top + s.height - 1
	}
	// col is checked in GoTo
	s.GoTo(row, col)
}

func (s *Screen) Ensure(row, col int) {
	for len(s.cells) <= row {
		s.cells = append(s.cells, nil)
	}
	for len(s.cells[row]) <= col {
		s.cells[row] = append(s.cells[row], Cell{' ', 0})
	}
}

func (s *Screen) View() string {
	str := ""
	curStyle := -1
	for row := s.top; row < min(len(s.cells), s.top+s.height); row++ {
		for _, cell := range s.cells[row] {
			if curStyle != cell.style {
				str += "\x1b[0m"
				str += s.styles.Style(cell.style).Esc()
				curStyle = cell.style
			}
			str += string(cell.value)
		}
		str += "\n"
	}
	return str
}
