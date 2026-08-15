// Package mpv controls a bundled or PATH mpv process over JSON IPC.
//
// Place mpv.zip, mpv, or mpv.exe in bundle/linux or bundle/windows before
// building to embed the player in the binary. Run scripts/fetch-mpv.sh to
// download current portable builds. The zip is extracted to the user cache on
// first run. Without a bundle the player uses mpv from PATH.
package mpv

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"sync"
	"time"
)

// Player is a long-lived mpv process driven over JSON IPC.
type Player struct {
	cmd    *exec.Cmd
	ipc    *ipc
	bin    string
	sock   string
	mu     sync.Mutex
	pos    time.Duration
	dur    time.Duration
	paused bool
	idle   bool
	eof    bool
	vol    int
	path   string
	tracks []Track
	onEOF  func()
}

// Status is a snapshot of observed mpv properties.
type Status struct {
	Pos    time.Duration
	Dur    time.Duration
	Paused bool
	Idle   bool
	EOF    bool
	Volume int
	Path   string
}

// Start launches mpv in idle mode and connects JSON IPC.
func Start(onEOF func()) (*Player, error) {
	bin, err := resolveBinary()
	if err != nil {
		return nil, err
	}
	sock, err := ipcPath()
	if err != nil {
		return nil, err
	}
	_ = os.Remove(sock)
	cmd := exec.Command(bin, mpvArgs(sock)...) // #nosec G204
	prepareCmd(cmd)
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("start mpv: %w", err)
	}
	conn, err := dialIPC(sock, ipcTimeout)
	if err != nil {
		_ = cmd.Process.Kill()
		return nil, fmt.Errorf("mpv ipc: %w", err)
	}
	p := &Player{cmd: cmd, ipc: conn, bin: bin, sock: sock, vol: defaultVolume, idle: true, onEOF: onEOF}
	go p.readLoop()
	for _, prop := range []string{
		propTime, propDuration, propPause, propVolume, propEOF, propIdle, propPath, propTrack,
	} {
		_ = p.ipc.observe(prop)
	}
	return p, nil
}

func mpvArgs(sock string) []string {
	return []string{
		"--idle=yes",
		"--force-window=yes",
		"--keep-open=yes",
		"--no-terminal",
		"--hwdec=auto",
		ytdlNo,
		"--alang=" + alangEng,
		"--sid=" + sidNo,
		"--input-ipc-server=" + sock,
		"--title=" + titleName,
	}
}

// Close kills mpv and the IPC connection.
func (p *Player) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.ipc != nil {
		_ = p.ipc.close()
		p.ipc = nil
	}
	if p.cmd != nil && p.cmd.Process != nil {
		_ = p.cmd.Process.Kill()
		_, _ = p.cmd.Process.Wait()
		p.cmd = nil
	}
	if p.sock != "" && runtime.GOOS != "windows" {
		_ = os.Remove(p.sock)
	}
	return nil
}

// Load replaces the current file. start is an optional absolute offset.
func (p *Player) Load(url, title string, start time.Duration, headers []string) error {
	p.mu.Lock()
	p.eof = false
	p.idle = false
	p.mu.Unlock()
	if len(headers) > 0 {
		_ = p.ipc.cmd(ipcSetProp, "http-header-fields", headers)
	}
	if err := p.ipc.cmd("loadfile", url, flagReplace, 0, loadOpts(title, start)); err != nil {
		return err
	}
	_ = p.ipc.cmd(ipcSetProp, propPause, false)
	return nil
}

func loadOpts(title string, start time.Duration) map[string]string {
	opts := map[string]string{}
	if title != "" {
		opts[optTitle] = title
	}
	if start > 0 {
		opts[optStart] = strconv.FormatInt(int64(start.Seconds()), 10)
	}
	return opts
}

// PauseToggle cycles the pause property.
func (p *Player) PauseToggle() error {
	return p.ipc.cmd("cycle", propPause)
}

// Stop unloads the current file.
func (p *Player) Stop() error {
	return p.ipc.cmd("stop")
}

// Seek moves relative to the current position.
func (p *Player) Seek(delta time.Duration) error {
	return p.ipc.cmd("seek", delta.Seconds(), "relative")
}

// SeekAbs seeks to an absolute timestamp.
func (p *Player) SeekAbs(pos time.Duration) error {
	return p.ipc.cmd("seek", pos.Seconds(), "absolute")
}

// Volume adds delta to the current volume.
func (p *Player) Volume(delta int) error {
	return p.ipc.cmd("add", propVolume, delta)
}

// Tracks returns the last observed audio and subtitle tracks.
func (p *Player) Tracks() []Track {
	p.mu.Lock()
	defer p.mu.Unlock()
	out := make([]Track, len(p.tracks))
	copy(out, p.tracks)
	return out
}

// SetAudio selects an audio track by mpv id.
func (p *Player) SetAudio(id int) error {
	return p.ipc.cmd(ipcSetProp, propAid, id)
}

// SetSubtitle selects a subtitle track by mpv id.
func (p *Player) SetSubtitle(id int) error {
	return p.ipc.cmd(ipcSetProp, propSid, id)
}

// SetSubtitleNone turns subtitles off.
func (p *Player) SetSubtitleNone() error {
	return p.ipc.cmd(ipcSetProp, propSid, sidNo)
}

// Status returns the last observed property snapshot.
func (p *Player) Status() Status {
	p.mu.Lock()
	defer p.mu.Unlock()
	return Status{
		Pos:    p.pos,
		Dur:    p.dur,
		Paused: p.paused,
		Idle:   p.idle,
		EOF:    p.eof,
		Volume: p.vol,
		Path:   p.path,
	}
}
