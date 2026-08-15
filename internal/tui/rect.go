package tui

// Rect is a half-open pixel rectangle in cell coordinates.
type Rect struct {
	X, Y, W, H int
}

func (r Rect) MaxX() int { return r.X + r.W }
func (r Rect) MaxY() int { return r.Y + r.H }

func (r Rect) Empty() bool { return r.W <= 0 || r.H <= 0 }

func (r Rect) Contains(x, y int) bool {
	return x >= r.X && x < r.MaxX() && y >= r.Y && y < r.MaxY()
}

func (r Rect) Inset(n int) Rect {
	return Rect{X: r.X + n, Y: r.Y + n, W: r.W - 2*n, H: r.H - 2*n}
}

func (r Rect) InsetXY(x, y int) Rect {
	return Rect{X: r.X + x, Y: r.Y + y, W: r.W - 2*x, H: r.H - 2*y}
}

// SplitX cuts a left strip of width n. n may be negative to measure from the right.
func (r Rect) SplitX(n int) (left, right Rect) {
	if n < 0 {
		n = r.W + n
	}
	if n < 0 {
		n = 0
	}
	if n > r.W {
		n = r.W
	}
	left = Rect{r.X, r.Y, n, r.H}
	right = Rect{r.X + n, r.Y, r.W - n, r.H}
	return
}

// SplitY cuts a top strip of height n. n may be negative to measure from the bottom.
func (r Rect) SplitY(n int) (top, bottom Rect) {
	if n < 0 {
		n = r.H + n
	}
	if n < 0 {
		n = 0
	}
	if n > r.H {
		n = r.H
	}
	top = Rect{r.X, r.Y, r.W, n}
	bottom = Rect{r.X, r.Y + n, r.W, r.H - n}
	return
}

func (r Rect) Intersect(o Rect) Rect {
	x1 := max(r.X, o.X)
	y1 := max(r.Y, o.Y)
	x2 := min(r.MaxX(), o.MaxX())
	y2 := min(r.MaxY(), o.MaxY())
	if x2 <= x1 || y2 <= y1 {
		return Rect{}
	}
	return Rect{x1, y1, x2 - x1, y2 - y1}
}
