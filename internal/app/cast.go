package app

import (
	"jellyfin-tui/internal/jellyfin"
	"jellyfin-tui/internal/tui"
)

type castState struct {
	item   jellyfin.Item
	people []jellyfin.Person
	sel    int
	off    int
}

func (a *App) openCast() {
	it, ok := a.selectedItem()
	if !ok || a.jf == nil {
		return
	}
	if len(it.People) == 0 {
		go func() {
			full, err := a.jf.Item(it.ID)
			a.apply(func() {
				if err != nil {
					a.errMsg = err.Error()
					return
				}
				a.showCast(full)
			})
		}()
		return
	}
	a.showCast(it)
}

func (a *App) showCast(it jellyfin.Item) {
	people := filterCast(it.People)
	if len(people) == 0 {
		a.errMsg = "no cast listed"
		return
	}
	a.cast = castState{item: it, people: people}
	a.prev = a.mode
	a.mode = modeCast
}

func filterCast(list []jellyfin.Person) []jellyfin.Person {
	out := make([]jellyfin.Person, 0, len(list))
	for _, p := range list {
		if p.Name == "" {
			continue
		}
		out = append(out, p)
	}
	return out
}

func (a *App) closeCast() {
	a.mode = a.prev
	if a.mode == modeCast {
		a.mode = modeHome
	}
}

func (a *App) handleCast(e tui.Event) {
	if e.Is(tui.KeyEsc) || e.RuneKey('c') {
		a.closeCast()
		return
	}
	a.cast.sel = moveSel(a.cast.sel, len(a.cast.people), e)
	if e.Is(tui.KeyEnter) && a.cast.sel < len(a.cast.people) {
		a.searchPerson(a.cast.people[a.cast.sel].Name)
	}
}

func (a *App) searchPerson(name string) {
	if a.jf == nil {
		return
	}
	a.closeCast()
	a.search = searchState{}
	a.search.query.set(`actor:"` + name + `"`)
	a.mode = modeSearch
	a.runSearch()
}

func (a *App) paintCast(p *tui.Painter, th tui.Theme) {
	r := p.Bounds()
	w := min(castBoxW, r.W-4)
	h := min(castBoxH, r.H-2)
	box := tui.Rect{X: r.X + (r.W-w)/2, Y: r.Y + (r.H-h)/2, W: w, H: h}
	p.Fill(box, th.Panel())
	title := "cast  " + a.cast.item.Name
	p.Frame(box, title, th.Chrome().WithBg(th.Basin).WithFg(th.Ember))
	labels := make([]string, len(a.cast.people))
	for i, person := range a.cast.people {
		role := person.Role
		if role == "" {
			role = person.Type
		}
		if role != "" {
			labels[i] = person.Name + "  (" + role + ")"
		} else {
			labels[i] = person.Name
		}
	}
	inner := box.Inset(2)
	a.cast.off = paintRows(p, inner, rowList{
		labels: labels, sel: a.cast.sel, active: true, hit: a.track(hitCast),
	}, th)
}
