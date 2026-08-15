package app

import "jellyfin-tui/internal/tui"

type hitKind uint8

const (
	hitNone hitKind = iota
	hitRail
	hitHome
	hitBrowse
	hitSearch
	hitSeason
	hitEpisode
)

type hitZone struct {
	r    tui.Rect
	kind hitKind
	idx  int
}

func pickHit(zones []hitZone, x, y int) (hitZone, bool) {
	for i := len(zones) - 1; i >= 0; i-- {
		if zones[i].r.Contains(x, y) {
			return zones[i], true
		}
	}
	return hitZone{}, false
}

func (a *App) track(kind hitKind) func(tui.Rect, int) {
	return func(r tui.Rect, idx int) {
		a.hits = append(a.hits, hitZone{r: r, kind: kind, idx: idx})
	}
}

func (a *App) handleMouseOverlay(e tui.Event) bool {
	if a.mode != modeTracks && a.mode != modeCast {
		return false
	}
	if !e.Click() {
		return true
	}
	z, ok := pickHit(a.hits, e.X, e.Y)
	if ok {
		a.clickZone(z, e)
		return true
	}
	if a.mode == modeTracks {
		a.closeTracks()
	} else {
		a.closeCast()
	}
	return true
}

func (a *App) handleMouse(e tui.Event) bool {
	if e.Kind != tui.KindMouse {
		return false
	}
	if a.mode == modeHelp {
		if e.Click() {
			a.mode = a.prev
		}
		return true
	}
	if a.handleMouseOverlay(e) {
		return true
	}
	if w := e.Wheel(); w != 0 {
		a.wheel(e, w)
		return true
	}
	if !e.Click() || (a.typing() && a.mode != modeLogin) {
		return e.Click()
	}
	z, ok := pickHit(a.hits, e.X, e.Y)
	if !ok {
		return true
	}
	a.clickZone(z, e)
	return true
}

func (a *App) wheel(e tui.Event, delta int) {
	z, ok := pickHit(a.hits, e.X, e.Y)
	if ok {
		a.focusHit(z.kind)
	}
	step := tui.Event{Kind: tui.KindKey, Key: tui.KeyDown}
	if delta < 0 {
		step.Key = tui.KeyUp
	}
	a.handleMode(step)
}

func (a *App) focusHit(k hitKind) {
	switch k {
	case hitRail, hitSeason:
		a.setPane(0)
	case hitHome, hitBrowse, hitSearch, hitEpisode, hitCast:
		a.setPane(1)
	case hitTrackAudio:
		if a.mode == modeTracks {
			a.tracks.focus = 0
		}
	case hitTrackSub:
		if a.mode == modeTracks {
			a.tracks.focus = 1
		}
	}
}

func (a *App) setPane(n int) {
	switch a.mode {
	case modeHome:
		a.home.focus = n
	case modeSeries:
		a.series.focus = n
	}
}

func (a *App) clickZone(z hitZone, e tui.Event) {
	switch z.kind {
	case hitRail:
		a.clickRail(z.idx)
	case hitHome:
		a.clickHome(z.idx)
	case hitBrowse:
		a.clickBrowse(z.idx)
	case hitSearch:
		a.clickSearch(z.idx)
	case hitSeason:
		a.clickSeason(z.idx)
	case hitEpisode:
		a.clickEpisode(z.idx)
	default:
		a.clickExtra(z, e)
	}
}

func (a *App) clickRail(i int) {
	a.home.focus = 0
	if i == a.home.railSel {
		a.handleHomeRail(tui.Event{Kind: tui.KindKey, Key: tui.KeyEnter})
		return
	}
	a.home.railSel = i
}

func (a *App) clickHome(i int) {
	a.home.focus = 1
	if i == a.home.sel {
		if i < len(a.home.flat) {
			a.openItem(a.home.flat[i].item)
		}
		return
	}
	a.home.sel = i
}

func (a *App) clickBrowse(i int) {
	if i == a.browse.sel {
		if it, ok := a.browse.current(); ok {
			a.openItem(it)
		}
		return
	}
	a.browse.sel = i
}

func (a *App) clickSearch(i int) {
	a.search.edit = false
	if i == a.search.sel {
		if i < len(a.search.items) {
			a.openItem(a.search.items[i])
		}
		return
	}
	a.search.sel = i
}

func (a *App) clickSeason(i int) {
	a.series.focus = 0
	if i != a.series.selSeason {
		a.loadSeason(i)
	}
	a.series.focus = 1
}

func (a *App) clickEpisode(i int) {
	a.series.focus = 1
	if i == a.series.selEp {
		if i < len(a.series.episodes) {
			a.playItem(a.series.episodes[i])
		}
		return
	}
	a.series.selEp = i
}
