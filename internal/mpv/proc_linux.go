package mpv

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"
)

func ipcPath() (string, error) {
	return filepath.Join(os.TempDir(), fmt.Sprintf("jellyfin-tui-%d.sock", os.Getpid())), nil
}

func dialIPC(path string, timeout time.Duration) (*ipc, error) {
	deadline := time.Now().Add(timeout)
	var last error
	for time.Now().Before(deadline) {
		c, err := net.DialTimeout("unix", path, ipcDialTry)
		if err == nil {
			return newIPC(c), nil
		}
		last = err
		time.Sleep(ipcRetry)
	}
	if last == nil {
		last = fmt.Errorf("timeout connecting to mpv socket")
	}
	return nil, last
}

func prepareCmd(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
}
