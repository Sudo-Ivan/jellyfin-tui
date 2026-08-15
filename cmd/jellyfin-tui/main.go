package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"jellyfin-tui/internal/app"
	"jellyfin-tui/internal/config"
	"jellyfin-tui/internal/tui"
)

func main() {
	if err := run(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	term, err := tui.Open()
	if err != nil {
		return err
	}
	defer func() { _ = term.Restore() }()
	return app.New(term, cfg).Run(ctx)
}
