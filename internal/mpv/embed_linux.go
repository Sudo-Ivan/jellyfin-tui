//go:build linux

package mpv

import (
	"embed"
	"io/fs"
)

//go:embed bundle/linux
var bundleFS embed.FS

func bundleBytes() ([]byte, error) {
	for _, name := range []string{"bundle/linux/mpv.zip", "bundle/linux/mpv"} {
		b, err := bundleFS.ReadFile(name)
		if err == nil && len(b) > 0 {
			return b, nil
		}
	}
	return nil, fs.ErrNotExist
}
