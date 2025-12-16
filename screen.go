package term

import (
	"fmt"
	"strings"
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

func (s Screen) Height() int {
	return s.height
}

func (s Screen) Width() int {
	return s.width
}

func (s Screen) Row() int {
	return s.row
}

func (s Screen) Col() int {
	return s.col
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
