package term

import (
	"fmt"

	"github.com/ab36245/go-ansi"
	"github.com/rs/zerolog/log"
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

func (s Style) IsItalic() bool {
	return s.Flags&sfItalic == sfItalic
}

func (s Style) IsUnderlined() bool {
	return s.Flags&sfUnderlined == sfUnderlined
}

func (s Style) IsBlink() bool {
	return s.Flags&sfBlink == sfBlink
}

func (s Style) IsInverse() bool {
	return s.Flags&sfInverse == sfInverse
}

func (s Style) IsInvisible() bool {
	return s.Flags&sfInvisible == sfInvisible
}

func (s Style) IsCrossedOut() bool {
	return s.Flags&sfCrossedOut == sfCrossedOut
}

func (s Style) IsDoublyUnderlined() bool {
	return s.Flags&sfDoublyUnderlined == sfDoublyUnderlined
}

// TODO

func (s Style) Esc() string {
	esc := ""
	add := func(s string) {
		if esc != "" {
			esc += ";"
		}
		esc += s
	}

	if s.IsBold() {
		add("1")
	}
	if s.IsFaint() {
		add("2")
	}
	if s.IsItalic() {
		add("3")
	}
	if s.IsUnderlined() {
		add("4")
	}
	if s.IsBlink() {
		add("5")
	}
	if s.IsInverse() {
		add("7")
	}
	if s.IsInvisible() {
		add("8")
	}
	if s.IsCrossedOut() {
		add("9")
	}
	if s.IsDoublyUnderlined() {
		add("21")
	}
	if s.Fg.IsValid() {
		add("3" + s.Fg.Esc())
	}
	if s.Bg.IsValid() {
		add("4" + s.Bg.Esc())
	}
	return "\x1b[" + esc + "m"
}

func (s Style) ApplySgr(sgr ansi.Csi) Style {
	if sgr.Action() != 'm' {
		return s
	}
	log.Debug().Stringer("sgr", sgr).Msg("applying sgr")
	new := s.Copy()
	params := sgr.Params()
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
			new.Fg = NewColor4(n - 30)
		case n == 38:
			switch n := params.Get(0); n {
			case 2:
				r := params.Get(0)
				g := params.Get(0)
				b := params.Get(0)
				new.Fg = NewColor24(r, g, b)
			case 5:
				v := params.Get(0)
				new.Fg = NewColor8(v)
			default:
				// TODO log warning?
			}
		case n == 39:
			new.Fg = Color{}

		case n >= 40 && n <= 47:
			new.Bg = NewColor4(n - 40)
		case n == 48:
			switch n := params.Get(0); n {
			case 2:
				r := params.Get(0)
				g := params.Get(0)
				b := params.Get(0)
				new.Bg = NewColor24(r, g, b)
			case 5:
				v := params.Get(0)
				new.Bg = NewColor8(v)
			default:
				// TODO log warning?
			}
		case n == 49:
			new.Bg = Color{}

		case n >= 90 && n <= 97:
			// TODO

		case n >= 100 && n <= 107:
			// TODO
		default:
			log.Warn().Stringer("sgr", sgr).Msg("unhandled sgr")
		}
	}
	log.Debug().Stringer("style", new).Msg("new style")
	return new
}

func (s Style) Copy() Style {
	return Style{
		s.Flags,
		s.Fg,
		s.Bg,
	}
}

func (s Style) String() string {
	str := fmt.Sprintf("flags 0x%04x", s.Flags)
	str += fmt.Sprintf(" fg %s", s.Fg)
	str += fmt.Sprintf(" bg %s", s.Bg)
	return str
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
