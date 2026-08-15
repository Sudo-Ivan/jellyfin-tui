package tui

import (
	"io"
	"strconv"
	"unicode/utf8"
)

// Present writes dirty cells in the Front raster to w using VT sequences.
// Adjacent dirty cells on a row are coalesced into a single cursor+SGR+text
// run. After a successful write, Front is copied into Back.
func (b *Buffer) Present(w io.Writer) error {
	if b.W == 0 || b.H == 0 {
		return nil
	}
	full := b.force || b.DirtyRatio() >= fullRedrawRatio
	buf := make([]byte, 0, presentCap)
	if full {
		buf = append(buf, cupHome...)
	}
	buf = b.appendDirty(buf, full)
	buf = append(buf, sgrReset...)
	if _, err := w.Write(buf); err != nil {
		return err
	}
	copy(b.Back, b.Front)
	b.force = false
	return nil
}

func (b *Buffer) appendDirty(buf []byte, full bool) []byte {
	cur := presentCur{x: -1, y: -1}
	for y := 0; y < b.H; y++ {
		buf, cur = b.appendRow(buf, y, full, cur)
	}
	return buf
}

type presentCur struct {
	x, y   int
	styled bool
	last   Style
}

func (b *Buffer) appendRow(buf []byte, y int, full bool, cur presentCur) ([]byte, presentCur) {
	rowF := b.Front[y*b.W : y*b.W+b.W]
	rowB := b.Back[y*b.W : y*b.W+b.W]
	x := 0
	for x < b.W {
		if rowF[x].Wide {
			x++
			continue
		}
		if !full && rowF[x].equal(rowB[x]) {
			x++
			continue
		}
		start := x
		x = runEnd(rowF, rowB, start, full)
		if cur.x != start || cur.y != y {
			buf = appendCUP(buf, start, y)
			cur.x, cur.y = start, y
		}
		st := rowF[start].Style
		if !cur.styled || st != cur.last {
			buf = appendSGR(buf, st)
			cur.last = st
			cur.styled = true
		}
		for i := start; i < x; i++ {
			if rowF[i].Wide {
				continue
			}
			buf = utf8.AppendRune(buf, rowF[i].rune())
			cur.x++
		}
	}
	return buf, cur
}

func runEnd(rowF, rowB []Cell, start int, full bool) int {
	x := start + 1
	w := len(rowF)
	for x < w {
		if rowF[x].Wide {
			x++
			continue
		}
		if !full && rowF[x].equal(rowB[x]) {
			break
		}
		if rowF[x].Style != rowF[start].Style {
			break
		}
		x++
	}
	return x
}

func appendCUP(dst []byte, x, y int) []byte {
	dst = append(dst, "\x1b["...)
	dst = strconv.AppendInt(dst, int64(y+1), 10)
	dst = append(dst, ';')
	dst = strconv.AppendInt(dst, int64(x+1), 10)
	return append(dst, 'H')
}

func appendSGR(dst []byte, s Style) []byte {
	dst = append(dst, "\x1b[0"...)
	if s.Attr&AttrBold != 0 {
		dst = append(dst, ";1"...)
	}
	if s.Attr&AttrDim != 0 {
		dst = append(dst, ";2"...)
	}
	if s.Attr&AttrItalic != 0 {
		dst = append(dst, ";3"...)
	}
	if s.Attr&AttrUnderline != 0 {
		dst = append(dst, ";4"...)
	}
	if s.Attr&AttrReverse != 0 {
		dst = append(dst, ";7"...)
	}
	if s.Fg.Set() {
		dst = append(dst, ';')
		dst = s.Fg.appendSGR(dst, true)
	}
	if s.Bg.Set() {
		dst = append(dst, ';')
		dst = s.Bg.appendSGR(dst, false)
	}
	return append(dst, 'm')
}
