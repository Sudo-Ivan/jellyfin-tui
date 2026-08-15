package mpv

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"time"
	"unsafe"
)

func ipcPath() (string, error) {
	return fmt.Sprintf(`\\.\pipe\jellyfin-tui-%d`, os.Getpid()), nil
}

func dialIPC(path string, timeout time.Duration) (*ipc, error) {
	deadline := time.Now().Add(timeout)
	var last error
	for time.Now().Before(deadline) {
		f, err := openPipe(path)
		if err == nil {
			return newIPC(f), nil
		}
		last = err
		time.Sleep(ipcRetry)
	}
	if last == nil {
		last = fmt.Errorf("timeout connecting to mpv pipe")
	}
	return nil, last
}

func prepareCmd(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: false}
}

var (
	k32            = syscall.NewLazyDLL("kernel32.dll")
	procCreateFile = k32.NewProc("CreateFileW")
)

const (
	genericRead    = 0x80000000
	genericWrite   = 0x40000000
	openExisting   = 3
	fileAttrNormal = 0x80
	fileShareRead  = 0x1
	fileShareWrite = 0x2
	invalidHandle  = ^uintptr(0)
)

func openPipe(path string) (*os.File, error) {
	p, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return nil, err
	}
	r1, _, e := procCreateFile.Call(
		uintptr(unsafe.Pointer(p)), // #nosec G103
		uintptr(genericRead|genericWrite),
		uintptr(fileShareRead|fileShareWrite),
		0,
		uintptr(openExisting),
		uintptr(fileAttrNormal),
		0,
	)
	if r1 == 0 || r1 == invalidHandle {
		if e != syscall.Errno(0) {
			return nil, e
		}
		return nil, syscall.EINVAL
	}
	return os.NewFile(r1, path), nil
}
