package tui

// Painter is the immediate-mode drawing surface. Every frame the app
// stamps glyphs into the Front raster through a clip stack. Nothing
// here retains widgets between frames.
type Painter struct {
	buf   *Buffer
	theme Theme
	clip  Rect
}

func NewPainter(buf *Buffer, theme Theme) *Painter {
	return &Painter{buf: buf, theme: theme, clip: buf.Bounds()}
}

func (p *Painter) Theme() Theme { return p.theme }
func (p *Painter) Bounds() Rect { return p.clip }
func (p *Painter) Size() (w, h int) {
	return p.buf.W, p.buf.H
}

func (p *Painter) WithClip(r Rect) *Painter {
	c := *p
	c.clip = r.Intersect(p.clip)
	return &c
}

func (p *Painter) Clear() {
	p.buf.Clear(p.theme.Body())
}

func (p *Painter) Fill(r Rect, st Style) {
	r = r.Intersect(p.clip)
	p.buf.Fill(r, Cell{Ch: ' ', Style: st})
}

func (p *Painter) Put(x, y int, r rune, st Style) {
	if !p.clip.Contains(x, y) {
		return
	}
	w := RuneWidth(r)
	if w <= 0 {
		return
	}
	p.buf.Put(x, y, Cell{Ch: r, Style: st})
	if w == 2 && p.clip.Contains(x+1, y) {
		p.buf.Put(x+1, y, Cell{Ch: 0, Style: st, Wide: true})
	}
}

func (p *Painter) Text(x, y int, s string, st Style) {
	cx := x
	for _, r := range s {
		w := RuneWidth(r)
		if w == 0 {
			continue
		}
		p.Put(cx, y, r, st)
		cx += w
	}
}

func (p *Painter) TextClip(r Rect, s string, st Style) {
	r = r.Intersect(p.clip)
	if r.Empty() {
		return
	}
	p.Fill(r, st)
	s = ellipsize(s, r.W)
	p.Text(r.X, r.Y, s, st)
}

func (p *Painter) HLine(x, y, w int, r rune, st Style) {
	for i := range w {
		p.Put(x+i, y, r, st)
	}
}

func (p *Painter) VLine(x, y, h int, r rune, st Style) {
	for i := range h {
		p.Put(x, y+i, r, st)
	}
}

func (p *Painter) Frame(r Rect, title string, st Style) {
	if r.W < 2 || r.H < 2 {
		return
	}
	p.Put(r.X, r.Y, '╭', st)
	p.Put(r.MaxX()-1, r.Y, '╮', st)
	p.Put(r.X, r.MaxY()-1, '╰', st)
	p.Put(r.MaxX()-1, r.MaxY()-1, '╯', st)
	p.HLine(r.X+1, r.Y, r.W-2, '─', st)
	p.HLine(r.X+1, r.MaxY()-1, r.W-2, '─', st)
	p.VLine(r.X, r.Y+1, r.H-2, '│', st)
	p.VLine(r.MaxX()-1, r.Y+1, r.H-2, '│', st)
	if title != "" && r.W > 4 {
		label := " " + ellipsize(title, r.W-6) + " "
		p.Text(r.X+2, r.Y, label, st)
	}
}

func (p *Painter) Gauge(r Rect, frac float64, fill, empty Style) {
	r = r.Intersect(p.clip)
	if r.Empty() {
		return
	}
	if frac < 0 {
		frac = 0
	}
	if frac > 1 {
		frac = 1
	}
	filled := int(frac * float64(r.W))
	for x := 0; x < r.W; x++ {
		st := empty
		ch := '░'
		if x < filled {
			st = fill
			ch = '█'
		}
		p.Put(r.X+x, r.Y, ch, st)
	}
}

func (p *Painter) Swatch(x, y int, c Color) {
	st := Style{Fg: c, Bg: p.theme.Void}
	p.Put(x, y, '▍', st)
}
