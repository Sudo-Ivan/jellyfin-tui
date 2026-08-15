# jellyfin-tui

Terminal browser for a Jellyfin server. Video and audio play in [mpv](https://mpv.io/). The Go module has no third-party requires: HTTP, JSON, zip extract, and terminal I/O are all standard library.

Requires Go 1.26.6. Builds with `CGO_ENABLED=0` on linux/amd64, linux/arm64, and windows/amd64.

The TUI stamps each frame into a retained cell raster, then writes only dirty runs to the terminal. mpv is a child process on JSON IPC (unix socket on Linux, named pipe on Windows).

## Build

```
make
```

That formats, runs `go fix`, tests, revive, gosec, and produces `./bin/jellyfin-tui`. Separate targets: `make build`, `make test`, `make lint`, `make gosec`, `make windows`.

Install the linters once with `make tools`.

## Run

```
./bin/jellyfin-tui
```

Needs a VT-capable terminal and a Jellyfin server. mpv must be on `PATH` unless you embedded it.

If mpv is on `PATH`, playback works with no extra files. To ship mpv inside the binary, drop `mpv.zip`, `mpv`, or `mpv.exe` in `internal/mpv/bundle/` and rebuild. A zip is extracted to the user cache on first run (this is how Windows builds that include DLLs work). The placeholder file in that directory is only there so `go:embed` has something to compile.

## Config

On first run the login screen asks for server URL, username, and password. The session token is written to `config.json` under the user config dir (`jellyfin-tui/`), not the password. Device id is generated once and reused.

`AutoNext` defaults to true: when an episode ends, the next one loads.

## Keys

| Key | Action |
| --- | --- |
| arrows / hjkl | move |
| enter | open or play |
| tab | switch panes |
| / | search (`genre:` `actor:` `year:`) or filter |
| m | random movie |
| a | newly added |
| space | pause |
| n | next episode |
| < > | seek 10s |
| , . | seek 60s |
| - + | volume |
| r | reload home |
| ctrl-l | full redraw |
| ? | key list |
| q | quit |
| esc | back |
| click | select, click again to open |
| wheel | scroll the pane under the pointer |

Resume position from Jellyfin is passed to mpv as `--start`. Playback progress is posted back every 10 seconds.

## Search

`/` opens search against `GET /Items`. Words without a prefix are a name (`SearchTerm`). Prefixes map to Jellyfin query binders:

```
genre:Action actor:"Harrison Ford" year:1999 blade
g:Sci-Fi y:2010-2012
genre:Action|Thriller
```

`genre` / `g` is pipe-delimited on the wire. `year` / `y` is comma-delimited, or a range (`2010-2012`) expanded client-side. `actor` / `person` / `p` sets `Person` plus `PersonTypes=Actor`. The same prefixes filter a library list after `/` in browse (matched against `Genres`, `People`, and `ProductionYear` on the loaded rows).

`m` plays one movie from `SortBy=Random` (scoped to the current movie library when you are inside one). `a` opens a DateCreated-descending list. Home already shows a shorter newly-added strip.

The client talks to current Jellyfin routes: `/Items`, `/Items/{id}`, `/UserViews`, `/UserItems/Resume`, with `userId` as a query parameter. The old `/Users/{id}/Items` paths are not used.

## LICENSE

[0BSD](LICENSE)
