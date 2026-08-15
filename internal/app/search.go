package app

import (
	"jellyfin-tui/internal/jellyfin"
	"jellyfin-tui/internal/tui"
)

type searchState struct {
	query field
	edit  bool
	items []jellyfin.Item
	sel   int
	off   int
}

func (a *App) openSearch() {
	a.search = searchState{edit: true}
	a.mode = modeSearch
}

func (a *App) handleSearch(e tui.Event) {
	if a.search.edit {
		if e.Is(tui.KeyEsc) {
			a.search.edit = false
			if a.search.query.string() == "" {
				a.mode = modeHome
			}
			return
		}
		if e.Is(tui.KeyEnter) {
			a.search.edit = false
			a.runSearch()
			return
		}
		if e.Is(tui.KeyDown) {
			a.search.edit = false
			return
		}
		a.search.query.handle(e)
		return
	}
	if e.Is(tui.KeyEsc) {
		a.mode = modeHome
		return
	}
	if e.RuneKey('q') {
		a.requestQuit()
		return
	}
	if a.handleItemMeta(e) {
		return
	}
	if e.RuneKey('/') || e.Is(tui.KeyEnter) && len(a.search.items) == 0 {
		a.search.edit = true
		return
	}
	a.search.sel = moveSel(a.search.sel, len(a.search.items), e)
	if e.Is(tui.KeyEnter) && a.search.sel < len(a.search.items) {
		a.openItem(a.search.items[a.search.sel])
	}
}

func (a *App) runSearch() {
	raw := a.search.query.string()
	f := jellyfin.ParseFilter(raw)
	if f.Empty() {
		return
	}
	a.status = statusSearching
	go func() {
		items, err := a.jf.Query(f.Query(searchLimit))
		a.apply(func() {
			a.status = ""
			if err != nil {
				a.errMsg = err.Error()
				return
			}
			a.search.items = items
			a.search.sel = 0
		})
	}()
}

func (a *App) paintSearch(p *tui.Painter, th tui.Theme) {
	r := p.Bounds()
	head, rest := r.SplitY(2)
	st := th.Accent()
	if a.search.edit {
		st = th.Cursor()
	}
	p.TextClip(tui.Rect{X: head.X + 1, Y: head.Y, W: head.W - 2, H: 1}, "/"+a.search.query.string(), st)
	p.Text(head.X+1, head.Y+1, keysSearch, th.Muted())
	labels := make([]string, len(a.search.items))
	ids := make([]string, len(a.search.items))
	for i, it := range a.search.items {
		labels[i] = it.Kind() + "  " + it.Label()
		ids[i] = it.ID
	}
	list := rest.InsetXY(1, 0)
	a.search.off = paintRows(p, list, rowList{
		labels: labels, ids: ids, sel: a.search.sel, off: a.search.off, active: !a.search.edit,
		hit: a.track(hitSearch),
	}, th)
}
