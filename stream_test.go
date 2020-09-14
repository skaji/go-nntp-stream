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
	options := []nntp.Option{
		nntp.WithAddr("nntp.perl.org:119"),
		nntp.WithGroup("perl.cpan.uploads"),
		nntp.WithOffset(-2),
		nntp.WithTick(5 * time.Second),
		nntp.WithTimeout(20 * time.Second),
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

	ch := nntp.Stream(ctx, options...)
	for event := range ch {
		switch v := event.(type) {
		case *nntp.Article:
			log.Println(v.ID, v.Header.Get("Subject"))
		case *nntp.Log:
			log.Println(v.Level, v.Message)
		}
	}
}
