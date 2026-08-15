package app

import (
	"time"

	"jellyfin-tui/internal/jellyfin"
	"jellyfin-tui/internal/mpv"
)

func (a *App) openItem(it jellyfin.Item) {
	switch it.Type {
	case jellyfin.TypeSeries:
		a.openSeries(it)
	case jellyfin.TypeSeason:
		if it.SeriesID != "" {
			a.openSeries(jellyfin.Item{ID: it.SeriesID, Name: it.SeriesName, Type: jellyfin.TypeSeries})
		} else {
			a.openBrowse(it)
		}
	case jellyfin.TypeFolder, jellyfin.TypeCollection, jellyfin.TypePlaylist, jellyfin.TypeBoxSet, jellyfin.TypeHome:
		a.openBrowse(it)
	case jellyfin.TypeMovie, jellyfin.TypeEpisode, jellyfin.TypeVideo, jellyfin.TypeMusicVideo:
		a.playItem(it)
	default:
		if it.CollectionType != "" {
			a.openBrowse(it)
		} else {
			a.playItem(it)
		}
	}
}

func (a *App) ensurePlayer() error {
	if a.player != nil && !a.player.IsClosed() {
		return nil
	}
	if a.player != nil {
		_ = a.player.Close()
		a.player = nil
	}
	p, err := mpv.Start(func() {
		a.apply(func() { a.handleEOF() })
	})
	if err != nil {
		return err
	}
	a.player = p
	return nil
}

func (a *App) playRandomMovie() {
	if a.jf == nil && !a.offline {
		return
	}
	parent := ""
	if a.mode == modeBrowse && a.browse.parent.CollectionType == jellyfin.CollectionMovies {
		parent = a.browse.parent.ID
	}
	a.status = statusRandom
	go func() {
		it, err := a.jf.RandomMovie(parent)
		a.apply(func() {
			a.status = ""
			if err != nil {
				a.errMsg = err.Error()
				return
			}
			a.playItem(it)
		})
	}()
}

func (a *App) prefetchNext(it jellyfin.Item) {
	if it.Type != jellyfin.TypeEpisode || a.jf == nil {
		a.now.next = nil
		return
	}
	seriesID := it.SeriesID
	if seriesID == "" {
		return
	}
	cur := it
	go func() {
		list, err := a.jf.AdjacentEpisodes(seriesID, cur.ID)
		if err != nil || len(list) == 0 {
			list, err = a.jf.Episodes(seriesID, "")
		}
		if err != nil {
			return
		}
		next, ok := jellyfin.NextEpisode(list, cur.ID)
		if !ok {
			return
		}
		a.apply(func() {
			if a.now != nil && a.now.item.ID == cur.ID {
				n := next
				a.now.next = &n
			}
		})
	}()
}

func (a *App) playNext() {
	if a.now == nil || a.now.next == nil {
		a.status = statusNoNext
		return
	}
	a.stopCurrent()
	a.playItem(*a.now.next)
}

func (a *App) handleEOF() {
	if a.now == nil || a.player == nil {
		return
	}
	if a.tryTranscodeFallback() {
		return
	}
	st := a.player.Status()
	if !st.EOF && !st.Idle {
		return
	}
	next := a.now.next
	a.stopCurrent()
	if a.player.IsClosed() {
		_ = a.player.Close()
		a.player = nil
	}
	if a.autoNext && next != nil {
		a.playItem(*next)
		return
	}
	a.now = nil
	a.status = ""
}

func (a *App) stopCurrent() {
	if a.now == nil || a.player == nil {
		return
	}
	st := a.player.Status()
	id := a.now.item.ID
	pos, dur := st.Pos, st.Dur
	t := a.now.play
	t.ItemID = id

	it := a.now.item
	if st.EOF {
		it.UserData.Played = true
		it.UserData.PlaybackPositionTicks = 0
	} else {
		it.UserData.PlaybackPositionTicks = int64(pos.Nanoseconds() / 100)
		if dur > 0 && float64(pos)/float64(dur) > playedThreshold {
			it.UserData.Played = true
			it.UserData.PlaybackPositionTicks = 0
		}
	}
	a.updateItemState(it)

	if a.jf != nil {
		go func() {
			_ = a.jf.Stopped(t, pos, dur)
			updated, err := a.jf.Item(id)
			if err == nil {
				a.apply(func() {
					a.updateItemState(updated)
				})
			}
		}()
	}
}

func (a *App) updateItemState(it jellyfin.Item) {
	for sIdx, s := range a.home.sections {
		for i, item := range s.items {
			if item.ID == it.ID {
				a.home.sections[sIdx].items[i] = it
			}
		}
	}
	a.rebuildHomeFlat()

	for i, item := range a.browse.items {
		if item.ID == it.ID {
			a.browse.items[i] = it
		}
	}
	a.browse.filterShown()

	for i, item := range a.series.episodes {
		if item.ID == it.ID {
			a.series.episodes[i] = it
		}
	}
	for i, item := range a.series.allEpisodes {
		if item.ID == it.ID {
			a.series.allEpisodes[i] = it
		}
	}
}

func (a *App) maybeProgress() {
	if a.now == nil || a.jf == nil || a.player == nil {
		return
	}
	if time.Since(a.lastProg) < progressEvery {
		return
	}
	a.lastProg = time.Now()
	st := a.player.Status()
	t := a.now.play
	t.ItemID = a.now.item.ID
	go func() {
		_ = a.jf.Progress(t, st.Pos, st.Dur, st.Paused)
	}()
}
