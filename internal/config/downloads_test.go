package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

const (
	dlItemID   = "movie-1"
	dlItemName = "Heat"
	dlItemType = "Movie"
	dlItemPath = "/data/heat.mkv"
)

func TestDownloadStoreUpsertFind(t *testing.T) {
	st := &DownloadStore{}
	ent := DownloadEntry{ItemID: dlItemID, Name: dlItemName, Type: dlItemType, Path: dlItemPath}
	st.Upsert(ent)
	got, ok := st.Find(dlItemID)
	if !ok || got.Path != dlItemPath {
		t.Fatalf("find %+v ok=%v", got, ok)
	}
	ent.Path = "/data/heat2.mkv"
	st.Upsert(ent)
	got, _ = st.Find(dlItemID)
	if got.Path != "/data/heat2.mkv" {
		t.Fatal("upsert replace")
	}
	if _, ok := st.Find("missing"); ok {
		t.Fatal("missing")
	}
}

func TestDownloadStoreSaveLoad(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, downloadsFile)
	st := &DownloadStore{
		Items: []DownloadEntry{{ItemID: dlItemID, Name: dlItemName, Path: dlItemPath}},
		path:  p,
	}
	if err := st.Save(); err != nil {
		t.Fatal(err)
	}
	raw, err := os.ReadFile(p)
	if err != nil {
		t.Fatal(err)
	}
	var got DownloadStore
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatal(err)
	}
	if len(got.Items) != 1 || got.Items[0].ItemID != dlItemID {
		t.Fatalf("items %+v", got.Items)
	}
}
