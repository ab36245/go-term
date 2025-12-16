package term

func NewBuffer(width, height int) *Buffer {
	b := &Buffer{
		height: height,
		width:  width,
		cells:  make([]BCell, width*height),
	}
	b.ClearAll()
	return b
}

type Buffer struct {
	height int
	width  int
	cells  []BCell
}

type BCell struct {
	value rune
	style int
}

func (b *Buffer) Height() int {
	return b.height
}

func (b *Buffer) Width() int {
	return b.width
}

func (b *Buffer) AssignCell(row, col int, value rune, style int) {
	index := b.index(row, col)
	b.cells[index] = BCell{value, style}
}

func (b *Buffer) AssignString(row, col int, str string, style int) {
	index := b.index(row, col)
	for _, value := range str {
		b.cells[index] = BCell{value, style}
		index++
		if index >= len(b.cells) {
			break
		}
	}
}

func (b *Buffer) ClearAll() {
	for index := range len(b.cells) {
		b.cells[index] = BCell{' ', 0}
	}
}

func (b *Buffer) ClearCell(row, col int) {
	index := b.index(row, col)
	b.cells[index] = BCell{' ', 0}
}

func (b *Buffer) ClearRange(startRow, startCol, endRow, endCol int) {
	startIndex := b.index(startRow, startCol)
	endIndex := b.index(endRow, endCol)
	for i := startIndex; i <= endIndex; i++ {
		b.cells[i] = BCell{' ', 0}
	}
}

func (b *Buffer) CopyRange(startRow, startCol, endRow, endCol, toRow, toCol int) {
	startIndex := b.index(startRow, startCol)
	endIndex := b.index(endRow, endCol)
	size := endIndex - startIndex + 1
	toIndex := b.index(toRow, toCol)
	if startIndex < toIndex {
		for i := size - 1; i >= 0; i-- {
			b.cells[toIndex+i] = b.cells[startIndex+i]
		}
	} else {
		for i := range size {
			b.cells[toIndex+i] = b.cells[startIndex+i]
		}
	}
}

func (b *Buffer) View(styles *Styles) []string {
	var view []string
	index := 0
	style := 0
	for range b.height {
		line := ""
		for range b.width {
			cell := b.cells[index]
			index++
			if cell.value != 0 {
				if style != cell.style {
					line += styles.Style(cell.style).Esc()
					style = cell.style
				}
				line += string(cell.value)
			}
		}
		if style != 0 {
			line += styles.Style(0).Esc()
			style = 0
		}
		view = append(view, line)
	}
	return view
}

func (b *Buffer) index(row, col int) int {
	if row < 0 {
		row = 0
	} else if row > b.height-1 {
		row = b.height - 1
	}
	if col < 0 {
		col = 0
	} else if col > b.width-1 {
		col = b.width - 1
	}
	return row*b.width + col
}
