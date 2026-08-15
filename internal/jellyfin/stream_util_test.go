package jellyfin

import (
	"net/url"
	"strings"
	"testing"
)

const (
	srcDirectID = "direct-1"
	srcTransID  = "trans-1"
	baseURL     = "https://jellyfin.example"
)

func TestPickSourceDirectPrefersPlayable(t *testing.T) {
	list := []mediaSource{
		{ID: "bad", SupportsDirectPlay: false, SupportsDirectStream: false},
		{ID: srcDirectID, SupportsDirectStream: true},
	}
	got := pickSource(list, false)
	if got.ID != srcDirectID {
		t.Fatalf("direct %q", got.ID)
	}
}

func TestPickSourceTranscodePrefersTranscoding(t *testing.T) {
	list := []mediaSource{
		{ID: "plain", SupportsDirectPlay: true},
		{ID: srcTransID, SupportsTranscoding: true, TranscodingURL: "/tr"},
	}
	got := pickSource(list, true)
	if got.ID != srcTransID {
		t.Fatalf("transcode %q", got.ID)
	}
}

func TestAbsURLRelativeAndAbsolute(t *testing.T) {
	c := New(baseURL, "dev", "d1")
	if c.absURL("http://other/x") != "http://other/x" {
		t.Fatal("absolute")
	}
	if c.absURL("/Videos/1/stream") != baseURL+"/Videos/1/stream" {
		t.Fatal("rooted")
	}
	if c.absURL("Videos/1/stream") != baseURL+"/Videos/1/stream" {
		t.Fatal("bare")
	}
}

func TestWithTokenPreservesExistingAPIKey(t *testing.T) {
	c := New(baseURL, "dev", "d1")
	c.Token = "secret"
	raw := baseURL + "/x?api_key=keep"
	got := c.withToken(raw)
	u, err := url.Parse(got)
	if err != nil {
		t.Fatal(err)
	}
	if u.Query().Get(qAPIKey) != "keep" {
		t.Fatal("existing key")
	}
	bare := c.withToken(baseURL + "/x")
	u2, _ := url.Parse(bare)
	if u2.Query().Get(qAPIKey) != "secret" {
		t.Fatal("injected key")
	}
}

func TestFirstContainerTable(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"mkv,mp4", "mkv"},
		{" MP4 ", "mp4"},
		{"", ""},
		{"webm", "webm"},
	}
	for _, tc := range cases {
		if got := firstContainer(tc.in); got != tc.want {
			t.Fatalf("firstContainer(%q)=%q", tc.in, got)
		}
	}
}

func TestBuildStreamIncludesContainerExt(t *testing.T) {
	c := New(baseURL, "dev", "d1")
	c.Token = "tok"
	got := c.buildStream("id1", "src1", "mkv,mp4", playSessID, "etag")
	if !strings.Contains(got, "/Videos/id1/stream.mkv") {
		t.Fatalf("path %q", got)
	}
	if !strings.Contains(got, qPlaySession+"="+playSessID) {
		t.Fatal("session")
	}
}
