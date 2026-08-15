package app

import (
	"time"

	"jellyfin-tui/internal/config"
	"jellyfin-tui/internal/jellyfin"
)

func (a *App) requestQuit() {
	a.quit = true
}

func (a *App) shutdown() {
	a.stopCurrent()
	if a.player != nil {
		_ = a.player.Close()
		a.player = nil
	}
	a.now = nil
}

func (a *App) bootSession() {
	a.status = statusReconnect
	a.login.busy = true
	go func() {
		c := jellyfin.New(a.cfg.Server, config.DeviceName(), a.cfg.DeviceID)
		c.Token = a.cfg.Token
		c.UserID = a.cfg.UserID
		err := probeMe(c)
		a.apply(func() {
			a.login.busy = false
			if err != nil {
				a.status = ""
				if jellyfin.IsAuth(err) {
					a.errMsg = errExpired
					a.cfg.ClearSession()
					_ = a.cfg.Save()
					return
				}
				a.enterOffline(err)
				return
			}
			a.offline = false
			a.jf = c
			a.mode = modeHome
			a.status = ""
			a.errMsg = ""
			a.reloadHome()
		})
	}()
}

func (a *App) enterOffline(err error) {
	a.offline = true
	a.jf = nil
	a.mode = modeHome
	a.status = statusOffline
	if err != nil {
		a.errMsg = err.Error()
	}
	a.setupOfflineHome()
}

func (a *App) setupOfflineHome() {
	home := jellyfin.Item{ID: idHome, Name: nameHome, Type: jellyfin.TypeHome}
	dl := jellyfin.Item{ID: idDownloads, Name: nameDownloads, Type: jellyfin.TypeHome}
	a.home.rail = []jellyfin.Item{home, dl}
	a.showDownloads()
	a.home.ready = true
}

func probeMe(c *jellyfin.Client) error {
	var last error
	for i := range loginTries {
		if i > 0 {
			time.Sleep(reconnectWait * time.Duration(1<<(i-1)))
		}
		_, last = c.Me()
		if last == nil || jellyfin.IsAuth(last) {
			return last
		}
	}
	return last
}
