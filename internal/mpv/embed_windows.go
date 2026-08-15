//go:build windows

package mpv

import (
	"embed"
	"io/fs"
)

//go:embed bundle/windows
var bundleFS embed.FS

func bundleBytes() ([]byte, error) {
	for _, name := range []string{"bundle/windows/mpv.zip", "bundle/windows/mpv.exe"} {
		b, err := bundleFS.ReadFile(name)
		if err == nil && len(b) > 0 {
			return b, nil
		}
	}
	return nil, fs.ErrNotExist
}
