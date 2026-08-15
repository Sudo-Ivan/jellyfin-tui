package app

import (
	"fmt"
	"sync"

	"jellyfin-tui/internal/jellyfin"
	"jellyfin-tui/internal/tui"
)

type seriesState struct {
	series       jellyfin.Item
	seasons      []jellyfin.Item
	episodes     []jellyfin.Item
	allEpisodes  []jellyfin.Item
	selSeason    int
	selEp, epOff int
	seasonOff    int
	focus        int
	from         mode
}

func (a *App) openSeries(it jellyfin.Item) {
	a.series = seriesState{series: it, from: a.mode}
	a.mode = modeSeries
	a.status = statusLoading
	go func() {
		var (
			seasons []jellyfin.Item
			allEps  []jellyfin.Item
			errS    error
			errE    error
		)
		var wg sync.WaitGroup
		wg.Add(2)
		go func() {
			defer wg.Done()
			seasons, errS = a.jf.Seasons(it.ID)
		}()
		go func() {
			defer wg.Done()
			allEps, errE = a.jf.Episodes(it.ID, "")
		}()
		wg.Wait()

		a.apply(func() {
			a.status = ""
			if errS != nil {
				a.errMsg = errS.Error()
				return
			}
			if errE != nil {
				a.errMsg = errE.Error()
				return
			}
			a.series.seasons = seasons
			a.series.allEpisodes = allEps
			a.loadSeason(0)
		})
	}()
}

func (a *App) loadSeason(idx int) {
	if idx < 0 || idx >= len(a.series.seasons) {
		a.status = ""
		return
	}
	a.series.selSeason = idx
	a.series.selEp = 0
	sid := a.series.seasons[idx].ID
	var eps []jellyfin.Item
	for _, ep := range a.series.allEpisodes {
		if ep.SeasonID == sid {
			eps = append(eps, ep)
		}
	}
	a.series.episodes = eps
}

func (a *App) handleSeries(e tui.Event) {
	if a.handleSeriesMeta(e) {
		return
	}
	if a.series.focus == 0 {
		a.handleSeriesSeasons(e)
		return
	}
	if e.Is(tui.KeyLeft) || e.RuneKey('h') {
		a.series.focus = 0
		return
	}
	a.series.selEp = moveSel(a.series.selEp, len(a.series.episodes), e)
	if e.Is(tui.KeyEnter) && a.series.selEp < len(a.series.episodes) {
		a.playItem(a.series.episodes[a.series.selEp])
	}
}

func (a *App) handleSeriesMeta(e tui.Event) bool {
	if a.handleItemMeta(e) {
		return true
	}
	switch {
	case e.RuneKey('q'):
		a.requestQuit()
	case e.Is(tui.KeyEsc) || (e.RuneKey('h') && a.series.focus == 0):
		a.mode = a.series.from
		if a.mode == modeSeries {
			a.mode = modeHome
		}
	case e.RuneKey('/'):
		a.openSearch()
	case e.Is(tui.KeyTab):
		a.series.focus = 1 - a.series.focus
	default:
		return false
	}
	return true
}

func (a *App) handleSeriesSeasons(e tui.Event) {
	if e.Is(tui.KeyRight) || e.RuneKey('l') {
		a.series.focus = 1
		return
	}
	ns := moveSel(a.series.selSeason, len(a.series.seasons), e)
	if ns != a.series.selSeason {
		a.loadSeason(ns)
	}
	if e.Is(tui.KeyEnter) {
		a.series.focus = 1
	}
}

func (a *App) paintSeries(p *tui.Painter, th tui.Theme) {
	r := p.Bounds()
	head, rest := r.SplitY(2)
	title := a.series.series.Name
	if a.series.series.ProductionYear > 0 {
		title = fmt.Sprintf("%s  (%d)", title, a.series.series.ProductionYear)
	}
	p.TextClip(tui.Rect{X: head.X + 1, Y: head.Y, W: head.W - 2, H: 1}, title, th.Title())
	p.TextClip(tui.Rect{X: head.X + 1, Y: head.Y + 1, W: head.W - 2, H: 1}, clip(a.series.series.Overview, head.W-2), th.Muted())
	left, right := rest.SplitX(seasonWidth)
	p.Fill(left, th.Panel())
	p.VLine(left.MaxX()-1, left.Y, left.H, '│', th.Chrome().WithBg(th.Basin))
	slabels := make([]string, len(a.series.seasons))
	sids := make([]string, len(a.series.seasons))
	for i, s := range a.series.seasons {
		slabels[i] = s.Label()
		sids[i] = s.ID
	}
	srect := left.InsetXY(1, 1)
	a.series.seasonOff = paintRows(p.WithClip(srect), srect, rowList{
		labels: slabels, ids: sids, sel: a.series.selSeason, off: a.series.seasonOff, active: a.series.focus == 0,
		hit: a.track(hitSeason),
	}, th)

	elabels := make([]string, len(a.series.episodes))
	eids := make([]string, len(a.series.episodes))
	for i, ep := range a.series.episodes {
		name := ep.Name
		if ep.IndexNumber > 0 {
			name = fmt.Sprintf("E%d  %s", ep.IndexNumber, ep.Name)
		}
		if ep.UserData.Played {
			name = tickDone + name
		} else if ep.ResumeSeconds() > 0 {
			name = tickMid + name
		}
		elabels[i] = name
		eids[i] = ep.ID
	}
	erect := right.InsetXY(1, 1)
	a.series.epOff = paintRows(p.WithClip(erect), erect, rowList{
		labels: elabels, ids: eids, sel: a.series.selEp, off: a.series.epOff, active: a.series.focus == 1,
		hit: a.track(hitEpisode),
	}, th)
}

func clip(s string, n int) string {
	if n <= 0 {
		return ""
	}
	r := []rune(s)
	if len(r) <= n {
		return s
	}
	return string(r[:n])
}
