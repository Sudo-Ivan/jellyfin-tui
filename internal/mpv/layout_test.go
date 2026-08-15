package mpv

import (
	"os"
	"path/filepath"
	"testing"
)

const ldLinuxTest = "ld-linux-x86-64.so.2"

func TestEnsureSharunSharedLinksLib(t *testing.T) {
	dir := t.TempDir()
	lib := filepath.Join(dir, libDir)
	if err := os.Mkdir(lib, zipDirPerm); err != nil {
		t.Fatal(err)
	}
	ld := filepath.Join(lib, ldLinuxTest)
	if err := os.WriteFile(ld, []byte("x"), zipFilePerm); err != nil {
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
	if err := os.WriteFile(filepath.Join(lib, ldLinuxTest), []byte("x"), zipFilePerm); err != nil {
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
