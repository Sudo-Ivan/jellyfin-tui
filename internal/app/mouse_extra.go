package app

import (
	"time"

	"jellyfin-tui/internal/tui"
)

const (
	hitHeader hitKind = iota + 20
	hitPlayPause
	hitGauge
	hitLoginField
	hitReconnect
	hitTrackAudio
	hitTrackSub
	hitCast
)

func (a *App) clickExtra(z hitZone, e tui.Event) {
	switch z.kind {
	case hitHeader:
		a.mode = modeHome
	case hitPlayPause:
		if a.player != nil {
			_ = a.player.PauseToggle()
		}
	case hitGauge:
		a.clickGauge(z.r, e.X)
	case hitLoginField:
		a.login.cur = z.idx
	case hitReconnect:
		a.bootSession()
	case hitTrackAudio:
		a.tracks.focus = 0
		a.tracks.selAudio = z.idx
		a.applyTrackChoice()
		a.closeTracks()
	case hitTrackSub:
		a.tracks.focus = 1
		a.tracks.selSub = z.idx - 1
		a.applyTrackChoice()
		a.closeTracks()
	case hitCast:
		a.clickCast(z.idx)
	}
}

func (a *App) clickGauge(r tui.Rect, x int) {
	if a.player == nil || a.now == nil {
		return
	}
	st := a.player.Status()
	if st.Dur <= 0 || r.W <= 0 {
		return
	}
	rel := max(x-r.X, 0)
	if rel >= r.W {
		rel = r.W - 1
	}
	frac := float64(rel) / float64(r.W)
	pos := time.Duration(float64(st.Dur) * frac)
	_ = a.player.SeekAbs(pos)
}

func (a *App) clickCast(idx int) {
	if idx < 0 || idx >= len(a.cast.people) {
		return
	}
	a.cast.sel = idx
	a.searchPerson(a.cast.people[idx].Name)
}
