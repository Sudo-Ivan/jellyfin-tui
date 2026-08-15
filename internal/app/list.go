package app

import "jellyfin-tui/internal/tui"

type rowList struct {
	labels []string
	ids    []string
	sel    int
	off    int
	active bool
	hit    func(tui.Rect, int)
}

func paintRows(p *tui.Painter, r tui.Rect, list rowList, th tui.Theme) int {
	vis := r.H
	if vis < 1 {
		return list.off
	}
	n := len(list.labels)
	sel := list.sel
	off := list.off
	if sel < 0 {
		sel = 0
	}
	if n > 0 && sel >= n {
		sel = n - 1
	}
	if sel < off {
		off = sel
	}
	if sel >= off+vis {
		off = sel - vis + 1
	}
	if off < 0 {
		off = 0
	}
	for i := range vis {
		paintRow(p, r, rowPaint{list: list, idx: off + i, rowOff: i, sel: sel, n: n}, th)
	}
	return off
}

type rowPaint struct {
	list        rowList
	idx, rowOff int
	sel, n      int
}

func paintRow(p *tui.Painter, r tui.Rect, rp rowPaint, th tui.Theme) {
	row := tui.Rect{X: r.X, Y: r.Y + rp.rowOff, W: r.W, H: 1}
	if rp.idx >= rp.n {
		p.Fill(row, th.Body())
		return
	}
	st := th.Body()
	if rp.idx == rp.sel && rp.list.active {
		st = th.Focus()
	} else if rp.idx == rp.sel {
		st = th.Accent().WithBg(th.Well)
	}
	p.Fill(row, st)
	if rp.list.hit != nil {
		rp.list.hit(row, rp.idx)
	}
	if rp.idx < len(rp.list.ids) && rp.list.ids[rp.idx] != "" {
		p.Put(row.X, row.Y, swatchCh, st.WithFg(ink(rp.list.ids[rp.idx], th)))
		p.TextClip(tui.Rect{X: row.X + 2, Y: row.Y, W: row.W - 3, H: 1}, rp.list.labels[rp.idx], st)
		return
	}
	p.TextClip(tui.Rect{X: row.X + 1, Y: row.Y, W: row.W - 2, H: 1}, rp.list.labels[rp.idx], st)
}

func moveSel(sel, n int, e tui.Event) int {
	if n <= 0 {
		return 0
	}
	switch {
	case e.Is(tui.KeyUp) || e.RuneKey('k'):
		sel--
	case e.Is(tui.KeyDown) || e.RuneKey('j'):
		sel++
	case e.Is(tui.KeyPgUp):
		sel -= pageJump
	case e.Is(tui.KeyPgDn):
		sel += pageJump
	case e.Is(tui.KeyHome):
		sel = 0
	case e.Is(tui.KeyEnd):
		sel = n - 1
	}
	if sel < 0 {
		sel = 0
	}
	if sel >= n {
		sel = n - 1
	}
	return sel
}

func ink(id string, th tui.Theme) tui.Color {
	var h uint32 = fnvOffset
	for i := 0; i < len(id); i++ {
		h ^= uint32(id[i])
		h *= fnvPrime
	}
	pal := []tui.Color{th.Ember, th.Arc, th.Bloom, th.Signal, th.Filament}
	return pal[int(h)%len(pal)]
}
