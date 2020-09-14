package nntp

import (
	"context"
	"fmt"
	"time"
)

func Stream(ctx context.Context, options ...Option) <-chan interface{} {
	cfg := defaultConfig()
	for _, option := range options {
		option(cfg)
	}
	var (
		ch      chan interface{} = make(chan interface{})
		cli     *Client
		current int = -1
	)
	oneTick := func(baseContext context.Context) error {
		if cfg.subscribeLog {
			ch <- &Log{Level: LogLevelDebug, Message: "one tick"}
		}
		ctx, ctxCancel := context.WithTimeout(baseContext, cfg.timeout)
		defer ctxCancel()
		if cli == nil {
			c, err := Connect(ctx, cfg.addr)
			if err != nil {
				return err
			}
			if cfg.subscribeLog {
				ch <- &Log{Level: LogLevelInfo, Message: "successfully connected"}
			}
			cli = c
		}

		cancel := cli.SetDeadline(ctx)
		defer cancel()

		group, err := cli.Group(cfg.group)
		if err != nil {
			return err
		}
		high := group.High
		if cfg.subscribeLog {
			ch <- &Log{Level: LogLevelDebug, Message: fmt.Sprintf("%s group high is %d", cfg.group, high)}
		}
		if current == -1 {
			current = high + cfg.offset
		}
		for current < high {
			a, err := cli.Article(current + 1)
			if err != nil {
				return err
			}
			ch <- a
			current++
		}
		return nil
	}
	reset := func() {
		if cli != nil {
			if cfg.subscribeLog {
				ch <- &Log{Level: LogLevelInfo, Message: "close connection"}
			}
			cli.Close()
			cli = nil
		}
	}

	go func() {
		defer close(ch)

		ticker := time.NewTicker(cfg.tick)
		defer ticker.Stop()
		for {
			if err := oneTick(ctx); err != nil {
				if cfg.subscribeLog {
					ch <- &Log{Level: LogLevelError, Message: err.Error()}
				}
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
