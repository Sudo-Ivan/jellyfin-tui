package app

import (
	"time"

	"jellyfin-tui/internal/jellyfin"
)

func (a *App) playItem(it jellyfin.Item) {
	if err := a.ensurePlayer(); err != nil {
		a.errMsg = err.Error()
		return
	}
	start := time.Duration(it.ResumeSeconds()) * time.Second
	target, headers, local := a.resolvePlayTarget(it)
	if target.URL == "" {
		a.errMsg = "no playable source"
		return
	}
	if err := a.player.Load(target.URL, it.Label(), start, headers); err != nil {
		if local || target.Transcoding || a.jf == nil {
			a.errMsg = err.Error()
			return
		}
		a.playTranscode(it, start)
		return
	}
	a.beginPlay(it, target, start, local)
}

func (a *App) playTranscode(it jellyfin.Item, start time.Duration) {
	if a.jf == nil {
		a.errMsg = "transcoding unavailable offline"
		return
	}
	target := a.jf.OpenTranscodeStream(it.ID)
	target.ItemID = it.ID
	a.status = statusTranscode
	if err := a.player.Load(target.URL, it.Label(), start, a.jf.StreamHeaders()); err != nil {
		a.status = ""
		a.errMsg = err.Error()
		return
	}
	a.beginPlay(it, target, start, false)
}

func (a *App) beginPlay(it jellyfin.Item, target jellyfin.PlayTarget, start time.Duration, local bool) {
	a.now = &session{item: it, play: target}
	a.playStart = time.Now()
	a.errMsg = ""
	if target.Transcoding {
		a.status = statusTranscode
	} else {
		a.status = statusPlaying
	}
	if !local && a.jf != nil {
		go func() {
			_ = a.jf.Playing(target, start, time.Duration(it.Duration())*time.Second, false)
		}()
	}
	a.prefetchNext(it)
}

func (a *App) resolvePlayTarget(it jellyfin.Item) (jellyfin.PlayTarget, []string, bool) {
	if path, ok := a.localPath(it.ID); ok {
		return jellyfin.PlayTarget{URL: path, ItemID: it.ID, MediaSourceID: it.ID}, nil, true
	}
	if a.jf == nil {
		return jellyfin.PlayTarget{}, nil, false
	}
	target := a.jf.OpenStream(it.ID)
	target.ItemID = it.ID
	return target, a.jf.StreamHeaders(), false
}

func (a *App) tryTranscodeFallback() bool {
	if a.now == nil || a.player == nil || a.jf == nil || a.now.play.Transcoding || a.now.fallback {
		return false
	}
	if time.Since(a.playStart) >= fallbackWait {
		return false
	}
	st := a.player.Status()
	if !st.EOF && !st.Idle {
		return false
	}
	it := a.now.item
	start := time.Duration(it.ResumeSeconds()) * time.Second
	a.now.fallback = true
	a.stopCurrent()
	a.playTranscode(it, start)
	return true
}
