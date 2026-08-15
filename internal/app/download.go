package app

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"jellyfin-tui/internal/config"
	"jellyfin-tui/internal/jellyfin"
)

type downloadState struct {
	store  *config.DownloadStore
	dir    string
	active map[string]int
}

func (a *App) initDownloads() {
	store, err := config.LoadDownloads()
	if err != nil {
		return
	}
	dir, err := config.DownloadsDir()
	if err != nil {
		return
	}
	a.dl.store = store
	a.dl.dir = dir
	a.dl.active = map[string]int{}
}

func (a *App) localPath(itemID string) (string, bool) {
	if a.dl.store == nil {
		return "", false
	}
	ent, ok := a.dl.store.Find(itemID)
	if !ok {
		return "", false
	}
	if _, err := os.Stat(ent.Path); err != nil {
		return "", false
	}
	return ent.Path, true
}

func (a *App) startDownload() {
	it, ok := a.selectedItem()
	if !ok || a.jf == nil || a.dl.store == nil {
		return
	}
	switch it.Type {
	case jellyfin.TypeMovie, jellyfin.TypeEpisode, jellyfin.TypeVideo, jellyfin.TypeMusicVideo:
	default:
		a.errMsg = "only playable items can download"
		return
	}
	if _, ok := a.localPath(it.ID); ok {
		a.status = "already downloaded"
		return
	}
	if _, busy := a.dl.active[it.ID]; busy {
		return
	}
	a.dl.active[it.ID] = 0
	a.status = statusDownload
	go a.runDownload(it)
}

func (a *App) runDownload(it jellyfin.Item) {
	target := a.jf.OpenStream(it.ID)
	req, err := http.NewRequest(http.MethodGet, target.URL, nil)
	if err != nil {
		a.finishDownload(it.ID, err)
		return
	}
	for _, h := range a.jf.StreamHeaders() {
		name, val, ok := splitHeader(h)
		if ok {
			req.Header.Add(name, val)
		}
	}
	resp, err := a.jf.HTTP.Do(req)
	if err != nil {
		a.finishDownload(it.ID, err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		a.finishDownload(it.ID, fmt.Errorf("download http %d", resp.StatusCode))
		return
	}
	ext := filepath.Ext(resp.Request.URL.Path)
	if ext == "" {
		ext = ".bin"
	}
	path := filepath.Join(a.dl.dir, it.ID+ext)
	root, err := os.OpenRoot(a.dl.dir)
	if err != nil {
		a.finishDownload(it.ID, err)
		return
	}
	defer root.Close()
	f, err := root.Create(it.ID + ext)
	if err != nil {
		a.finishDownload(it.ID, err)
		return
	}
	total := resp.ContentLength
	written, err := copyWithProgress(f, resp.Body, total, func(pct int) {
		a.apply(func() { a.dl.active[it.ID] = pct })
	})
	_ = f.Close()
	if err != nil {
		_ = os.Remove(path)
		a.finishDownload(it.ID, err)
		return
	}
	if written == 0 {
		_ = os.Remove(path)
		a.finishDownload(it.ID, fmt.Errorf("empty download"))
		return
	}
	ent := config.DownloadEntry{ItemID: it.ID, Name: it.Label(), Type: it.Type, Path: path}
	a.dl.store.Upsert(ent)
	_ = a.dl.store.Save()
	a.apply(func() {
		delete(a.dl.active, it.ID)
		a.status = "download complete"
	})
}

func (a *App) finishDownload(id string, err error) {
	a.apply(func() {
		delete(a.dl.active, id)
		a.status = ""
		if err != nil {
			a.errMsg = err.Error()
		}
	})
}

func copyWithProgress(w io.Writer, r io.Reader, total int64, fn func(int)) (int64, error) {
	buf := make([]byte, 32*1024)
	var n int64
	last := -1
	for {
		nr, er := r.Read(buf)
		if nr > 0 {
			nw, ew := w.Write(buf[:nr])
			n += int64(nw)
			if total > 0 {
				pct := int(n * 100 / total)
				if pct != last {
					last = pct
					fn(pct)
				}
			}
			if ew != nil {
				return n, ew
			}
		}
		if er == io.EOF {
			return n, nil
		}
		if er != nil {
			return n, er
		}
	}
}

func splitHeader(h string) (name string, val string, ok bool) {
	for i := 0; i < len(h); i++ {
		if h[i] == ':' {
			name = h[:i]
			val = h[i+1:]
			if len(val) > 0 && val[0] == ' ' {
				val = val[1:]
			}
			return name, val, true
		}
	}
	return "", "", false
}

func (a *App) downloadItems() []config.DownloadEntry {
	if a.dl.store == nil {
		return nil
	}
	return a.dl.store.Items
}

func (a *App) downloadProgress() string {
	for _, pct := range a.dl.active {
		return fmt.Sprintf("%s: %d%%", statusDownload, pct)
	}
	return ""
}
