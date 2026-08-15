package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

const (
	testServer = "https://jellyfin.example"
	testUser   = "alice"
	testToken  = "tok-abc"
	testUID    = "uid-1"
	testDevID  = "dev-deadbeef"
)

func TestFileSaveLoadRoundTrip(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, configFile)
	f := &File{
		Server:   testServer,
		User:     testUser,
		Token:    testToken,
		UserID:   testUID,
		DeviceID: testDevID,
		AutoNext: true,
		path:     p,
	}
	if err := f.Save(); err != nil {
		t.Fatal(err)
	}
	raw, err := os.ReadFile(p)
	if err != nil {
		t.Fatal(err)
	}
	got := &File{AutoNext: true, path: p}
	if err := json.Unmarshal(raw, got); err != nil {
		t.Fatal(err)
	}
	if got.Server != testServer || got.User != testUser || got.Token != testToken {
		t.Fatalf("round trip %+v", got)
	}
	if got.DeviceID != testDevID {
		t.Fatal("device id")
	}
}

func TestLoadCorruptJSON(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, configFile)
	if err := os.WriteFile(p, []byte("{not json"), filePerm); err != nil {
		t.Fatal(err)
	}
	root, err := os.OpenRoot(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer root.Close()
	raw, err := root.ReadFile(configFile)
	if err != nil {
		t.Fatal(err)
	}
	var f File
	if err := json.Unmarshal(raw, &f); err == nil {
		t.Fatal("want unmarshal error")
	}
}

func TestClearSession(t *testing.T) {
	f := &File{Token: testToken, UserID: testUID, ServerID: "srv"}
	f.ClearSession()
	if f.Token != "" || f.UserID != "" || f.ServerID != "" {
		t.Fatalf("session %+v", f)
	}
}

func TestNewDeviceIDLength(t *testing.T) {
	id := newDeviceID()
	if id == fallbackID {
		t.Fatal("fallback")
	}
	if len(id) != deviceBytes*2 {
		t.Fatalf("len %d", len(id))
	}
}

func TestReadConfigMissing(t *testing.T) {
	dir := t.TempDir()
	_, err := readConfig(dir)
	if !os.IsNotExist(err) {
		t.Fatalf("err %v", err)
	}
}
