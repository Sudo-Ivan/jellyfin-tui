# Agent notes

Stdlib-only Go TUI for Jellyfin. Do not add third-party Go modules.

## Stack

- Go 1.26.6, `CGO_ENABLED=0`
- No bubbletea, tcell, lipgloss, or HTTP client libraries
- mpv as a child process over JSON IPC (unix socket on Linux, named pipe on Windows)
- Jellyfin REST on current 10.11+ routes: `/Items`, `/Items/{id}`, `/UserViews`, `/UserItems/Resume` with `userId` as a query parameter
- Do not use the deprecated `/Users/{id}/Items` paths

## Layout

- `cmd/jellyfin-tui/` process entry, SIGINT/SIGTERM
- `internal/tui/` cell raster, SGR mouse, dirty present
- `internal/jellyfin/` REST client, retry, GET cache
- `internal/mpv/` player process and optional embed
- `internal/app/` screens and input
- `internal/config/` session file

## Build

`make all` formats, runs `go fix`, tests, revive, gosec, and writes `./bin/jellyfin-tui`.

`make windows` writes `./bin/jellyfin-tui.exe`.

Revive: files max 220 lines (skip comments and blanks), functions 70/90, cyclomatic 15. `add-constant` is strict. Name unique integers in tests.

Comments are Go doc comments. No emojis, TODO markers, em dashes, or semicolons in comments.

## Mouse

SGR 1000/1006 on Linux. Windows console mouse input with quick-edit off.

Click selects a row. Click the selected row again to open or play. Wheel scrolls the pane under the pointer. Help closes on a click.

## HTTP

GET responses are cached in memory for a short TTL. `Invalidate` on home reload. Temporary failures (429, 502, 503, 504, timeouts, refused) retry with backoff. 401/403 do not retry and clear the saved session only on auth failure at boot.

SIGINT/SIGTERM and `q` post Stopped, close mpv, then restore the terminal.

## Tests

Table tests and `httptest`. Oracle tests print `QUERY_ORACLE_PROVED`. Cover decode, query encoding, retry, cache, and hit testing.

## Backlog

- subtitle and audio track picking
- watched / unwatched toggle
- transcoding fallback when direct play fails
- poster / image ASCII
- multi-server profiles
- click header and transport bar
- reconnect from the login screen without retyping the password
- click-to-focus login fields
