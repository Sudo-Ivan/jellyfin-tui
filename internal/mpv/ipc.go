package mpv

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"sync"
	"sync/atomic"
)

type ipc struct {
	c      io.ReadWriteCloser
	w      *bufio.Writer
	r      *bufio.Reader
	mu     sync.Mutex
	nextID atomic.Int64
}

func newIPC(c io.ReadWriteCloser) *ipc {
	return &ipc{
		c: c,
		w: bufio.NewWriter(c),
		r: bufio.NewReader(c),
	}
}

func (i *ipc) cmd(args ...any) error {
	id := i.nextID.Add(1)
	msg := map[string]any{"command": args, "request_id": id}
	raw, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	i.mu.Lock()
	defer i.mu.Unlock()
	if _, err := i.w.Write(raw); err != nil {
		return err
	}
	if err := i.w.WriteByte('\n'); err != nil {
		return err
	}
	return i.w.Flush()
}

func (i *ipc) observe(name string) error {
	id := i.nextID.Add(1)
	return i.cmd("observe_property", id, name)
}

func (i *ipc) read() (map[string]any, error) {
	line, err := i.r.ReadBytes('\n')
	if err != nil {
		return nil, err
	}
	var m map[string]any
	if err := json.Unmarshal(line, &m); err != nil {
		return nil, fmt.Errorf("mpv ipc json: %w", err)
	}
	return m, nil
}

func (i *ipc) close() error {
	return i.c.Close()
}
