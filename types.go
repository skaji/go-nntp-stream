package nntp

import "net/mail"

type Article struct {
	ID     int
	Header mail.Header
	Body   []byte
}

type Log struct {
	Level   LogLevel
	Message string
}

type LogLevel string

var (
	LogLevelDebug LogLevel = "debug"
	LogLevelInfo  LogLevel = "info"
	LogLevelError LogLevel = "error"
)

type Group struct {
	High int
}
