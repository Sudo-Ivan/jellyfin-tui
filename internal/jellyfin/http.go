package jellyfin

import (
	"net"
	"net/http"
)

func newHTTP() *http.Client {
	return &http.Client{Timeout: httpTimeout, Transport: newTransport()}
}

func newTransport() *http.Transport {
	if dt, ok := http.DefaultTransport.(*http.Transport); ok {
		t := dt.Clone()
		applyPool(t)
		return t
	}
	t := &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		ForceAttemptHTTP2:     true,
		TLSHandshakeTimeout:   tlsHandshake,
		ExpectContinueTimeout: expectContinue,
	}
	t.DialContext = (&net.Dialer{Timeout: dialTimeout, KeepAlive: keepAlive}).DialContext
	applyPool(t)
	return t
}

func applyPool(t *http.Transport) {
	t.MaxIdleConns = idleConns
	t.MaxIdleConnsPerHost = idleConnsPerHost
	t.MaxConnsPerHost = maxConnsPerHost
	t.IdleConnTimeout = idleConnTimeout
	t.DisableKeepAlives = false
}
