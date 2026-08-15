package app

import (
	"jellyfin-tui/internal/tui"
)

func (a *App) handle(e tui.Event) {
	if e.Kind == tui.KindResize {
		a.term.Buffer().Resize(e.W, e.H)
		_ = a.term.Resize()
		return
	}
	if a.handleMouse(e) {
		return
	}
	if a.handleSystem(e) {
		return
	}
	if a.handlePlayback(e) {
		return
	}
	a.handleMode(e)
}

func (a *App) handleSystem(e tui.Event) bool {
	if e.Is(tui.KeyCtrlC) || e.Is(tui.KeyCtrlD) {
		a.requestQuit()
		return true
	}
	if e.Is(tui.KeyCtrlL) {
		a.term.Buffer().Invalidate()
		return true
	}
	if a.mode == modeHelp {
		if e.Is(tui.KeyEsc) || e.Is(tui.KeyEnter) || e.RuneKey('?') || e.RuneKey('q') {
			a.mode = a.prev
		}
		return true
	}
	if e.RuneKey('?') && a.mode != modeLogin && !a.search.edit && !a.browse.filtering {
		a.prev = a.mode
		a.mode = modeHelp
		return true
	}
	return false
}

func (a *App) handlePlaybackSeek(e tui.Event) bool {
	switch {
	case e.RuneKey('>'):
		_ = a.player.Seek(seekShort)
	case e.RuneKey('<'):
		_ = a.player.Seek(-seekShort)
	case e.RuneKey('.'):
		_ = a.player.Seek(seekLong)
	case e.RuneKey(','):
		_ = a.player.Seek(-seekLong)
	case e.RuneKey('+') || e.RuneKey('='):
		_ = a.player.Volume(volStep)
	case e.RuneKey('-'):
		_ = a.player.Volume(-volStep)
	default:
		return false
	}
	return true
}

func (a *App) handlePlayback(e tui.Event) bool {
	if a.now == nil || a.mode == modeLogin || a.typing() || a.player == nil {
		return false
	}
	if e.RuneKey('t') {
		a.openTracks()
		return true
	}
	switch {
	case e.Is(tui.KeySpace):
		_ = a.player.PauseToggle()
	case e.RuneKey('n'):
		a.playNext()
	default:
		return a.handlePlaybackSeek(e)
	}
	return true
}

func (a *App) handleMode(e tui.Event) {
	switch a.mode {
	case modeLogin:
		a.handleLogin(e)
	case modeHome:
		a.handleHome(e)
	case modeBrowse:
		a.handleBrowse(e)
	case modeSeries:
		a.handleSeries(e)
	case modeSearch:
		a.handleSearch(e)
	case modeTracks:
		a.handleTracks(e)
	case modeCast:
		a.handleCast(e)
	}
}
