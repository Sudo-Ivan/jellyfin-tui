package mpv

import "strings"

// Track is one audio or subtitle stream from mpv track-list.
type Track struct {
	ID       int
	Type     string
	Lang     string
	Title    string
	Selected bool
}

func parseTracks(data any) []Track {
	raw, ok := data.([]any)
	if !ok {
		return nil
	}
	out := make([]Track, 0, len(raw))
	for _, item := range raw {
		m, ok := item.(map[string]any)
		if !ok {
			continue
		}
		typ, _ := m["type"].(string)
		if typ != trackAudio && typ != trackSub {
			continue
		}
		tr := Track{
			ID:       asInt(m["id"]),
			Type:     typ,
			Lang:     strings.ToLower(asString(m["lang"])),
			Title:    asString(m["title"]),
			Selected: asBool(m["selected"]),
		}
		if tr.Title == "" {
			tr.Title = tr.Lang
		}
		out = append(out, tr)
	}
	return out
}

func asString(v any) string {
	s, _ := v.(string)
	return s
}

func asBool(v any) bool {
	b, _ := v.(bool)
	return b
}
