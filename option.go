package nntp

import "time"

type config struct {
	addr         string
	group        string
	offset       int
	tick         time.Duration
	timeout      time.Duration
	subscribeLog bool
}

func defaultConfig() *config {
	return &config{
		offset:       0,
		tick:         30 * time.Second,
		timeout:      25 * time.Second,
		subscribeLog: true,
	}
}

type Option func(c *config)

func WithAddr(addr string) Option {
	return func(c *config) { c.addr = addr }
}

func WithGroup(group string) Option {
	return func(c *config) { c.group = group }
}

func WithOffset(offset int) Option {
	return func(c *config) { c.offset = offset }
}

func WithTick(tick time.Duration) Option {
	return func(c *config) { c.tick = tick }
}

func WithTimeout(timeout time.Duration) Option {
	return func(c *config) { c.timeout = timeout }
}

func WithSubscribeLog(b bool) Option {
	return func(c *config) { c.subscribeLog = b }
}
