package app

import (
	"sync"

	"jellyfin-tui/internal/jellyfin"
	"jellyfin-tui/internal/tui"
)

type homeState struct {
	rail     []jellyfin.Item
	railSel  int
	railOff  int
	focus    int
	sections []section
	flat     []homeHit
	sel      int
	off      int
	ready    bool
}

type section struct {
	title string
	items []jellyfin.Item
}

type homeHit struct {
	sec  int
	item jellyfin.Item
}

func (a *App) reloadHome() {
	if a.jf == nil {
		return
	}
	a.jf.Invalidate()
	a.status = statusLoading
	a.home.ready = false
	go func() {
		var (
			views  []jellyfin.Item
			resume []jellyfin.Item
			next   []jellyfin.Item
			latest []jellyfin.Item
			err1   error
		)
		var wg sync.WaitGroup
		wg.Add(4)
		go func() {
			defer wg.Done()
			views, err1 = a.jf.Views()
		}()
		go func() {
			defer wg.Done()
			resume, _ = a.jf.Resume(homeResumeN)
		}()
		go func() {
			defer wg.Done()
			next, _ = a.jf.NextUp(homeNextUpN)
		}()
		go func() {
			defer wg.Done()
			latest, _ = a.jf.NewlyAdded(homeLatestN)
		}()
		wg.Wait()

		a.apply(func() {
			a.applyHome(views, resume, next, latest, err1)
		})
	}()
}

func (a *App) applyHome(views, resume, next, latest []jellyfin.Item, err error) {
	a.status = ""
	if err != nil {
		a.errMsg = err.Error()
	}
	home := jellyfin.Item{ID: idHome, Name: nameHome, Type: jellyfin.TypeHome}
	a.home.rail = append([]jellyfin.Item{home}, views...)
	if len(a.downloadItems()) > 0 {
		a.home.rail = append(a.home.rail, jellyfin.Item{
			ID: idDownloads, Name: nameDownloads, Type: jellyfin.TypeHome,
		})
	}
	a.home.sections = nil
	if len(resume) > 0 {
		a.home.sections = append(a.home.sections, section{"continue", resume})
	}
	if len(next) > 0 {
		a.home.sections = append(a.home.sections, section{"next up", next})
	}
	if len(latest) > 0 {
		a.home.sections = append(a.home.sections, section{nameAdded, latest})
	}
	a.rebuildHomeFlat()
	a.home.ready = true
}

func (a *App) rebuildHomeFlat() {
	a.home.flat = nil
	for i, s := range a.home.sections {
		for _, it := range s.items {
			a.home.flat = append(a.home.flat, homeHit{sec: i, item: it})
		}
	}
	if a.home.sel >= len(a.home.flat) {
		a.home.sel = max(0, len(a.home.flat)-1)
	}
}

func (a *App) handleHome(e tui.Event) {
	if a.handleHomeMeta(e) {
		return
	}
	if a.home.focus == 0 {
		a.handleHomeRail(e)
		return
	}
	a.home.sel = moveSel(a.home.sel, len(a.home.flat), e)
	if e.Is(tui.KeyEnter) && a.home.sel < len(a.home.flat) {
		a.openItem(a.home.flat[a.home.sel].item)
	}
}

func (a *App) handleHomeMeta(e tui.Event) bool {
	if a.handleItemMeta(e) {
		return true
	}
	switch {
	case e.RuneKey('q'):
		a.requestQuit()
	case e.RuneKey('/'):
		a.openSearch()
	case e.RuneKey('m'):
		a.playRandomMovie()
	case e.RuneKey('a'):
		a.openNewlyAdded()
	case e.RuneKey('r'):
		a.reloadHome()
	case e.Is(tui.KeyTab):
		a.home.focus = 1 - a.home.focus
	case e.Is(tui.KeyLeft) || e.RuneKey('h'):
		a.home.focus = 0
	case e.Is(tui.KeyRight) || e.RuneKey('l'):
		a.handleHomeRight()
	default:
		return false
	}
	return true
}

func (a *App) handleHomeRight() {
	if a.home.focus == 0 && a.home.railSel > 0 && a.home.railSel < len(a.home.rail) {
		a.openBrowse(a.home.rail[a.home.railSel])
		return
	}
	a.home.focus = 1
}

func (a *App) handleHomeRail(e tui.Event) {
	a.home.railSel = moveSel(a.home.railSel, len(a.home.rail), e)
	if !e.Is(tui.KeyEnter) || a.home.railSel < 0 || a.home.railSel >= len(a.home.rail) {
		return
	}
	it := a.home.rail[a.home.railSel]
	if it.ID == idDownloads {
		a.showDownloads()
		a.home.focus = 1
		return
	}
	if it.Type == jellyfin.TypeHome {
		a.home.focus = 1
		return
	}
	a.openBrowse(it)
}
