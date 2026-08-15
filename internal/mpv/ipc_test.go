package mpv

import (
	"bufio"
	"bytes"
	"encoding/json"
	"net"
	"strings"
	"sync"
	"testing"
)

const (
	ipcCmdLoad   = "loadfile"
	ipcPropPause = "pause"
)

func TestIPCCmdFraming(t *testing.T) {
	srv, cli := net.Pipe()
	defer srv.Close()
	defer cli.Close()
	ipc := newIPC(cli)
	var wg sync.WaitGroup
	wg.Go(func() {
		line, err := bufio.NewReader(srv).ReadString('\n')
		if err != nil {
			t.Errorf("read: %v", err)
			return
		}
		var msg map[string]any
		if err := json.Unmarshal([]byte(line), &msg); err != nil {
			t.Errorf("json: %v", err)
			return
		}
		cmd, _ := msg["command"].([]any)
		if len(cmd) != 2 || cmd[0] != ipcCmdLoad || cmd[1] != "http://x" {
			t.Errorf("command %v", cmd)
		}
		if msg["request_id"] != float64(1) {
			t.Errorf("id %v", msg["request_id"])
		}
	})
	if err := ipc.cmd(ipcCmdLoad, "http://x"); err != nil {
		t.Fatal(err)
	}
	wg.Wait()
}

func TestIPCReadMalformedJSON(t *testing.T) {
	rwc := &readWriteBuf{buf: bytes.NewBufferString("not-json\n")}
	ipc := newIPC(rwc)
	if _, err := ipc.read(); err == nil || !strings.Contains(err.Error(), "json") {
		t.Fatalf("err %v", err)
	}
}

type readWriteBuf struct {
	buf *bytes.Buffer
}

func (r *readWriteBuf) Read(p []byte) (int, error) { return r.buf.Read(p) }
func (*readWriteBuf) Write(p []byte) (int, error)  { return len(p), nil }
func (*readWriteBuf) Close() error                 { return nil }

func TestIPCObserveSendsCommand(t *testing.T) {
	srv, cli := net.Pipe()
	defer srv.Close()
	defer cli.Close()
	ipc := newIPC(cli)
	done := make(chan string, 1)
	go func() {
		line, _ := bufio.NewReader(srv).ReadString('\n')
		done <- line
	}()
	if err := ipc.observe(ipcPropPause); err != nil {
		t.Fatal(err)
	}
	line := <-done
	if !strings.Contains(line, "observe_property") || !strings.Contains(line, ipcPropPause) {
		t.Fatalf("line %q", line)
	}
}
