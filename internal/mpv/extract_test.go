package mpv

import (
	"archive/zip"
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

const (
	zipSafeName  = "bin/mpv"
	zipEvilUp    = "../outside"
	zipEvilAbs   = "/etc/passwd"
	zipEvilDot   = "."
	zipEvilParen = ".."
	zipPayload   = "#!/bin/sh\necho mpv\n"
)

func TestZipRelTable(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{zipSafeName, zipSafeName},
		{"bin\\mpv", "bin/mpv"},
		{zipEvilUp, ""},
		{zipEvilAbs, ""},
		{zipEvilDot, ""},
		{zipEvilParen, ""},
		{"nested/../mpv", "mpv"},
		{"foo/../../bar", ""},
	}
	for _, tc := range cases {
		got := zipRel(tc.in)
		if got != tc.want {
			t.Fatalf("zipRel(%q)=%q want %q", tc.in, got, tc.want)
		}
	}
}

func TestUnzipRootSkipsTraversal(t *testing.T) {
	dir := t.TempDir()
	root, err := os.OpenRoot(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer root.Close()
	raw := makeTestZip(t, map[string]string{
		zipSafeName:     zipPayload,
		zipEvilUp:       "evil",
		zipEvilAbs:      "evil",
		"sub/nested/ok": "ok",
	})
	if err := unzipRoot(root, raw); err != nil {
		t.Fatal(err)
	}
	if _, err := root.Stat(zipSafeName); err != nil {
		t.Fatal("safe file missing")
	}
	if _, err := root.Stat("sub/nested/ok"); err != nil {
		t.Fatal("nested safe file missing")
	}
	if _, err := os.Stat(filepath.Join(dir, "outside")); !os.IsNotExist(err) {
		t.Fatal("traversal wrote outside root")
	}
}

func TestFindExeNested(t *testing.T) {
	dir := t.TempDir()
	sub := filepath.Join(dir, "deep", "bin")
	if err := os.MkdirAll(sub, dirPerm); err != nil {
		t.Fatal(err)
	}
	want := filepath.Join(sub, mpvName())
	if err := os.WriteFile(want, []byte(zipPayload), binPerm); err != nil {
		t.Fatal(err)
	}
	got := findExe(dir, mpvName())
	if got != want {
		t.Fatalf("findExe=%q want %q", got, want)
	}
}

func makeTestZip(t *testing.T, files map[string]string) []byte {
	t.Helper()
	buf := new(bytes.Buffer)
	zw := zip.NewWriter(buf)
	for name, body := range files {
		w, err := zw.Create(name)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := w.Write([]byte(body)); err != nil {
			t.Fatal(err)
		}
	}
	if err := zw.Close(); err != nil {
		t.Fatal(err)
	}
	return buf.Bytes()
}
