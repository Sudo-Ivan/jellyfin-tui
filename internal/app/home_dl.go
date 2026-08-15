package app

import "jellyfin-tui/internal/jellyfin"

func (a *App) showDownloads() {
	items := make([]jellyfin.Item, 0, len(a.downloadItems()))
	for _, d := range a.downloadItems() {
		items = append(items, jellyfin.Item{ID: d.ItemID, Name: d.Name, Type: d.Type})
	}
	a.home.sections = []section{{"downloads", items}}
	a.rebuildHomeFlat()
	a.home.sel = 0
}
