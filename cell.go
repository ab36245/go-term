package term

type Cell struct {
	value rune
	style int // index into Styles table
}

// func (c Cell) Width() int {
// 	return runewidth.RuneWidth(c.value)
// }
