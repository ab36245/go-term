package term

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
