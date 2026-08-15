.PHONY: all build test fmt fix lint gosec tools tidy windows clean fetch-mpv build-embed

GO      ?= go
BINDIR  := bin
BIN     := $(BINDIR)/jellyfin-tui
PKG     := ./cmd/jellyfin-tui
REVIVE  ?= revive
GOSEC   ?= gosec
REVIVE_VER := github.com/mgechev/revive@v1.13.0
GOSEC_VER  := github.com/securego/gosec/v2/cmd/gosec@v2.22.10

all: fmt fix test lint gosec build

build:
	mkdir -p $(BINDIR)
	$(GO) build -trimpath -ldflags='-s -w' -o $(BIN) $(PKG)

fetch-mpv:
	bash scripts/fetch-mpv.sh --all

build-embed: fetch-mpv build

windows-embed: fetch-mpv
	mkdir -p $(BINDIR)
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GO) build -trimpath -ldflags='-s -w' -o $(BIN).exe $(PKG)

windows:
	mkdir -p $(BINDIR)
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GO) build -trimpath -ldflags='-s -w' -o $(BIN).exe $(PKG)

test:
	$(GO) test ./...

fmt:
	$(GO) fmt ./...

fix:
	$(GO) fix ./...

lint:
	$(REVIVE) -config revive.toml -formatter friendly ./...

gosec:
	$(GOSEC) -quiet ./...

tools:
	$(GO) install $(REVIVE_VER)
	$(GO) install $(GOSEC_VER)

tidy:
	$(GO) mod tidy

clean:
	rm -rf $(BINDIR)
	rm -rf .cache/mpv-fetch-*
	rm -f jellyfin-tui jellyfin-tui.exe
