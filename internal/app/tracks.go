package app

import (
	"strings"

	"jellyfin-tui/internal/mpv"
	"jellyfin-tui/internal/tui"
)

type tracksState struct {
	audio    []mpv.Track
	subs     []mpv.Track
	selAudio int
	selSub   int
	focus    int
}

func (a *App) openTracks() {
	if a.player == nil || a.now == nil {
		return
	}
	tracks := a.player.Tracks()
	st := tracksState{selSub: -1}
	for _, tr := range tracks {
		switch tr.Type {
		case "audio":
			st.audio = append(st.audio, tr)
			if tr.Selected {
				st.selAudio = len(st.audio) - 1
			}
		case "sub":
			st.subs = append(st.subs, tr)
			if tr.Selected {
				st.selSub = len(st.subs) - 1
			}
		}
	}
	if st.selSub < 0 {
		st.selSub = -1
	}
	a.tracks = st
	a.prev = a.mode
	a.mode = modeTracks
}

func (a *App) closeTracks() {
	a.mode = a.prev
	if a.mode == modeTracks {
		a.mode = modeHome
	}
}

func (a *App) applyTrackChoice() {
	if a.player == nil {
		return
	}
	if a.tracks.focus == 0 && a.tracks.selAudio >= 0 && a.tracks.selAudio < len(a.tracks.audio) {
		_ = a.player.SetAudio(a.tracks.audio[a.tracks.selAudio].ID)
	}
	if a.tracks.focus == 1 {
		if a.tracks.selSub < 0 {
			_ = a.player.SetSubtitleNone()
		} else if a.tracks.selSub < len(a.tracks.subs) {
			_ = a.player.SetSubtitle(a.tracks.subs[a.tracks.selSub].ID)
		}
	}
}

func (a *App) handleTracks(e tui.Event) {
	if e.Is(tui.KeyEsc) || e.RuneKey('t') {
		a.closeTracks()
		return
	}
	if e.Is(tui.KeyTab) || e.Is(tui.KeyLeft) || e.Is(tui.KeyRight) || e.RuneKey('h') || e.RuneKey('l') {
		a.tracks.focus = 1 - a.tracks.focus
		return
	}
	if a.tracks.focus == 0 {
		a.tracks.selAudio = moveSel(a.tracks.selAudio, len(a.tracks.audio), e)
	} else {
		a.tracks.selSub = moveSubSel(a.tracks.selSub, len(a.tracks.subs), e)
	}
	if e.Is(tui.KeyEnter) {
		a.applyTrackChoice()
		a.closeTracks()
	}
}

func moveSubSel(sel, n int, e tui.Event) int {
	if e.Is(tui.KeyUp) || e.RuneKey('k') {
		if sel < 0 {
			return n - 1
		}
		return sel - 1
	}
	if e.Is(tui.KeyDown) || e.RuneKey('j') {
		if sel >= n-1 {
			return -1
		}
		return sel + 1
	}
	return moveSel(sel, n, e)
}

func trackLabel(tr mpv.Track) string {
	label := tr.Title
	if label == "" {
		label = tr.Lang
	}
	if label == "" {
		label = "track"
	}
	if tr.Lang != "" && !strings.EqualFold(tr.Lang, label) {
		label += " (" + tr.Lang + ")"
	}
	return label
}

type trackColPaint struct {
	title  string
	labels []string
	sel    int
	active bool
	kind   hitKind
}

func (a *App) paintTracks(p *tui.Painter, th tui.Theme) {
	r := p.Bounds()
	w := min(tracksBoxW, r.W-4)
	h := min(tracksBoxH, r.H-2)
	box := tui.Rect{X: r.X + (r.W-w)/2, Y: r.Y + (r.H-h)/2, W: w, H: h}
	p.Fill(box, th.Panel())
	p.Frame(box, "tracks", th.Chrome().WithBg(th.Basin).WithFg(th.Ember))
	inner := box.Inset(2)
	left, right := inner.SplitX(inner.W / 2)
	audioLabels := make([]string, len(a.tracks.audio))
	for i, tr := range a.tracks.audio {
		audioLabels[i] = trackLabel(tr)
	}
	a.paintTrackCol(p, left, trackColPaint{
		title: "audio", labels: audioLabels, sel: a.tracks.selAudio,
		active: a.tracks.focus == 0, kind: hitTrackAudio,
	}, th)
	subLabels := []string{"off"}
	for _, tr := range a.tracks.subs {
		subLabels = append(subLabels, trackLabel(tr))
	}
	a.paintTrackCol(p, right, trackColPaint{
		title: "subtitles", labels: subLabels, sel: a.tracks.selSub + 1,
		active: a.tracks.focus == 1, kind: hitTrackSub,
	}, th)
}

func (a *App) paintTrackCol(p *tui.Painter, r tui.Rect, spec trackColPaint, th tui.Theme) {
	p.Text(r.X, r.Y, spec.title, th.Title().WithBg(th.Basin))
	list := tui.Rect{X: r.X, Y: r.Y + 1, W: r.W, H: r.H - 1}
	if len(spec.labels) == 0 {
		p.Text(list.X, list.Y, "none", th.Muted().WithBg(th.Basin))
		return
	}
	_ = paintRows(p, list, rowList{
		labels: spec.labels, sel: spec.sel, active: spec.active, hit: a.track(spec.kind),
	}, th)
}
