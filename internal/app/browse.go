package app

import (
	"jellyfin-tui/internal/jellyfin"
	"jellyfin-tui/internal/tui"
)

type browseState struct {
	parent    jellyfin.Item
	items     []jellyfin.Item
	sel, off  int
	filter    field
	filtering bool
	shown     []int
}

func (a *App) openBrowse(it jellyfin.Item) {
	a.browse = browseState{parent: it}
	a.mode = modeBrowse
	a.status = statusLoading
	go func() {
		include := ""
		switch it.CollectionType {
		case jellyfin.CollectionMovies:
			include = jellyfin.TypeMovie
		case jellyfin.CollectionTV:
			include = jellyfin.TypeSeries
		}
		items, err := a.jf.Children(it.ID, include)
		a.apply(func() {
			a.status = ""
			if err != nil {
				a.errMsg = err.Error()
				return
			}
			a.browse.items = items
			a.browse.filterShown()
		})
	}()
}

func (b *browseState) filterShown() {
	f := jellyfin.ParseFilter(b.filter.string())
	b.shown = b.shown[:0]
	for i, it := range b.items {
		if it.Matches(f) {
			b.shown = append(b.shown, i)
		}
	}
	if b.sel >= len(b.shown) {
		b.sel = max(0, len(b.shown)-1)
	}
}

func (b *browseState) current() (jellyfin.Item, bool) {
	if b.sel < 0 || b.sel >= len(b.shown) {
		return jellyfin.Item{}, false
	}
	return b.items[b.shown[b.sel]], true
}

func (a *App) handleBrowse(e tui.Event) {
	if a.browse.filtering {
		if e.Is(tui.KeyEsc) || e.Is(tui.KeyEnter) {
			a.browse.filtering = false
			return
		}
		a.browse.filter.handle(e)
		a.browse.filterShown()
		return
	}
	if e.RuneKey('/') {
		a.browse.filtering = true
		return
	}
	if e.RuneKey('m') {
		a.playRandomMovie()
		return
	}
	if a.handleItemMeta(e) {
		return
	}
	if e.RuneKey('q') {
		a.requestQuit()
		return
	}
	if e.Is(tui.KeyEsc) || e.RuneKey('h') || e.Is(tui.KeyLeft) {
		a.mode = modeHome
		return
	}
	a.browse.sel = moveSel(a.browse.sel, len(a.browse.shown), e)
	if e.Is(tui.KeyEnter) || e.RuneKey('l') || e.Is(tui.KeyRight) {
		if it, ok := a.browse.current(); ok {
			a.openItem(it)
		}
	}
}

func (a *App) paintBrowse(p *tui.Painter, th tui.Theme) {
	r := p.Bounds()
	title, rest := r.SplitY(1)
	p.TextClip(title.InsetXY(1, 0), a.browse.parent.Name, th.Title())
	if a.browse.filtering || a.browse.filter.string() != "" {
		filter, rest2 := rest.SplitY(1)
		rest = rest2
		label := "/" + a.browse.filter.string()
		st := th.Accent()
		if a.browse.filtering {
			st = th.Cursor()
		}
		p.TextClip(filter.InsetXY(1, 0), label, st)
	}
	labels := make([]string, len(a.browse.shown))
	ids := make([]string, len(a.browse.shown))
	for i, idx := range a.browse.shown {
		it := a.browse.items[idx]
		labels[i] = it.Label()
		ids[i] = it.ID
	}
	list := rest.InsetXY(1, 0)
	a.browse.off = paintRows(p, list, rowList{
		labels: labels, ids: ids, sel: a.browse.sel, off: a.browse.off, active: true,
		hit: a.track(hitBrowse),
	}, th)
}

func (a *App) openNewlyAdded() {
	a.browse = browseState{parent: jellyfin.Item{ID: idAdded, Name: nameAdded, Type: jellyfin.TypeHome}}
	a.mode = modeBrowse
	a.status = statusLoading
	go func() {
		items, err := a.jf.NewlyAdded(newlyLimit)
		a.apply(func() {
			a.status = ""
			if err != nil {
				a.errMsg = err.Error()
				return
			}
			a.browse.items = items
			a.browse.filterShown()
		})
	}()
}
