// Package tui draws each frame immediately into a retained cell raster,
// then presents only dirty runs to the host terminal.
package tui

// Buffer is a retained front/back cell raster. Immediate-mode drawing
// writes the whole frame into Front. Present copies only dirty runs to
// the terminal, then swaps Front into Back.
type Buffer struct {
	W, H  int
	Front []Cell
	Back  []Cell
	force bool
}

func NewBuffer(w, h int) *Buffer {
	b := &Buffer{}
	b.Resize(w, h)
	return b
}

func (b *Buffer) Resize(w, h int) {
	if w < 1 {
		w = 1
	}
	if h < 1 {
		h = 1
	}
	if w == b.W && h == b.H && b.Front != nil {
		return
	}
	n := w * h
	b.W, b.H = w, h
	b.Front = make([]Cell, n)
	b.Back = make([]Cell, n)
	b.force = true
}

func (b *Buffer) Bounds() Rect {
	return Rect{W: b.W, H: b.H}
}

func (b *Buffer) idx(x, y int) int { return y*b.W + x }

func (b *Buffer) At(x, y int) Cell {
	if x < 0 || y < 0 || x >= b.W || y >= b.H {
		return Cell{}
	}
	return b.Front[b.idx(x, y)]
}

func (b *Buffer) Put(x, y int, c Cell) {
	if x < 0 || y < 0 || x >= b.W || y >= b.H {
		return
	}
	b.Front[b.idx(x, y)] = c
}

func (b *Buffer) Fill(r Rect, c Cell) {
	r = r.Intersect(b.Bounds())
	for y := r.Y; y < r.MaxY(); y++ {
		row := b.Front[y*b.W : y*b.W+b.W]
		for x := r.X; x < r.MaxX(); x++ {
			row[x] = c
		}
	}
}

func (b *Buffer) Clear(style Style) {
	c := Cell{Ch: ' ', Style: style}
	for i := range b.Front {
		b.Front[i] = c
	}
}

// Invalidate forces a full rewrite on the next Present.
func (b *Buffer) Invalidate() { b.force = true }

func (b *Buffer) DirtyRatio() float64 {
	if len(b.Front) == 0 {
		return 1
	}
	n := 0
	for i := range b.Front {
		if b.force || !b.Front[i].equal(b.Back[i]) {
			n++
		}
	}
	return float64(n) / float64(len(b.Front))
}
