package term

import "fmt"

type Color struct {
	size    byte
	x, y, z byte
}

func (c Color) Esc() string {
	switch c.size {
	case 4:
		return fmt.Sprintf("%d", c.x)
	case 24:
		return fmt.Sprintf("8;2;%d;%d;%d", c.x, c.y, c.z)
	default:
		// TODO
		return ""
	}
}

func (c Color) IsValid() bool {
	return c.size == 4 || c.size == 8 || c.size == 24
}

func (c Color) RGB() (int, int, int) {
	return int(c.x), int(c.y), int(c.z)
}

func (c Color) Size() int {
	return int(c.size)
}

func (c Color) String() string {
	switch c.size {
	case 0:
		return "(none)"
	case 4, 8:
		return fmt.Sprintf("%d 0x%02x", c.size, c.x)
	case 24:
		return fmt.Sprintf("rgb #%02X%02X%02x", c.x, c.y, c.z)
	default:
		return fmt.Sprintf("(bad color size %d)", c.size)
	}
}

func (c Color) Value() int {
	return int(c.x)
}

func newColor4(n int) Color {
	return Color{size: 4, x: byte(n)}
}

func newColor8(n int) Color {
	return Color{size: 8, x: byte(n)}
}

func newColor24(r, g, b int) Color {
	return Color{size: 24, x: byte(r), y: byte(g), z: byte(b)}
}
