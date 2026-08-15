package app

import (
	"strings"

	"jellyfin-tui/internal/tui"
)

type homeLine struct {
	header bool
	text   string
	id     string
	idx    int
}

func (a *App) paintHome(p *tui.Painter, th tui.Theme) {
	r := p.Bounds()
	left, right := r.SplitX(railWidth)
	p.Fill(left, th.Panel())
	p.VLine(left.MaxX()-1, left.Y, left.H, '│', th.Chrome().WithBg(th.Basin))
	labels := make([]string, len(a.home.rail))
	ids := make([]string, len(a.home.rail))
	for i, it := range a.home.rail {
		labels[i] = it.Name
		ids[i] = it.ID
	}
	rail := left.InsetXY(1, 1)
	a.home.railOff = paintRows(p.WithClip(rail), rail, rowList{
		labels: labels, ids: ids, sel: a.home.railSel, off: a.home.railOff, active: a.home.focus == 0,
		hit: a.track(hitRail),
	}, th)
	if !a.home.ready {
		p.Text(right.X+2, right.Y+1, "loading library", th.Muted())
		return
	}
	a.paintHomeMain(p, right, th)
}

func (a *App) paintHomeMain(p *tui.Painter, right tui.Rect, th tui.Theme) {
	lines := a.homeLines()
	if len(lines) == 0 {
		p.Text(right.X+2, right.Y+1, "nothing queued  pick a library", th.Muted())
		return
	}
	off := clampOff(homeSelLine(lines, a.home.sel), a.home.off, right.H)
	a.home.off = off
	for i := 0; i < right.H; i++ {
		paintHomeLine(p, homeLinePaint{
			right: right, lines: lines, li: off + i, rowOff: i, sel: a.home.sel, focus: a.home.focus == 1,
			hit: a.track(hitHome),
		}, th)
	}
}

func (a *App) homeLines() []homeLine {
	var lines []homeLine
	flatIdx := 0
	for _, sec := range a.home.sections {
		lines = append(lines, homeLine{header: true, text: strings.ToUpper(sec.title)})
		for _, it := range sec.items {
			lines = append(lines, homeLine{text: it.Label(), id: it.ID, idx: flatIdx})
			flatIdx++
		}
	}
	return lines
}

func homeSelLine(lines []homeLine, sel int) int {
	for i, ln := range lines {
		if !ln.header && ln.idx == sel {
			return i
		}
	}
	return 0
}

func clampOff(sel, off, vis int) int {
	if sel < off {
		off = sel
	}
	if vis > 0 && sel >= off+vis {
		off = sel - vis + 1
	}
	if off < 0 {
		return 0
	}
	return off
}

type homeLinePaint struct {
	right  tui.Rect
	lines  []homeLine
	li     int
	rowOff int
	sel    int
	focus  bool
	hit    func(tui.Rect, int)
}

func paintHomeLine(p *tui.Painter, spec homeLinePaint, th tui.Theme) {
	row := tui.Rect{X: spec.right.X + 1, Y: spec.right.Y + spec.rowOff, W: spec.right.W - 2, H: 1}
	if spec.li >= len(spec.lines) {
		p.Fill(row, th.Body())
		return
	}
	ln := spec.lines[spec.li]
	if ln.header {
		p.Fill(row, th.Body())
		p.TextClip(tui.Rect{X: row.X + 1, Y: row.Y, W: row.W - 1, H: 1}, ln.text, th.Accent())
		return
	}
	st := th.Body()
	if ln.idx == spec.sel && spec.focus {
		st = th.Focus()
	}
	p.Fill(row, st)
	if spec.hit != nil {
		spec.hit(row, ln.idx)
	}
	p.Put(row.X, row.Y, swatchCh, st.WithFg(ink(ln.id, th)))
	p.TextClip(tui.Rect{X: row.X + 2, Y: row.Y, W: row.W - 3, H: 1}, ln.text, st)
}
