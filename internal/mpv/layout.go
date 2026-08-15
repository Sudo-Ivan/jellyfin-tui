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
	if _, err := os.Lstat(shared); err == nil {
		return nil
	}
	if !hasSharunLoader(filepath.Join(dir, libDir)) {
		return nil
	}
	return os.Symlink(libDir, shared)
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
