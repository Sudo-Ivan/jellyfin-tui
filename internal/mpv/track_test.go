package mpv

import "testing"

func TestParseTracks(t *testing.T) {
	raw := []any{
		map[string]any{"id": float64(1), "type": "audio", "lang": "eng", "title": "English", "selected": true},
		map[string]any{"id": float64(2), "type": "sub", "lang": "spa", "title": "Spanish", "selected": false},
		map[string]any{"id": float64(3), "type": "video"},
	}
	got := parseTracks(raw)
	if len(got) != trackParseWant {
		t.Fatalf("len %d", len(got))
	}
	if got[0].Type != trackAudio || !got[0].Selected {
		t.Fatalf("audio %+v", got[0])
	}
	if got[1].Lang != "spa" {
		t.Fatalf("sub %+v", got[1])
	}
}

const trackParseWant = 2
