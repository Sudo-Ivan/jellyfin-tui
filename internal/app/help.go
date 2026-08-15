package app

import "jellyfin-tui/internal/tui"

func (*App) paintHelp(p *tui.Painter, th tui.Theme) {
	r := p.Bounds()
	w := min(helpBoxW, r.W-4)
	h := min(helpBoxH, r.H-2)
	box := tui.Rect{X: r.X + (r.W-w)/2, Y: r.Y + (r.H-h)/2, W: w, H: h}
	p.Fill(box, th.Panel())
	p.Frame(box, "keys", th.Chrome().WithBg(th.Basin).WithFg(th.Ember))
	lines := []string{
		"j k / arrows     move",
		"enter            open or play",
		"h l / tab        switch panes",
		"/                search  genre: actor: year:",
		"m                random movie",
		"a                newly added",
		"space            pause mpv",
		"n                next episode",
		"< >              seek 10s",
		", .              seek 60s",
		"- +              volume",
		"r                reload home",
		"ctrl-l           redraw",
		"click            select, again to open",
		"wheel            scroll pane under pointer",
		"t                audio and subtitle tracks",
		"c                cast and crew",
		"d                download for offline",
		"r login          reconnect saved session",
		"q / esc          back or quit",
	}
	inner := box.Inset(2)
	for i, line := range lines {
		if i >= inner.H {
			break
		}
		p.Text(inner.X, inner.Y+i, line, th.Panel())
	}
}
