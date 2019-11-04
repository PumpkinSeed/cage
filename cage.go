package cage

import (
	"bytes"
	"io"
	"os"
	"strings"
	"time"
)

type container struct {
	backup  *os.File
	writer  *os.File
	data    string
	channel chan string

	Data []string
}

func Start() *container {
	r, w, _ := os.Pipe()
	c := &container{
		backup: os.Stdout,
		writer: w,

		channel: make(chan string),
	}
	os.Stdout = w

	go func(out chan string, reader *os.File) {
		for {
			var buf bytes.Buffer
			_, _ = io.Copy(&buf, r)
			if buf.Len() > 0 {
				out <- buf.String()
			}
		}
	}(c.channel, r)

	go func(c *container) {
		for {
			select {
			case out := <-c.channel:
				c.data = out
			}
		}
	}(c)

	return c
}

func Stop(c *container) {
	_ = c.writer.Close()
	time.Sleep(10 * time.Millisecond)

	os.Stdout = c.backup

	c.Data = strings.Split(c.data, "\n")
	c.Data = c.Data[:len(c.Data)-1]
}
