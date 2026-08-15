package mpv

import (
	"archive/zip"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

func resolveBinary() (string, error) {
	if p, err := extractBundle(); err == nil && p != "" {
		return p, nil
	}
	if p, err := exec.LookPath(mpvName()); err == nil {
		return p, nil
	}
	return "", errors.New("mpv not found on PATH and no embedded bundle in internal/mpv/bundle")
}

func extractBundle() (string, error) {
	raw, err := bundleBytes()
	if err != nil || len(raw) == 0 {
		return "", err
	}
	sum := sha256.Sum256(raw)
	tag := hex.EncodeToString(sum[:hashBytes])
	cache, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(cache, cacheApp, cacheMpv+tag)
	exe := mpvName()
	if st, err := os.Stat(dir); err == nil && st.IsDir() {
		if found := findExe(dir, exe); found != "" {
			return finishBundle(dir, found)
		}
	}
	target := filepath.Join(dir, exe)
	if err := os.MkdirAll(dir, dirPerm); err != nil {
		return "", err
	}
	root, err := os.OpenRoot(dir)
	if err != nil {
		return "", err
	}
	defer root.Close()
	if bytes.HasPrefix(raw, []byte(zipMagic)) {
		if err := unzipRoot(root, raw); err != nil {
			return "", err
		}
		found := findExe(dir, exe)
		if found == "" {
			return "", fmt.Errorf("embedded zip has no %s", exe)
		}
		return finishBundle(dir, found)
	}
	if err := root.WriteFile(exe, raw, binPerm); err != nil {
		return "", err
	}
	return finishBundle(dir, target)
}

func unzipRoot(root *os.Root, raw []byte) error {
	zr, err := zip.NewReader(bytes.NewReader(raw), int64(len(raw)))
	if err != nil {
		return err
	}
	for _, f := range zr.File {
		if err := unzipOne(root, f); err != nil {
			return err
		}
	}
	return nil
}

func unzipOne(root *os.Root, f *zip.File) error {
	name := zipRel(f.Name)
	if name == "" {
		return nil
	}
	if f.FileInfo().IsDir() {
		return root.MkdirAll(name, zipDirPerm)
	}
	if dir := path.Dir(name); dir != "." {
		if err := root.MkdirAll(dir, zipDirPerm); err != nil {
			return err
		}
	}
	if f.Mode()&os.ModeSymlink != 0 {
		return writeZipSymlink(root, f, name)
	}
	return writeZipFile(root, f, name)
}

func zipRel(name string) string {
	name = strings.ReplaceAll(name, `\`, "/")
	name = path.Clean(name)
	if name == "." || name == ".." || strings.HasPrefix(name, "../") || path.IsAbs(name) {
		return ""
	}
	return name
}

func writeZipSymlink(root *os.Root, f *zip.File, name string) error {
	rc, err := f.Open()
	if err != nil {
		return err
	}
	defer rc.Close()
	raw, err := io.ReadAll(io.LimitReader(rc, maxSymlinkTarget+1))
	if err != nil {
		return err
	}
	if len(raw) > maxSymlinkTarget {
		return fmt.Errorf("zip symlink too long")
	}
	target := zipSymlinkTarget(name, string(raw))
	if target == "" {
		return nil
	}
	return root.Symlink(target, name)
}

func zipSymlinkTarget(name, target string) string {
	target = strings.ReplaceAll(strings.TrimSpace(target), `\`, "/")
	if target == "" || path.IsAbs(target) {
		return ""
	}
	cleaned := path.Clean(path.Join(path.Dir(name), target))
	if cleaned == ".." || strings.HasPrefix(cleaned, "../") {
		return ""
	}
	return target
}

func writeZipFile(root *os.Root, f *zip.File, name string) error {
	rc, err := f.Open()
	if err != nil {
		return err
	}
	defer rc.Close()
	mode := f.Mode()
	if mode == 0 {
		mode = zipFilePerm
	}
	out, err := root.OpenFile(name, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, mode)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, io.LimitReader(rc, maxZipMember))
	return err
}

func findExe(dir, name string) string {
	var found string
	_ = filepath.Walk(dir, func(fpath string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		if info.Name() == name {
			found = fpath
			return io.EOF
		}
		return nil
	})
	return found
}

func mpvName() string {
	if runtime.GOOS == "windows" {
		return exeWin
	}
	return exeUnix
}
