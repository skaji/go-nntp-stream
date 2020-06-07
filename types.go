package nntp

import "net/mail"

type EventType string

var (
	EventTypeArticle EventType = "article"
	EventTypeDebug   EventType = "debug"
	EventTypeInfo    EventType = "info"
	EventTypeError   EventType = "error"
)

type Event struct {
	Type    EventType
	Message string
	Article *Article
}

type Article struct {
	ID     int
	Header mail.Header
	Body   []byte
}

type Group struct {
	High int
}
