// Package jellyfin is a stdlib HTTP client for the Jellyfin REST API.
package jellyfin

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// Client talks to a Jellyfin server with the REST API.
type Client struct {
	Base     string
	Token    string
	UserID   string
	ServerID string
	UserName string
	Device   string
	DeviceID string
	HTTP     *http.Client
	cache    *respCache
}

// New returns a client pointed at base with no session token yet.
func New(base, device, deviceID string) *Client {
	return &Client{
		Base:     strings.TrimRight(base, "/"),
		Device:   device,
		DeviceID: deviceID,
		HTTP:     newHTTP(),
		cache:    newCache(cacheTTL),
	}
}

func (c *Client) authHeader() string {
	parts := []string{
		"MediaBrowser Client=" + quoted(clientName),
		"Device=" + quoted(escapeAuth(c.Device)),
		"DeviceId=" + quoted(c.DeviceID),
		"Version=" + quoted(clientVersion),
	}
	if c.Token != "" {
		parts = append(parts, "Token="+quoted(c.Token))
	}
	return strings.Join(parts, ", ")
}

func quoted(s string) string {
	return authQuote + s + authQuote
}

func escapeAuth(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `"`, `\"`)
	return s
}

func (c *Client) do(method, path string, query url.Values, body any, out any) error {
	var payload []byte
	if body != nil {
		raw, err := json.Marshal(body)
		if err != nil {
			return err
		}
		payload = raw
	}
	raw, err := c.doRetry(method, path, query, payload)
	if err != nil {
		return err
	}
	return decodeOut(raw, out)
}

func (c *Client) doRetry(method, path string, query url.Values, payload []byte) ([]byte, error) {
	var last error
	for i := range maxTries {
		if i > 0 {
			time.Sleep(retryBackoff * time.Duration(1<<(i-1)))
		}
		raw, err := c.doOnce(method, path, query, payload)
		if err == nil {
			return raw, nil
		}
		last = err
		if !IsTemp(err) {
			return nil, err
		}
	}
	return nil, last
}

func (c *Client) doOnce(method, path string, query url.Values, payload []byte) ([]byte, error) {
	u := c.Base + path
	if query != nil {
		u += "?" + query.Encode()
	}
	var rdr io.Reader
	if payload != nil {
		rdr = bytes.NewReader(payload)
	}
	req, err := http.NewRequestWithContext(context.Background(), method, u, rdr)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", headerJSON)
	req.Header.Set(headerAuth, c.authHeader())
	if c.Token != "" {
		req.Header.Set(headerToken, c.Token)
	}
	if payload != nil {
		req.Header.Set("Content-Type", headerJSON)
	}
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	raw, err := io.ReadAll(io.LimitReader(resp.Body, maxBodyBytes))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < httpOKMin || resp.StatusCode >= httpOKMax {
		return nil, statusErr(method, path, resp.StatusCode, resp.Status, raw)
	}
	return raw, nil
}

func statusErr(method, path string, code int, status string, raw []byte) error {
	msg := strings.TrimSpace(string(raw))
	if len(msg) > maxErrRunes {
		msg = msg[:maxErrRunes]
	}
	return StatusError{Method: method, Path: path, Status: status, Code: code, Body: msg}
}

func decodeOut(raw []byte, out any) error {
	if out == nil || len(raw) == 0 {
		return nil
	}
	return json.Unmarshal(raw, out)
}

func (c *Client) get(path string, query url.Values, out any) error {
	q := ""
	if query != nil {
		q = query.Encode()
	}
	key := cacheKey(path, q)
	if raw, ok := c.cache.get(key); ok {
		return decodeOut(raw, out)
	}
	raw, err := c.doRetry(http.MethodGet, path, query, nil)
	if err != nil {
		return err
	}
	c.cache.put(key, raw)
	return decodeOut(raw, out)
}

func (c *Client) post(path string, query url.Values, body any, out any) error {
	return c.do(http.MethodPost, path, query, body, out)
}

func limitQ(n int) string { return strconv.Itoa(n) }

func ticks(d time.Duration) int64 {
	return d.Nanoseconds() / nsPerTick
}
