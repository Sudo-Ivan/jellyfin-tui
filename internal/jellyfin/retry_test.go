package jellyfin

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestDoRetriesTempThenOK(t *testing.T) {
	n := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		n++
		if n < maxTries {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		w.Header().Set("Content-Type", headerJSON)
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, emptyItems)
	}))
	defer srv.Close()
	if _, err := testClient(srv).Views(); err != nil {
		t.Fatal(err)
	}
	if n != maxTries {
		t.Fatalf("tries %d", n)
	}
}

func TestDoDoesNotRetryAuth(t *testing.T) {
	n := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		n++
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer srv.Close()
	_, err := testClient(srv).Views()
	if err == nil || !IsAuth(err) {
		t.Fatalf("err %v", err)
	}
	if n != 1 {
		t.Fatalf("tries %d", n)
	}
}

func TestGetCachesUntilInvalidate(t *testing.T) {
	n := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		n++
		w.Header().Set("Content-Type", headerJSON)
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, emptyItems)
	}))
	defer srv.Close()
	c := testClient(srv)
	if _, err := c.Views(); err != nil {
		t.Fatal(err)
	}
	if _, err := c.Views(); err != nil {
		t.Fatal(err)
	}
	if n != 1 {
		t.Fatalf("cache miss count %d", n)
	}
	c.Invalidate()
	if _, err := c.Views(); err != nil {
		t.Fatal(err)
	}
	if n != 2 {
		t.Fatalf("after invalidate %d", n)
	}
}

func TestRespCacheExpiry(t *testing.T) {
	c := newCache(time.Millisecond)
	c.put(errPath, []byte("a"))
	if _, ok := c.get(errPath); !ok {
		t.Fatal("fresh")
	}
	time.Sleep(3 * time.Millisecond)
	if _, ok := c.get(errPath); ok {
		t.Fatal("stale")
	}
}

type timeoutErr struct{}

func (timeoutErr) Error() string { return "i/o timeout" }
func (timeoutErr) Timeout() bool { return true }

const errPath = "/x"

func TestIsAuthAndTemp(t *testing.T) {
	if !IsAuth(StatusError{Code: httpAuthMin}) || !IsAuth(StatusError{Code: httpForbidden}) {
		t.Fatal("auth")
	}
	if IsAuth(StatusError{Code: httpTooMany}) {
		t.Fatal("429 is not auth")
	}
	if !IsTemp(StatusError{Code: httpUnavail}) || !IsTemp(timeoutErr{}) {
		t.Fatal("temp")
	}
	if IsTemp(errors.New("nope")) {
		t.Fatal("plain")
	}
}

func TestDoRetries429(t *testing.T) {
	n := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		n++
		if n < maxTries {
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		w.Header().Set("Content-Type", headerJSON)
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, emptyItems)
	}))
	defer srv.Close()
	if _, err := testClient(srv).Views(); err != nil {
		t.Fatal(err)
	}
	if n != maxTries {
		t.Fatalf("tries %d", n)
	}
}

func TestStatusErrTruncatesBody(t *testing.T) {
	long := strings.Repeat("x", maxErrRunes+50)
	err := statusErr("GET", errPath, http.StatusBadRequest, "400", []byte(long))
	se, ok := err.(StatusError)
	if !ok || len(se.Body) != maxErrRunes {
		t.Fatalf("body len %d", len(se.Body))
	}
}
