package app

import (
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"jellyfin-tui/internal/tui"
)

func (a *App) draw() error {
	a.hits = a.hits[:0]
	p := a.term.Painter()
	p.Clear()
	th := p.Theme()
	w, h := p.Size()
	root := tui.Rect{W: w, H: h}
	head, rest := root.SplitY(1)
	body, foot := rest.SplitY(-a.footerH())
	a.paintHeader(p.WithClip(head), th)
	switch a.mode {
	case modeLogin:
		a.paintLogin(p.WithClip(body), th)
	case modeHome:
		a.paintHome(p.WithClip(body), th)
	case modeBrowse:
		a.paintBrowse(p.WithClip(body), th)
	case modeSeries:
		a.paintSeries(p.WithClip(body), th)
	case modeSearch:
		a.paintSearch(p.WithClip(body), th)
	case modeHelp:
		a.paintHelp(p.WithClip(body), th)
	case modeTracks:
		a.paintTracks(p.WithClip(body), th)
	case modeCast:
		a.paintCast(p.WithClip(body), th)
	}
	a.paintFooter(p.WithClip(foot), th)
	return a.term.Present()
}

func (a *App) footerH() int {
	if a.now != nil {
		return 2
	}
	return 1
}

func (a *App) paintHeader(p *tui.Painter, th tui.Theme) {
	r := p.Bounds()
	p.Fill(r, th.Panel())
	titleR := tui.Rect{X: r.X + 1, Y: r.Y, W: len(appTitle) + 1, H: 1}
	p.Text(titleR.X, titleR.Y, appTitle, th.Title().WithBg(th.Basin))
	a.hits = append(a.hits, hitZone{r: titleR, kind: hitHeader})
	host := a.headerHost()
	if host != "" {
		p.TextClip(tui.Rect{X: r.X + headerHostX, Y: r.Y, W: r.W - headerHostX - 2, H: 1}, host, th.Muted().WithBg(th.Basin))
	}
	if a.errMsg != "" {
		n := utf8.RuneCountInString(a.errMsg)
		p.TextClip(tui.Rect{X: r.MaxX() - n - 2, Y: r.Y, W: n + 1, H: 1}, a.errMsg, th.Danger().WithBg(th.Basin))
	}
}

func (a *App) headerHost() string {
	if a.jf == nil {
		return ""
	}
	host := strings.TrimPrefix(strings.TrimPrefix(a.jf.Base, schemeHTTPS), schemeHTTP)
	if a.jf.UserName != "" {
		return a.jf.UserName + "  " + host
	}
	return host
}

func (a *App) paintFooter(p *tui.Painter, th tui.Theme) {
	r := p.Bounds()
	if r.Empty() {
		return
	}
	y := r.Y
	if a.now != nil && r.H >= 2 {
		a.paintTransport(p, tui.Rect{X: r.X, Y: y, W: r.W, H: 1}, th)
		y++
	}
	keys := keysMain
	if a.mode == modeLogin {
		keys = keysLogin
	}
	if a.status != "" {
		keys = a.status + "   " + keys
	}
	if prog := a.downloadProgress(); prog != "" {
		keys = prog + "   " + keys
	}
	p.Fill(tui.Rect{X: r.X, Y: y, W: r.W, H: 1}, th.Panel())
	p.TextClip(tui.Rect{X: r.X + 1, Y: y, W: r.W - 2, H: 1}, keys, th.Muted().WithBg(th.Basin))
}

func (a *App) paintTransport(p *tui.Painter, r tui.Rect, th tui.Theme) {
	p.Fill(r, styleWell(th))
	st := a.player.Status()
	mark := markPlay
	if st.Paused {
		mark = markPause
	}
	label := a.now.item.Label()
	if a.now.play.Transcoding {
		label += " [Transcoding]"
	}
	left := fmt.Sprintf("%s  %s  %s / %s", mark, label, fmtDur(st.Pos), fmtDur(st.Dur))
	if a.now.next != nil {
		left += "   next " + a.now.next.Label()
	}
	leftR := tui.Rect{X: r.X + 1, Y: r.Y, W: r.W/2 - 1, H: 1}
	p.TextClip(leftR, left, th.Playing().WithBg(th.Well))
	a.hits = append(a.hits, hitZone{r: leftR, kind: hitPlayPause})
	gauge := tui.Rect{X: r.X + r.W/2, Y: r.Y, W: r.W/2 - 1, H: 1}
	a.hits = append(a.hits, hitZone{r: gauge, kind: hitGauge})
	frac := 0.0
	if st.Dur > 0 {
		frac = float64(st.Pos) / float64(st.Dur)
	}
	p.Gauge(gauge, frac, th.BarFill(), th.BarTrack())
}

func styleWell(th tui.Theme) tui.Style {
	return tui.Style{Fg: th.Filament, Bg: th.Well}
}

func fmtDur(d time.Duration) string {
	if d < 0 {
		d = 0
	}
	s := int(d.Seconds())
	h := s / secondsPerHour
	m := (s % secondsPerHour) / secondsPerMin
	sec := s % secondsPerMin
	if h > 0 {
		return fmt.Sprintf("%d:%02d:%02d", h, m, sec)
	}
	return fmt.Sprintf("%d:%02d", m, sec)
}
