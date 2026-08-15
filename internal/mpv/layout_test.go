package mpv

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

const (
	ldLinuxTest  = "ld-linux-x86-64.so.2"
	ldPayload    = "x"
	minBundleZip = 1 << 20
	mpvVerFlag   = "--version"
)

func TestEnsureSharunSharedLinksLib(t *testing.T) {
	dir := t.TempDir()
	lib := filepath.Join(dir, libDir)
	if err := os.Mkdir(lib, zipDirPerm); err != nil {
		t.Fatal(err)
	}
	ld := filepath.Join(lib, ldLinuxTest)
	if err := os.WriteFile(ld, []byte(ldPayload), zipFilePerm); err != nil {
		t.Fatal(err)
	}
	if err := ensureSharunShared(dir); err != nil {
		t.Fatal(err)
	}
	got, err := os.Readlink(filepath.Join(dir, sharedDir))
	if err != nil || got != libDir {
		t.Fatalf("readlink %q %v", got, err)
	}
	if err := ensureSharunShared(dir); err != nil {
		t.Fatal(err)
	}
}

func TestEnsureSharunSharedSkipsWithoutLoader(t *testing.T) {
	dir := t.TempDir()
	if err := ensureSharunShared(dir); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Lstat(filepath.Join(dir, sharedDir)); !os.IsNotExist(err) {
		t.Fatal("unexpected shared")
	}
}

func TestFinishBundleChmodAndLink(t *testing.T) {
	dir := t.TempDir()
	lib := filepath.Join(dir, libDir)
	if err := os.Mkdir(lib, zipDirPerm); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(lib, ldLinuxTest), []byte(ldPayload), zipFilePerm); err != nil {
		t.Fatal(err)
	}
	exe := filepath.Join(dir, exeUnix)
	if err := os.WriteFile(exe, []byte(zipPayload), zipFilePerm); err != nil {
		t.Fatal(err)
	}
	got, err := finishBundle(dir, exe)
	if err != nil || got != exe {
		t.Fatalf("finishBundle %q %v", got, err)
	}
	st, err := os.Stat(exe)
	if err != nil || st.Mode().Perm()&binPerm != binPerm {
		t.Fatalf("chmod %v %v", st, err)
	}
	link, err := os.Readlink(filepath.Join(dir, sharedDir))
	if err != nil || link != libDir {
		t.Fatalf("shared %q %v", link, err)
	}
}

func TestLinuxBundleSharunExec(t *testing.T) {
	raw, err := bundleBytes()
	if err != nil || len(raw) < minBundleZip {
		t.Skip("no portable mpv zip")
	}
	dir := t.TempDir()
	root, err := os.OpenRoot(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer root.Close()
	if err := unzipRoot(root, raw); err != nil {
		t.Fatal(err)
	}
	found := findExe(dir, mpvName())
	if found == "" {
		t.Fatal("no mpv in zip")
	}
	if _, err := finishBundle(dir, found); err != nil {
		t.Fatal(err)
	}
	if !sharunReady(dir) {
		t.Fatal("sharun loader missing")
	}
	out, err := exec.Command(found, mpvVerFlag).CombinedOutput() // #nosec G204
	if err != nil || !strings.Contains(string(out), exeUnix) {
		t.Fatalf("mpv %s %v", out, err)
	}
}

func TestEnsureSharunSharedKeepsDir(t *testing.T) {
	dir := t.TempDir()
	shared := filepath.Join(dir, sharedDir)
	if err := os.Mkdir(shared, zipDirPerm); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(shared, ldLinuxTest), []byte(ldPayload), zipFilePerm); err != nil {
		t.Fatal(err)
	}
	if err := ensureSharunShared(dir); err != nil {
		t.Fatal(err)
	}
	st, err := os.Lstat(shared)
	if err != nil || !st.IsDir() {
		t.Fatalf("shared %v %v", st, err)
	}
	if !sharunReady(dir) {
		t.Fatal("not ready")
	}
}

func TestEnsureSharunSharedReplacesFile(t *testing.T) {
	dir := t.TempDir()
	lib := filepath.Join(dir, libDir)
	if err := os.Mkdir(lib, zipDirPerm); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(lib, ldLinuxTest), []byte(ldPayload), zipFilePerm); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, sharedDir), []byte(libDir), zipFilePerm); err != nil {
		t.Fatal(err)
	}
	if err := ensureSharunShared(dir); err != nil {
		t.Fatal(err)
	}
	got, err := os.Readlink(filepath.Join(dir, sharedDir))
	if err != nil || got != libDir {
		t.Fatalf("readlink %q %v", got, err)
	}
}
