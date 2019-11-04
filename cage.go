package cage

import (
	"bytes"
	"io"
	"os"
	"strings"
	"time"
)

type container struct {
	backupStdout *os.File
	writerStdout *os.File
	backupStderr *os.File
	writerStderr *os.File

	data         string
	channel      chan string

	Data []string
}

func Start() *container {
	rStdout, wStdout, _ := os.Pipe()
	rStderr, wStderr, _ := os.Pipe()
	c := &container{
		backupStdout: os.Stdout,
		writerStdout: wStdout,

		backupStderr: os.Stderr,
		writerStderr: wStderr,

		channel: make(chan string),
	}
	os.Stdout = c.writerStdout
	os.Stderr = c.writerStderr

	go func(out chan string, readerStdout *os.File, readerStderr *os.File) {
		var bufStdout bytes.Buffer
		_, _ = io.Copy(&bufStdout, readerStdout)
		if bufStdout.Len() > 0 {
			out <- bufStdout.String()
		}

		var bufStderr bytes.Buffer
		_, _ = io.Copy(&bufStderr, readerStderr)
		if bufStderr.Len() > 0 {
			out <- bufStderr.String()
		}
	}(c.channel, rStdout, rStderr)

	go func(c *container) {
		for {
			select {
			case out := <-c.channel:
				c.data += out
			}
		}
	}(c)

	return c
}

func Stop(c *container) {
	_ = c.writerStdout.Close()
	_ = c.writerStderr.Close()
	time.Sleep(10 * time.Millisecond)

	os.Stdout = c.backupStdout
	os.Stderr = c.backupStderr

	c.Data = strings.Split(c.data, "\n")
	if c.Data[len(c.Data)-1] == "" {
		c.Data = c.Data[:len(c.Data)-1]
	}
}
