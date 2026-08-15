package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const (
	downloadsDir  = "downloads"
	downloadsFile = "downloads.json"
)

// DownloadEntry is one saved offline item.
type DownloadEntry struct {
	ItemID string `json:"itemId"`
	Name   string `json:"name"`
	Type   string `json:"type"`
	Path   string `json:"path"`
}

// DownloadStore tracks saved offline media.
type DownloadStore struct {
	Items []DownloadEntry `json:"items"`
	path  string
}

func DownloadsDir() (string, error) {
	base, err := Dir()
	if err != nil {
		return "", err
	}
	d := filepath.Join(base, downloadsDir)
	return d, os.MkdirAll(d, dirPerm)
}

func LoadDownloads() (*DownloadStore, error) {
	dir, err := Dir()
	if err != nil {
		return nil, err
	}
	p := filepath.Join(dir, downloadsFile)
	st := &DownloadStore{path: p}
	root, err := os.OpenRoot(dir)
	if err != nil {
		return st, nil
	}
	defer root.Close()
	raw, err := root.ReadFile(downloadsFile)
	if err != nil {
		if os.IsNotExist(err) {
			return st, nil
		}
		return nil, err
	}
	if err := json.Unmarshal(raw, st); err != nil {
		return nil, err
	}
	st.path = p
	return st, nil
}

func (s *DownloadStore) Save() error {
	if s.path == "" {
		dir, err := Dir()
		if err != nil {
			return err
		}
		s.path = filepath.Join(dir, downloadsFile)
	}
	dir := filepath.Dir(s.path)
	if err := os.MkdirAll(dir, dirPerm); err != nil {
		return err
	}
	raw, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	root, err := os.OpenRoot(dir)
	if err != nil {
		return err
	}
	defer root.Close()
	return root.WriteFile(downloadsFile, raw, filePerm)
}

func (s *DownloadStore) Find(itemID string) (DownloadEntry, bool) {
	for _, it := range s.Items {
		if it.ItemID == itemID {
			return it, true
		}
	}
	return DownloadEntry{}, false
}

func (s *DownloadStore) Upsert(e DownloadEntry) {
	for i, it := range s.Items {
		if it.ItemID == e.ItemID {
			s.Items[i] = e
			return
		}
	}
	s.Items = append(s.Items, e)
}
