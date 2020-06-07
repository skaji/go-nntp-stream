package nntp_test

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"

	nntp "github.com/skaji/go-nntp-stream"
)

func TestStream(t *testing.T) {
	config := nntp.StreamConfig{
		Addr:    "nntp.perl.org:119",
		Group:   "perl.cpan.uploads",
		Offset:  -2,
		Tick:    5 * time.Second,
		Timeout: 20 * time.Second,
	}

	log.SetFlags(log.LstdFlags)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	go func() {
		select {
		case <-sig:
			cancel()
		case <-ctx.Done():
		}
	}()

	ch := nntp.Stream(ctx, config)
	for event := range ch {
		switch event.Type {
		case nntp.EventTypeArticle:
			a := event.Article
			log.Println(a.ID, a.Header.Get("Subject"))
		default:
			log.Println(event)
		}
	}
}
