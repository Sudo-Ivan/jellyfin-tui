package jellyfin

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

const (
	playItemID   = "item-1"
	playSourceID = "src-9"
	playSessID   = "sess-42"
	posSec       = 30
	durSec       = 120
	halfSec      = 500
)

func TestPlayBodyDirectVsTranscode(t *testing.T) {
	direct := playBody(PlayTarget{ItemID: playItemID}, time.Duration(posSec)*time.Second, time.Duration(durSec)*time.Second, false)
	if direct["PlayMethod"] != playDirect {
		t.Fatalf("direct %v", direct["PlayMethod"])
	}
	tc := playBody(PlayTarget{ItemID: playItemID, Transcoding: true}, 0, 0, true)
	paused, _ := tc["IsPaused"].(bool)
	if tc["PlayMethod"] != playTranscode || !paused {
		t.Fatalf("transcode %+v", tc)
	}
}

func TestSourceIDFallback(t *testing.T) {
	if sourceID(PlayTarget{ItemID: playItemID}) != playItemID {
		t.Fatal("item fallback")
	}
	if sourceID(PlayTarget{ItemID: playItemID, MediaSourceID: playSourceID}) != playSourceID {
		t.Fatal("source")
	}
}

func TestTicksRoundTrip(t *testing.T) {
	d := 5*time.Second + halfSec*time.Millisecond
	if ticks(d) != int64(d.Nanoseconds()/nsPerTick) {
		t.Fatal("ticks")
	}
}

func TestPlayingPostsExpectedBody(t *testing.T) {
	var body map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != pathPlaying {
			t.Fatalf("path %s", r.URL.Path)
		}
		raw, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(raw, &body)
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()
	c := testClient(srv)
	tgt := PlayTarget{ItemID: playItemID, MediaSourceID: playSourceID, PlaySessionID: playSessID}
	pos := time.Duration(posSec) * time.Second
	dur := time.Duration(durSec) * time.Second
	if err := c.Playing(tgt, pos, dur, false); err != nil {
		t.Fatal(err)
	}
	if body["ItemId"] != playItemID || body["MediaSourceId"] != playSourceID {
		t.Fatalf("ids %+v", body)
	}
	if body[qPlaySession] != playSessID {
		t.Fatalf("session %+v", body)
	}
}

func TestProgressSetsEventName(t *testing.T) {
	var body map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(raw, &body)
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()
	c := testClient(srv)
	if err := c.Progress(PlayTarget{ItemID: playItemID}, 0, 0, false); err != nil {
		t.Fatal(err)
	}
	if body["EventName"] != eventTime {
		t.Fatalf("event %+v", body)
	}
}
