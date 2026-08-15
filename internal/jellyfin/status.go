package jellyfin

import (
	"errors"
	"fmt"
	"net"
	"strings"
)

// StatusError is a non-2xx HTTP response from the Jellyfin server.
type StatusError struct {
	Method string
	Path   string
	Status string
	Code   int
	Body   string
}

func (e StatusError) Error() string {
	return fmt.Sprintf("%s %s: %s %s", e.Method, e.Path, e.Status, e.Body)
}

func IsAuth(err error) bool {
	var se StatusError
	if !errors.As(err, &se) {
		return false
	}
	return se.Code == httpAuthMin || se.Code == httpForbidden
}

func IsTemp(err error) bool {
	if err == nil {
		return false
	}
	var se StatusError
	if errors.As(err, &se) {
		return se.Code == httpTooMany || se.Code == httpBadGate || se.Code == httpUnavail || se.Code == httpTimeoutC
	}
	var ne net.Error
	if errors.As(err, &ne) {
		return true
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "timeout") || strings.Contains(msg, "connection refused") ||
		strings.Contains(msg, "temporary") || strings.Contains(msg, "eof")
}
