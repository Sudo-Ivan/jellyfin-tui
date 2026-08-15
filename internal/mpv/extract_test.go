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

func TestZipSymlinkTargetTable(t *testing.T) {
	cases := []struct {
		name, in, want string
	}{
		{sharedDir, libDir, libDir},
		{sharedDir, zipEvilAbs, ""},
		{zipSafeName, "../../x", ""},
		{"bin/a", "../" + libDir, "../" + libDir},
		{"link", zipEvilUp, ""},
		{sharedDir, "", ""},
	}
	for _, tc := range cases {
		got := zipSymlinkTarget(tc.name, tc.in)
		if got != tc.want {
			t.Fatalf("zipSymlinkTarget(%q,%q)=%q want %q", tc.name, tc.in, got, tc.want)
		}
	}
}

func TestUnzipSymlink(t *testing.T) {
	dir := t.TempDir()
	root, err := os.OpenRoot(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer root.Close()
	raw := makeLinkZip(t, map[string]string{"lib/ok": "x"}, map[string]string{sharedDir: libDir})
	if err := unzipRoot(root, raw); err != nil {
		t.Fatal(err)
	}
	got, err := os.Readlink(filepath.Join(dir, sharedDir))
	if err != nil || got != libDir {
		t.Fatalf("readlink %q %v", got, err)
	}
}

func TestUnzipSymlinkSkipsEscape(t *testing.T) {
	dir := t.TempDir()
	root, err := os.OpenRoot(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer root.Close()
	raw := makeLinkZip(t, map[string]string{zipSafeName: zipPayload}, map[string]string{"link": zipEvilUp})
	if err := unzipRoot(root, raw); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Lstat(filepath.Join(dir, "link")); !os.IsNotExist(err) {
		t.Fatal("escaped symlink was created")
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

func makeLinkZip(t *testing.T, files, links map[string]string) []byte {
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
	for name, target := range links {
		h := &zip.FileHeader{Name: name}
		h.SetMode(os.ModeSymlink | zipDirPerm)
		w, err := zw.CreateHeader(h)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := w.Write([]byte(target)); err != nil {
			t.Fatal(err)
		}
	}
	if err := zw.Close(); err != nil {
		t.Fatal(err)
	}
	return buf.Bytes()
}
