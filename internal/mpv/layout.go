package mpv

import (
	"os"
	"path/filepath"
	"strings"
)

func finishBundle(dir, exePath string) (string, error) {
	if err := os.Chmod(exePath, binPerm); err != nil {
		return "", err
	}
	if err := ensureSharunShared(dir); err != nil {
		return "", err
	}
	return exePath, nil
}

func ensureSharunShared(dir string) error {
	shared := filepath.Join(dir, sharedDir)
	st, err := os.Lstat(shared)
	if err == nil {
		if st.Mode()&os.ModeSymlink != 0 || st.IsDir() {
			return nil
		}
		if err := os.Remove(shared); err != nil {
			return err
		}
	} else if !os.IsNotExist(err) {
		return err
	}
	if !hasSharunLoader(filepath.Join(dir, libDir)) {
		return nil
	}
	return os.Symlink(libDir, shared)
}

func sharunReady(dir string) bool {
	if hasSharunLoader(filepath.Join(dir, sharedDir)) {
		return true
	}
	if hasSharunLoader(filepath.Join(dir, sharedDir, libDir)) {
		return true
	}
	if _, err := os.Lstat(filepath.Join(dir, sharedDir)); err != nil {
		return false
	}
	return hasSharunLoader(filepath.Join(dir, libDir))
}

func hasSharunLoader(lib string) bool {
	ents, err := os.ReadDir(lib)
	if err != nil {
		return false
	}
	for _, e := range ents {
		n := e.Name()
		if strings.HasPrefix(n, ldLinuxPrefix) || strings.HasPrefix(n, ldMuslPrefix) {
			return true
		}
	}
	return false
}
