package nntp

import (
	"context"
	"fmt"
	"time"
)

type StreamConfig struct {
	Addr    string
	Group   string
	Offset  int
	Tick    time.Duration
	Timeout time.Duration
}

func Stream(ctx context.Context, config StreamConfig) <-chan *Event {
	var (
		ch      chan *Event = make(chan *Event)
		cli     *Client
		current int = -1
	)
	oneTick := func(baseContext context.Context) error {
		ch <- &Event{Type: EventTypeDebug, Message: "one tick"}
		ctx, ctxCancel := context.WithTimeout(baseContext, config.Timeout)
		defer ctxCancel()
		if cli == nil {
			c, err := Connect(ctx, config.Addr)
			if err != nil {
				return err
			}
			ch <- &Event{Type: EventTypeInfo, Message: "successfully connected"}
			cli = c
		}

		cancel := cli.SetDeadline(ctx)
		defer cancel()

		group, err := cli.Group(config.Group)
		if err != nil {
			return err
		}
		high := group.High
		ch <- &Event{Type: EventTypeDebug, Message: fmt.Sprintf("%s group high is %d", config.Group, high)}
		if current == -1 {
			current = high + config.Offset
		}
		for current < high {
			a, err := cli.Article(current + 1)
			if err != nil {
				return err
			}
			ch <- &Event{Type: EventTypeArticle, Article: a}
			current++
		}
		return nil
	}
	reset := func() {
		if cli != nil {
			ch <- &Event{Type: EventTypeInfo, Message: "close connection"}
			cli.Close()
			cli = nil
		}
	}

	go func() {
		defer close(ch)

		ticker := time.NewTicker(config.Tick)
		defer ticker.Stop()
		for {
			if err := oneTick(ctx); err != nil {
				ch <- &Event{Type: EventTypeError, Message: err.Error()}
				reset()
			}
			select {
			case <-ctx.Done():
				reset()
				return
			case <-ticker.C:
			}
		}
	}()
	return ch
}
