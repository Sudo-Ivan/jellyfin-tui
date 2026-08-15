//go:build !linux && !windows

package mpv

import "io/fs"

func bundleBytes() ([]byte, error) {
	return nil, fs.ErrNotExist
}
