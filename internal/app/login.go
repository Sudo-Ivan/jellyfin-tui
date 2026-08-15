package app

import (
	"strings"

	"jellyfin-tui/internal/config"
	"jellyfin-tui/internal/jellyfin"
	"jellyfin-tui/internal/tui"
)

type loginState struct {
	server, user, pass field
	cur                int
	busy               bool
}

func (a *App) tryLoginReconnect(e tui.Event) bool {
	if !e.RuneKey('r') || a.cfg.Token == "" || a.cfg.Server == "" {
		return false
	}
	a.bootSession()
	return true
}

func (a *App) handleLoginField(e tui.Event) {
	switch a.login.cur {
	case 0:
		a.login.server.handle(e)
	case 1:
		a.login.user.handle(e)
	case 2:
		a.login.pass.handle(e)
	}
}

func (a *App) handleLogin(e tui.Event) {
	if e.Is(tui.KeyEsc) {
		a.requestQuit()
		return
	}
	if a.login.busy {
		return
	}
	if a.tryLoginReconnect(e) {
		return
	}
	switch {
	case e.Is(tui.KeyTab) || e.Is(tui.KeyDown):
		a.login.cur = (a.login.cur + 1) % loginFields
	case e.Is(tui.KeyBacktab) || e.Is(tui.KeyUp):
		a.login.cur = (a.login.cur + 2) % loginFields
	case e.Is(tui.KeyEnter):
		a.submitLogin()
	case e.RuneKey('q') && a.login.cur == 0 && len(a.login.server.runes) == 0:
		a.requestQuit()
	default:
		a.handleLoginField(e)
	}
}

func (a *App) submitLogin() {
	server := strings.TrimSpace(a.login.server.string())
	user := strings.TrimSpace(a.login.user.string())
	pass := a.login.pass.string()
	if server == "" || server == schemeHTTP || server == schemeHTTPS {
		a.errMsg = errNeedServer
		return
	}
	if !strings.Contains(server, "://") {
		server = schemeHTTP + server
	}
	server = strings.TrimRight(server, "/")
	if user == "" {
		a.errMsg = errNeedUser
		return
	}
	a.login.busy = true
	a.status = statusSigningIn
	a.errMsg = ""
	go func() {
		c := jellyfin.New(server, config.DeviceName(), a.cfg.DeviceID)
		err := c.Login(user, pass)
		a.apply(func() {
			a.login.busy = false
			a.status = ""
			if err != nil {
				a.errMsg = err.Error()
				return
			}
			a.jf = c
			a.cfg.Server = server
			a.cfg.User = user
			a.cfg.Token = c.Token
			a.cfg.UserID = c.UserID
			a.cfg.ServerID = c.ServerID
			_ = a.cfg.Save()
			a.login.pass.set("")
			a.mode = modeHome
			a.reloadHome()
		})
	}()
}

func (a *App) paintLogin(p *tui.Painter, th tui.Theme) {
	r := p.Bounds()
	boxW := min(loginBoxW, r.W-4)
	boxH := loginBoxH
	x := r.X + (r.W-boxW)/2
	y := r.Y + max(1, (r.H-boxH)/2)
	box := tui.Rect{X: x, Y: y, W: boxW, H: boxH}
	p.Fill(box, th.Panel())
	p.Frame(box, "connect", th.Chrome().WithBg(th.Basin).WithFg(th.Ember))
	inner := box.Inset(2)
	p.Text(inner.X, inner.Y, "oxide terminal for jellyfin", th.Muted().WithBg(th.Basin))
	row := tui.Rect{X: inner.X, Y: inner.Y + 2, W: inner.W, H: 1}
	paintField(p, row, fieldPaint{label: "server", value: a.login.server.string(), cur: a.login.server.cur, focus: a.login.cur == 0}, th)
	a.hits = append(a.hits, hitZone{r: row, kind: hitLoginField, idx: 0})
	row.Y += loginFieldGap
	paintField(p, row, fieldPaint{label: "user", value: a.login.user.string(), cur: a.login.user.cur, focus: a.login.cur == 1}, th)
	a.hits = append(a.hits, hitZone{r: row, kind: hitLoginField, idx: 1})
	row.Y += loginFieldGap
	paintField(p, row, fieldPaint{
		label: "password", value: a.login.pass.string(), cur: a.login.pass.cur,
		secret: true, focus: a.login.cur == 2,
	}, th)
	a.hits = append(a.hits, hitZone{r: row, kind: hitLoginField, idx: 2})
	hint := "enter to sign in"
	if a.cfg.Token != "" && a.cfg.Server != "" {
		hintRow := tui.Rect{X: inner.X, Y: inner.MaxY() - 2, W: inner.W, H: 1}
		p.Text(hintRow.X, hintRow.Y, "r reconnect with saved session", th.Accent().WithBg(th.Basin))
		a.hits = append(a.hits, hitZone{r: hintRow, kind: hitReconnect})
	}
	if a.login.busy {
		hint = "connecting"
	}
	p.Text(inner.X, inner.MaxY()-1, hint, th.Accent().WithBg(th.Basin))
}
