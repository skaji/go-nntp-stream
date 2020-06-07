# go nntp stream

```go
package main

import (
	"context"
	"log"
	"time"

	nntp "github.com/skaji/go-nntp-stream"
)

func main() {
	config := nntp.StreamConfig{
		Addr:    "nntp.perl.org:119",
		Group:   "perl.cpan.uploads",
		Tick:    30 * time.Second,
		Timeout: 20 * time.Second,
	}

	stream := nntp.Stream(context.Background(), config)
	for event := range stream {
		switch event.Type {
		case nntp.EventTypeArticle:
			article := event.Article
			log.Println(article.ID, article.Header.Get("Subject"))
		default:
			log.Println(event.Type, event.Message)
		}
	}
}
```

# See also

https://github.com/dustin/go-nntp

# Author

Shoichi Kaji

# License

MIT
