// Package config stores the server URL, session token, and device id.
package config

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
)

// File is the on-disk session. Token is stored, password is not.
type File struct {
	Server   string `json:"server"`
	User     string `json:"user"`
	Token    string `json:"token"`
	UserID   string `json:"userId"`
	ServerID string `json:"serverId"`
	DeviceID string `json:"deviceId"`
	AutoNext bool   `json:"autoNext"`
	path     string
}

func Dir() (string, error) {
	base, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(base, appName), nil
}

func Path() (string, error) {
	d, err := Dir()
	if err != nil {
		return "", err
	}
	return filepath.Join(d, configFile), nil
}

func Load() (*File, error) {
	d, err := Dir()
	if err != nil {
		return nil, err
	}
	p := filepath.Join(d, configFile)
	raw, err := readConfig(d)
	if err != nil {
		if os.IsNotExist(err) {
			f := &File{AutoNext: true, path: p, DeviceID: newDeviceID()}
			return f, nil
		}
		return nil, err
	}
	f := &File{AutoNext: true, path: p}
	if err := json.Unmarshal(raw, f); err != nil {
		return nil, err
	}
	f.path = p
	if f.DeviceID == "" {
		f.DeviceID = newDeviceID()
	}
	return f, nil
}

func readConfig(dir string) ([]byte, error) {
	root, err := os.OpenRoot(dir)
	if err != nil {
		return nil, err
	}
	defer root.Close()
	return root.ReadFile(configFile)
}

func (f *File) Save() error {
	if f.path == "" {
		p, err := Path()
		if err != nil {
			return err
		}
		f.path = p
	}
	dir := filepath.Dir(f.path)
	if err := os.MkdirAll(dir, dirPerm); err != nil {
		return err
	}
	raw, err := json.MarshalIndent(f, "", "  ")
	if err != nil {
		return err
	}
	root, err := os.OpenRoot(dir)
	if err != nil {
		return err
	}
	defer root.Close()
	if err := root.WriteFile(configTmp, raw, filePerm); err != nil {
		return err
	}
	return root.Rename(configTmp, configFile)
}

func (f *File) ClearSession() {
	f.Token = ""
	f.UserID = ""
	f.ServerID = ""
}

func newDeviceID() string {
	var b [deviceBytes]byte
	if _, err := rand.Read(b[:]); err != nil {
		return fallbackID
	}
	return hex.EncodeToString(b[:])
}

func DeviceName() string {
	h, err := os.Hostname()
	if err != nil || h == "" {
		h = runtime.GOOS
	}
	return h + deviceSuffix
}

func CacheDir() (string, error) {
	base, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	d := filepath.Join(base, appName)
	return d, os.MkdirAll(d, dirPerm)
}
