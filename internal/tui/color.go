package tui

import "strconv"

// Color is a 24-bit sRGB value. Zero means "unset" and the present
// path will omit that channel from the SGR sequence.
type Color uint32

const (
	ColorUnset Color = 0
	ColorRGB         = 1 << 24
)

// RGB packs a truecolor value. The high bit marks it as set so black
// is distinct from unset.
func RGB(r, g, b uint8) Color {
	return Color(ColorRGB | uint32(r)<<16 | uint32(g)<<8 | uint32(b))
}

func colorHex(n uint32) Color {
	return RGB(byteOf(n>>16), byteOf(n>>8), byteOf(n))
}

func byteOf(n uint32) uint8 {
	return uint8(n & byteMask) // #nosec G115
}

func (c Color) Set() bool { return c&ColorRGB != 0 }

func (c Color) RGB() (r, g, b uint8) {
	v := uint32(c)
	return byteOf(v >> 16), byteOf(v >> 8), byteOf(v)
}

func (c Color) appendSGR(dst []byte, fg bool) []byte {
	if !c.Set() {
		return dst
	}
	r, g, b := c.RGB()
	if fg {
		dst = append(dst, sgrFgTrue...)
	} else {
		dst = append(dst, sgrBgTrue...)
	}
	dst = strconv.AppendInt(dst, int64(r), 10)
	dst = append(dst, ';')
	dst = strconv.AppendInt(dst, int64(g), 10)
	dst = append(dst, ';')
	dst = strconv.AppendInt(dst, int64(b), 10)
	return dst
}

// Mix blends c toward other by t in 0..256.
func (c Color) Mix(other Color, t int) Color {
	if !c.Set() {
		return other
	}
	if !other.Set() {
		return c
	}
	if t < 0 {
		t = 0
	}
	if t > mixMax {
		t = mixMax
	}
	r1, g1, b1 := c.RGB()
	r2, g2, b2 := other.RGB()
	return RGB(mixCh(r1, r2, t), mixCh(g1, g2, t), mixCh(b1, b2, t))
}

func mixCh(a, b uint8, t int) uint8 {
	v := (int(a)*(mixMax-t) + int(b)*t) / mixMax
	if v < 0 {
		return 0
	}
	if v > byteMax {
		return byteMax
	}
	return uint8(v) // #nosec G115
}
