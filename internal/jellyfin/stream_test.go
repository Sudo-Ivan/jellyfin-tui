package jellyfin

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

const (
	idItem     = "abc"
	idSource   = "ms1"
	idPlay     = "ps1"
	extMKV     = "mkv"
	pathMKV    = "/Videos/abc/stream.mkv"
	bodyDirect = `{"PlaySessionId":"ps1","MediaSources":[{` +
		`"Id":"ms1","Container":"mkv","SupportsDirectPlay":true,` +
		`"SupportsDirectStream":true,` +
		`"DirectStreamUrl":"/Videos/abc/stream.mkv?Static=true&mediaSourceId=ms1"}]}`
	bodyContainer = `{"PlaySessionId":"ps1","MediaSources":[{` +
		`"Id":"ms1","Container":"mkv,webm","SupportsDirectStream":true}]}`
)

func TestOpenStreamPrefersDirectURL(t *testing.T) {
	body := bodyDirect
	var h hit
	srv := mockAPI(t, http.StatusOK, body, &h)
	defer srv.Close()
	got := testClient(srv).OpenStream(idItem)
	if h.method != http.MethodPost || !strings.HasSuffix(h.path, pathPlaybackInfo) {
		t.Fatalf("playbackinfo %s %s", h.method, h.path)
	}
	u, err := url.Parse(got.URL)
	if err != nil {
		t.Fatal(err)
	}
	if u.Path != pathMKV {
		t.Fatalf("path %s", u.Path)
	}
	if u.Query().Get(qAPIKey) != "sess" {
		t.Fatal("api_key")
	}
	if got.PlaySessionID != idPlay || got.MediaSourceID != idSource {
		t.Fatalf("%+v", got)
	}
}

func TestOpenStreamBuildsContainer(t *testing.T) {
	body := bodyContainer
	var h hit
	srv := mockAPI(t, http.StatusOK, body, &h)
	defer srv.Close()
	got := testClient(srv).OpenStream(idItem)
	u, err := url.Parse(got.URL)
	if err != nil {
		t.Fatal(err)
	}
	if u.Path != pathMKV {
		t.Fatalf("path %s", u.Path)
	}
	q := u.Query()
	if q.Get(qStatic) != boolTrue || q.Get(qMediaSourceID) != idSource || q.Get(qPlaySession) != idPlay {
		t.Fatalf("%v", q)
	}
}

func TestOpenStreamFallbackOnError(t *testing.T) {
	var h hit
	srv := mockAPI(t, http.StatusInternalServerError, "nope", &h)
	defer srv.Close()
	got := testClient(srv).OpenStream(idItem)
	u, err := url.Parse(got.URL)
	if err != nil {
		t.Fatal(err)
	}
	if u.Path != pathVideos+idItem+pathStream {
		t.Fatalf("path %s", u.Path)
	}
	if u.Query().Get(qStatic) != boolTrue {
		t.Fatal("static")
	}
}

func TestFirstContainer(t *testing.T) {
	if firstContainer("MKV,webm") != extMKV {
		t.Fatal("first")
	}
	if firstContainer("  ") != "" {
		t.Fatal("empty")
	}
}

func TestHTTPPool(t *testing.T) {
	c := New("http://jf", "d", "i")
	tr, ok := c.HTTP.Transport.(*http.Transport)
	if !ok {
		t.Fatal("transport")
	}
	if tr.DisableKeepAlives || tr.MaxIdleConnsPerHost != idleConnsPerHost {
		t.Fatal("reuse")
	}
	if tr.MaxIdleConns != idleConns || tr.MaxConnsPerHost != maxConnsPerHost {
		t.Fatal("pool size")
	}
}

func TestPlaybackInfoRetriesNotNeededForAuth(t *testing.T) {
	n := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		n++
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = io.WriteString(w, "nope")
	}))
	defer srv.Close()
	_ = testClient(srv).OpenStream(idItem)
	if n != 1 {
		t.Fatalf("tries %d", n)
	}
}
