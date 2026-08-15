package app

import (
	"jellyfin-tui/internal/jellyfin"
	"jellyfin-tui/internal/tui"
)

func (a *App) selectedItem() (jellyfin.Item, bool) {
	switch a.mode {
	case modeHome:
		if a.home.sel >= 0 && a.home.sel < len(a.home.flat) {
			return a.home.flat[a.home.sel].item, true
		}
	case modeBrowse:
		return a.browse.current()
	case modeSeries:
		if a.series.selEp >= 0 && a.series.selEp < len(a.series.episodes) {
			return a.series.episodes[a.series.selEp], true
		}
	case modeSearch:
		if !a.search.edit && a.search.sel >= 0 && a.search.sel < len(a.search.items) {
			return a.search.items[a.search.sel], true
		}
	}
	return jellyfin.Item{}, false
}

func (a *App) handleItemMeta(e tui.Event) bool {
	if e.RuneKey('c') {
		a.openCast()
		return true
	}
	if e.RuneKey('d') {
		a.startDownload()
		return true
	}
	return false
}
