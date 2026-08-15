// Package app is the jellyfin-tui session: screens, input, and playback control.
package app

import (
	"context"
	"time"

	"jellyfin-tui/internal/config"
	"jellyfin-tui/internal/jellyfin"
	"jellyfin-tui/internal/mpv"
	"jellyfin-tui/internal/tui"
)

type mode uint8

const (
	modeLogin mode = iota
	modeHome
	modeBrowse
	modeSeries
	modeSearch
	modeHelp
	modeTracks
	modeCast
)

// App is the process-wide UI and playback session.
type App struct {
	term   *tui.Term
	cfg    *config.File
	jf     *jellyfin.Client
	player *mpv.Player

	ui   chan func()
	mode mode
	prev mode
	quit bool

	status string
	errMsg string

	login  loginState
	home   homeState
	browse browseState
	series seriesState
	search searchState
	tracks tracksState
	cast   castState
	dl     downloadState
	hits   []hitZone

	now       *session
	playStart time.Time
	lastProg  time.Time
	autoNext  bool
	offline   bool
}

type session struct {
	item     jellyfin.Item
	next     *jellyfin.Item
	play     jellyfin.PlayTarget
	fallback bool
}

// New builds an App on an already-opened terminal.
func New(term *tui.Term, cfg *config.File) *App {
	a := &App{
		term:     term,
		cfg:      cfg,
		ui:       make(chan func(), uiQueueSize),
		autoNext: cfg.AutoNext,
		mode:     modeLogin,
	}
	a.login.server.set(cfg.Server)
	a.login.user.set(cfg.User)
	a.initDownloads()
	if cfg.Server == "" {
		a.login.server.set(schemeHTTP)
		a.login.cur = 0
	} else if cfg.User == "" {
		a.login.cur = 1
	} else {
		a.login.cur = 2
	}
	return a
}

func (a *App) apply(fn func()) {
	select {
	case a.ui <- fn:
	default:
		go func() { a.ui <- fn }()
	}
}

// Run is the input/draw loop. It returns after quit or ctx cancel.
func (a *App) Run(ctx context.Context) error {
	if a.cfg.Token != "" && a.cfg.Server != "" {
		a.bootSession()
	}
	stop := make(chan struct{})
	defer close(stop)
	ev := a.term.Listen(stop)
	if err := a.draw(); err != nil {
		return err
	}
	tick := time.NewTicker(frameTick)
	defer tick.Stop()
	for !a.quit {
		redraw := false
		select {
		case <-ctx.Done():
			a.requestQuit()
		case e := <-ev:
			a.handle(e)
			redraw = true
		case fn := <-a.ui:
			fn()
			redraw = true
		case <-tick.C:
			if a.now != nil {
				a.maybeProgress()
				redraw = true
			}
		}
		if a.quit {
			break
		}
		if !redraw {
			continue
		}
		if err := a.draw(); err != nil {
			a.shutdown()
			return err
		}
	}
	a.shutdown()
	return nil
}

func (a *App) typing() bool {
	return a.mode == modeLogin || a.search.edit || a.browse.filtering
}
