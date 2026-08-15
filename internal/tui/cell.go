package tui

// Attr is a bitset of cell attributes.
type Attr uint8

const (
	AttrBold Attr = 1 << iota
	AttrDim
	AttrItalic
	AttrUnderline
	AttrReverse
)

// Style is a complete cell appearance.
type Style struct {
	Fg, Bg Color
	Attr   Attr
}

func (s Style) WithFg(c Color) Style { s.Fg = c; return s }
func (s Style) WithBg(c Color) Style { s.Bg = c; return s }
func (s Style) Bold() Style          { s.Attr |= AttrBold; return s }
func (s Style) Dim() Style           { s.Attr |= AttrDim; return s }
func (s Style) Italic() Style        { s.Attr |= AttrItalic; return s }
func (s Style) Under() Style         { s.Attr |= AttrUnderline; return s }
func (s Style) Reverse() Style       { s.Attr |= AttrReverse; return s }

// Cell is one column in the raster. Ch==0 is treated as a space on present.
// Wide is set on the trailing column of a two-cell glyph so present skips it.
type Cell struct {
	Ch    rune
	Style Style
	Wide  bool
}

func (c Cell) equal(o Cell) bool {
	return c.Ch == o.Ch && c.Wide == o.Wide &&
		c.Style.Fg == o.Style.Fg &&
		c.Style.Bg == o.Style.Bg &&
		c.Style.Attr == o.Style.Attr
}

func (c Cell) rune() rune {
	if c.Ch == 0 {
		return ' '
	}
	return c.Ch
}
