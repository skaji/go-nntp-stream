package nntp

import (
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"net/mail"
	"net/textproto"
	"strconv"
	"strings"
	"time"
)

type Client struct {
	baseConn net.Conn
	conn     *textproto.Conn
}

func Connect(ctx context.Context, addr string) (*Client, error) {
	conn, err := (&net.Dialer{}).DialContext(ctx, "tcp", addr)
	if err != nil {
		return nil, err
	}

	c := &Client{baseConn: conn, conn: textproto.NewConn(conn)}
	cancel := c.SetDeadline(ctx)
	_, _, err = c.conn.ReadCodeLine(200)
	cancel()
	if err != nil {
		c.Close()
		return nil, err
	}
	return c, nil
}

var (
	noDeadline         = time.Time{}
	timeoutImmediately = time.Unix(1, 0)
)

func (c *Client) SetDeadline(ctx context.Context) func() {
	if deadline, ok := ctx.Deadline(); ok {
		c.baseConn.SetDeadline(deadline)
	} else {
		c.baseConn.SetDeadline(noDeadline)
	}
	stop := make(chan struct{})
	done := make(chan struct{})
	go func() {
		select {
		case <-ctx.Done():
			c.baseConn.SetDeadline(timeoutImmediately)
		case <-stop:
		}
		close(done)
	}()
	return func() {
		close(stop)
		<-done
		c.baseConn.SetDeadline(noDeadline)
	}
}

func (c *Client) Close() error {
	return c.baseConn.Close()
}

func (c *Client) Group(name string) (*Group, error) {
	_, msg, err := c.command("GROUP "+name, 211)
	if err != nil {
		return nil, err
	}
	// msg = 257662 1 257662 perl.cpan.uploads
	parts := strings.Split(msg, " ")
	if len(parts) != 4 {
		return nil, fmt.Errorf("failed to parse message '%s'", msg)
	}
	high, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse high count '%s'", parts[2])
	}
	return &Group{High: int(high)}, nil
}

func (c *Client) Article(id int) (*Article, error) {
	if err := c.conn.PrintfLine("ARTICLE %d", id); err != nil {
		return nil, err
	}
	if _, _, err := c.conn.ReadCodeLine(220); err != nil {
		return nil, err
	}
	msg, err := mail.ReadMessage(c.conn.DotReader())
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(msg.Body)
	if err != nil {
		return nil, err
	}
	return &Article{
		ID:     id,
		Body:   body,
		Header: msg.Header,
	}, nil
}

func (c *Client) command(cmd string, expectCode int) (int, string, error) {
	err := c.conn.PrintfLine(cmd)
	if err != nil {
		return 0, "", err
	}
	return c.conn.ReadCodeLine(expectCode)
}
