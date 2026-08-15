package mpv

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func openDevNull() *os.File {
	f, err := os.OpenFile(os.DevNull, os.O_RDWR, 0)
	if err != nil {
		return nil
	}
	return f
}

func abortPlayer(cmd *exec.Cmd, sock string, stderr *bytes.Buffer, err error) error {
	if cmd.Process != nil {
		_ = cmd.Process.Kill()
		_, _ = cmd.Process.Wait()
	}
	if sock != "" && runtime.GOOS != "windows" {
		_ = os.Remove(sock)
	}
	if msg := strings.TrimSpace(stderr.String()); msg != "" {
		return fmt.Errorf("mpv ipc: %w (%s)", err, msg)
	}
	return fmt.Errorf("mpv ipc: %w", err)
}
