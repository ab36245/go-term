package term

import (
	"github.com/ab36245/go-ansi"
)

type Style struct {
	Flags StyleFlag
	Fg    Color
	Bg    Color
}

func (s Style) IsBold() bool {
	return s.Flags&sfBold == sfBold
}

func (s Style) IsFaint() bool {
	return s.Flags&sfFaint == sfFaint
}

// TODO

func (s Style) ApplyCsi(csi ansi.Csi) Style {
	if csi.Action() != 'm' {
		return s
	}
	new := s.Copy()
	params := csi.Params()
	for params.More() {
		n := params.Get(0)
		switch {
		case n == 0:
			new = Style{}

		case n == 1:
			new.Flags |= sfBold
			new.Flags &= ^sfFaint
		case n == 2:
			new.Flags &= ^sfBold
			new.Flags |= sfFaint
		case n == 3:
			new.Flags |= sfItalic
		case n == 4:
			new.Flags |= sfUnderlined
		case n == 5:
			new.Flags |= sfBlink
		case n == 7:
			new.Flags |= sfInverse
		case n == 8:
			new.Flags |= sfInvisible
		case n == 9:
			new.Flags |= sfCrossedOut

		case n == 21:
			new.Flags |= sfDoublyUnderlined
		case n == 22:
			new.Flags &= ^sfBold
			new.Flags &= ^sfFaint
		case n == 23:
			new.Flags &= ^sfItalic
		case n == 24:
			new.Flags &= ^sfUnderlined
		case n == 25:
			new.Flags &= ^sfBlink
		case n == 27:
			new.Flags &= ^sfInverse
		case n == 28:
			new.Flags &= ^sfInvisible
		case n == 29:
			new.Flags &= ^sfCrossedOut

		case n >= 30 && n <= 37:
			new.Fg = newColor4(n - 30)
		case n == 38:
			switch n := params.Get(0); n {
			case 2:
				r := params.Get(0)
				g := params.Get(0)
				b := params.Get(0)
				new.Fg = newColor24(r, g, b)
			case 5:
				v := params.Get(0)
				new.Fg = newColor8(v)
			default:
				// TODO log warning?
			}
		case n == 39:
			new.Fg = Color{}

		case n >= 40 && n <= 47:
			new.Bg = newColor4(n - 40)
		case n == 48:
			switch n := params.Get(0); n {
			case 2:
				r := params.Get(0)
				g := params.Get(0)
				b := params.Get(0)
				new.Bg = newColor24(r, g, b)
			case 5:
				v := params.Get(0)
				new.Bg = newColor8(v)
			default:
				// TODO log warning?
			}
		case n == 49:
			new.Bg = Color{}

		case n >= 90 && n <= 97:
			// TODO

		case n >= 100 && n <= 107:
			// TODO
		}
	}
	return new
}

func (s Style) Copy() Style {
	return Style{
		s.Flags,
		s.Fg,
		s.Bg,
	}
}

type StyleFlag uint16

const (
	sfBold StyleFlag = 1 << iota
	sfFaint
	sfItalic
	sfUnderlined
	sfBlink
	sfInverse
	sfInvisible
	sfCrossedOut
	sfDoublyUnderlined
)
