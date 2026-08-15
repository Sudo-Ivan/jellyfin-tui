package jellyfin

import (
	"net/http"
	"strings"
	"testing"
)

const (
	devSpecial = `host"name`
)

func TestEscapeAuthQuotesAndBackslashes(t *testing.T) {
	if escapeAuth(`a"b`) != `a\"b` {
		t.Fatal("quote")
	}
	if escapeAuth(`a\b`) != `a\\b` {
		t.Fatal("backslash")
	}
	if escapeAuth(devSpecial) != `host\"name` {
		t.Fatal("special")
	}
}

func TestAuthHeaderOracle(t *testing.T) {
	c := New("https://x", devSpecial, "dev-1")
	c.Token = "sess-tok"
	got := c.authHeader()
	if !strings.Contains(got, `Client="`+clientName+`"`) {
		t.Fatalf("client %q", got)
	}
	if !strings.Contains(got, `Device="host\"name"`) {
		t.Fatalf("device %q", got)
	}
	if !strings.Contains(got, `Token="sess-tok"`) {
		t.Fatalf("token %q", got)
	}
	if !strings.Contains(got, `Version="`+clientVersion+`"`) {
		t.Fatalf("version %q", got)
	}
}

func TestAuthHeaderNoTokenOmitsField(t *testing.T) {
	c := New("https://x", "pc", "dev-1")
	got := c.authHeader()
	if strings.Contains(got, "Token=") {
		t.Fatalf("token present %q", got)
	}
}

func TestLoginRejectsEmptyToken(t *testing.T) {
	srv := mockAPI(t, http.StatusOK, `{"User":{"Id":"","Name":""},"AccessToken":""}`, &hit{})
	defer srv.Close()
	err := testClient(srv).Login("u", "p")
	if err == nil || !strings.Contains(err.Error(), "empty session") {
		t.Fatalf("err %v", err)
	}
}

func TestLoginStoresSession(t *testing.T) {
	body := `{"User":{"Id":"uid","Name":"Bob"},"AccessToken":"tok","ServerId":"srv"}`
	var h hit
	srv := mockAPI(t, http.StatusOK, body, &h)
	defer srv.Close()
	c := testClient(srv)
	c.Token = ""
	if err := c.Login("Bob", "secret"); err != nil {
		t.Fatal(err)
	}
	if c.Token != "tok" || c.UserID != "uid" || c.ServerID != "srv" {
		t.Fatalf("session %+v", c)
	}
	if h.method != "POST" || h.path != pathAuth {
		t.Fatalf("auth %s %s", h.method, h.path)
	}
}
