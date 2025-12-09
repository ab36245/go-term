package term

import "fmt"

type Color struct {
	size    byte
	x, y, z byte
}

func (c Color) RGB() (int, int, int) {
	return int(c.x), int(c.y), int(c.z)
}

func (c Color) Size() int {
	return int(c.size)
}

func (c Color) String() string {
	switch c.size {
	case 4, 8:
		return fmt.Sprintf("%d 0x%02x", c.size, c.x)
	case 24:
		return fmt.Sprintf("r 0x%02x g 0x%02x b 0x%02x", c.x, c.y, c.z)
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
	return Color{size: 8, x: byte(r), y: byte(g), z: byte(b)}
}
