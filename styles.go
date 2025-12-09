package term

import "github.com/ab36245/go-ansi"

type Styles struct {
	byIndex []Style
	byStyle map[Style]int
	current int
}

func (s *Styles) Current() (int, Style) {
	return s.current, s.byIndex[s.current]
}

func (s *Styles) Index(style Style) int {
	if index, ok := s.byStyle[style]; ok {
		return index
	}
	index := len(s.byIndex)
	s.byIndex = append(s.byIndex, style)
	s.byStyle[style] = index
	return index
}

func (s *Styles) Style(index int) Style {
	if index < len(s.byIndex) {
		return s.byIndex[index]
	}
	return Style{}
}

func (s *Styles) ApplyCsi(csi ansi.Csi) (int, Style) {
	index, style := s.Current()
	style = style.ApplyCsi(csi)
	index = s.Index(style)
	s.current = index
	return index, style
}

func newStyles() *Styles {
	style := Style{}
	byIndex := []Style{style}
	byStyle := map[Style]int{style: 0}
	return &Styles{
		byIndex: byIndex,
		byStyle: byStyle,
		current: 0,
	}
}
