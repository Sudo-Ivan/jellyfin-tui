package mpv

import (
	"encoding/json"
	"time"
)

func (p *Player) readLoop() {
	for {
		msg, err := p.ipc.read()
		if err != nil {
			break
		}
		p.apply(msg)
	}
	_ = p.Close()
	p.mu.Lock()
	p.idle = true
	p.eof = true
	p.mu.Unlock()
	if p.onEOF != nil {
		p.onEOF()
	}
}

func (p *Player) apply(m map[string]any) {
	if ev, _ := m["event"].(string); ev == evEndFile {
		p.applyEndFile(m)
		return
	}
	if m["event"] != evProp {
		return
	}
	p.applyProp(m)
}

func (p *Player) applyEndFile(m map[string]any) {
	reason, _ := m["reason"].(string)
	if reason != reasonEOF && reason != reasonStop {
		return
	}
	p.mu.Lock()
	p.eof = reason == reasonEOF
	p.mu.Unlock()
	if reason == reasonEOF && p.onEOF != nil {
		p.onEOF()
	}
}

func (p *Player) applyProp(m map[string]any) {
	name, _ := m["name"].(string)
	data := m["data"]
	p.mu.Lock()
	defer p.mu.Unlock()
	switch name {
	case propTime:
		p.pos = asDuration(data)
	case propDuration:
		p.dur = asDuration(data)
	case propPause:
		p.paused, _ = data.(bool)
	case propVolume:
		p.vol = asInt(data)
	case propEOF:
		p.eof, _ = data.(bool)
	case propIdle:
		p.idle, _ = data.(bool)
	case propPath:
		p.path, _ = data.(string)
	case propTrack:
		p.tracks = parseTracks(data)
	}
}

func asDuration(v any) time.Duration {
	switch n := v.(type) {
	case float64:
		return time.Duration(n * float64(time.Second))
	case json.Number:
		f, _ := n.Float64()
		return time.Duration(f * float64(time.Second))
	default:
		return 0
	}
}

func asInt(v any) int {
	switch n := v.(type) {
	case float64:
		return int(n)
	case json.Number:
		i, _ := n.Int64()
		return int(i)
	default:
		return 0
	}
}
