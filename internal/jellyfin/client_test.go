package jellyfin

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

type hit struct {
	method string
	path   string
	query  url.Values
}

func mockAPI(t *testing.T, status int, body string, h *hit) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.method = r.Method
		h.path = r.URL.Path
		h.query = r.URL.Query()
		w.Header().Set("Content-Type", headerJSON)
		w.WriteHeader(status)
		_, _ = io.WriteString(w, body)
	}))
}

func testClient(srv *httptest.Server) *Client {
	c := New(srv.URL, "dev", "dev-1")
	c.UserID = "user-1"
	c.Token = "sess"
	c.HTTP = srv.Client()
	return c
}

func assertPath(t *testing.T, got, want string) {
	t.Helper()
	if got != want {
		t.Fatalf("got path %q want %q", got, want)
	}
}

const emptyItems = `{"Items":[],"TotalRecordCount":0}`

func TestViewsUsesUserViews(t *testing.T) {
	var h hit
	srv := mockAPI(t, http.StatusOK, emptyItems, &h)
	defer srv.Close()
	_, err := testClient(srv).Views()
	if err != nil {
		t.Fatal(err)
	}
	assertPath(t, h.path, pathUserViews)
	if h.query.Get(qUserID) != "user-1" {
		t.Fatal("userId")
	}
}

func TestResumeUsesUserItemsResume(t *testing.T) {
	var h hit
	srv := mockAPI(t, http.StatusOK, emptyItems, &h)
	defer srv.Close()
	_, err := testClient(srv).Resume(8)
	if err != nil {
		t.Fatal(err)
	}
	assertPath(t, h.path, pathResume)
}

func TestItemUsesItemsID(t *testing.T) {
	var h hit
	srv := mockAPI(t, http.StatusOK, `{"Id":"x","Name":"X","Type":"Movie"}`, &h)
	defer srv.Close()
	it, err := testClient(srv).Item("x")
	if err != nil {
		t.Fatal(err)
	}
	if it.ID != "x" {
		t.Fatal(it.ID)
	}
	assertPath(t, h.path, pathItems+"/x")
}

func TestRandomMovieQuery(t *testing.T) {
	var h hit
	body := `{"Items":[{"Id":"m1","Name":"X","Type":"Movie"}],"TotalRecordCount":1}`
	srv := mockAPI(t, http.StatusOK, body, &h)
	defer srv.Close()
	it, err := testClient(srv).RandomMovie("lib")
	if err != nil {
		t.Fatal(err)
	}
	if it.ID != "m1" {
		t.Fatal(it.ID)
	}
	assertPath(t, h.path, pathItems)
	if h.query.Get(qSortBy) != sortRandom {
		t.Fatalf("sort %s", h.query.Get(qSortBy))
	}
	if h.query.Get(qInclude) != includeMovies {
		t.Fatal("include")
	}
	if h.query.Get(qLimit) != "1" {
		t.Fatal("limit")
	}
	if h.query.Get(qParent) != "lib" {
		t.Fatal("parent")
	}
}

func TestRandomMovieEmpty(t *testing.T) {
	var h hit
	srv := mockAPI(t, http.StatusOK, emptyItems, &h)
	defer srv.Close()
	_, err := testClient(srv).RandomMovie("")
	if err != errNoMovie {
		t.Fatalf("err %v", err)
	}
}

func TestNewlyAddedQuery(t *testing.T) {
	var h hit
	srv := mockAPI(t, http.StatusOK, emptyItems, &h)
	defer srv.Close()
	_, err := testClient(srv).NewlyAdded(limitAdded)
	if err != nil {
		t.Fatal(err)
	}
	assertPath(t, h.path, pathItems)
	if strings.Contains(h.path, "/Users/") {
		t.Fatal("deprecated users items path")
	}
	if h.query.Get(qSortBy) != sortCreated || h.query.Get(qSortOrd) != orderDesc {
		t.Fatalf("sort %s %s", h.query.Get(qSortBy), h.query.Get(qSortOrd))
	}
}

func TestSearchSendsGenreActorYear(t *testing.T) {
	var h hit
	srv := mockAPI(t, http.StatusOK, emptyItems, &h)
	defer srv.Close()
	_, err := testClient(srv).Search(`genre:Action actor:Ford year:1999 blade`, limitSearch)
	if err != nil {
		t.Fatal(err)
	}
	assertPath(t, h.path, pathItems)
	if h.query.Get(qGenres) != nameAction {
		t.Fatalf("genres %q", h.query.Get(qGenres))
	}
	if h.query.Get(qPerson) != nameFord || h.query.Get(qPersonTypes) != personActor {
		t.Fatalf("person %s %s", h.query.Get(qPerson), h.query.Get(qPersonTypes))
	}
	if h.query.Get(qYears) != "1999" {
		t.Fatalf("years %q", h.query.Get(qYears))
	}
	if h.query.Get(qSearch) != nameBlade {
		t.Fatalf("term %q", h.query.Get(qSearch))
	}
}

func TestHTTPErrorStatus(t *testing.T) {
	var h hit
	srv := mockAPI(t, http.StatusUnauthorized, "nope", &h)
	defer srv.Close()
	_, err := testClient(srv).Views()
	if err == nil {
		t.Fatal("want error")
	}
}

func TestStreamURL(t *testing.T) {
	c := New("http://jf", "d", "i")
	c.Token = "sess"
	u, err := url.Parse(c.StreamURL("abc"))
	if err != nil {
		t.Fatal(err)
	}
	if u.Path != "/Videos/abc/stream" {
		t.Fatalf("stream %s", u.Path)
	}
	q := u.Query()
	if q.Get("static") != "true" || q.Get("api_key") != "sess" {
		t.Fatalf("%v", q)
	}
}

func TestQueryResultDecode(t *testing.T) {
	var r QueryResult
	raw := []byte(`{"Items":[{"Id":"1","Genres":["A"],"People":[{"Name":"X","Type":"Actor"}],"ProductionYear":1999}],"TotalRecordCount":1}`)
	if err := json.Unmarshal(raw, &r); err != nil {
		t.Fatal(err)
	}
	if r.TotalRecordCount != 1 || r.Items[0].Genres[0] != "A" || r.Items[0].People[0].Type != personActor {
		t.Fatalf("%+v", r)
	}
}
