package nntp_test

import (
	"context"
	"testing"

	nntp "github.com/skaji/go-nntp-stream"
)

func TestClient(t *testing.T) {
	c, err := nntp.Connect(context.Background(), "nntp.perl.org:119")
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	g, err := c.Group("perl.cpan.uploads")
	if err != nil {
		t.Fatal(err)
	}
	a, err := c.Article(g.High)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("ID", a.ID)
	for k, v := range a.Header {
		t.Log(k, v)
	}
	t.Log(string(a.Body))
}
