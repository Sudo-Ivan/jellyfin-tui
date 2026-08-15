package jellyfin

import (
	"sync"
	"time"
)

type cacheEnt struct {
	at  time.Time
	raw []byte
}

type respCache struct {
	mu  sync.Mutex
	ttl time.Duration
	m   map[string]cacheEnt
}

func newCache(ttl time.Duration) *respCache {
	return &respCache{ttl: ttl, m: map[string]cacheEnt{}}
}

func (c *respCache) get(key string) ([]byte, bool) {
	if c == nil {
		return nil, false
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	ent, ok := c.m[key]
	if !ok || time.Since(ent.at) > c.ttl {
		delete(c.m, key)
		return nil, false
	}
	return ent.raw, true
}

func (c *respCache) put(key string, raw []byte) {
	if c == nil || key == "" {
		return
	}
	c.mu.Lock()
	c.m[key] = cacheEnt{at: time.Now(), raw: append([]byte(nil), raw...)}
	c.mu.Unlock()
}

func (c *respCache) clear() {
	if c == nil {
		return
	}
	c.mu.Lock()
	c.m = map[string]cacheEnt{}
	c.mu.Unlock()
}

// Invalidate drops cached GET bodies so the next fetch hits the server.
func (c *Client) Invalidate() {
	c.cache.clear()
}

func cacheKey(path string, query string) string {
	if query == "" {
		return path
	}
	return path + "?" + query
}
