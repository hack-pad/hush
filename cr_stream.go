package hush

import (
	"io"
)

type carriageReturnWriter struct {
	io.Writer
}

func newCarriageReturnWriter(dest io.Writer) (io.Writer, error) {
	return &carriageReturnWriter{dest}, nil
}

func (c *carriageReturnWriter) Write(buf []byte) (n int, err error) {
	for _, b := range buf {
		_, err = c.Writer.Write([]byte{b})
		if err != nil {
			return
		}
		if b == '\n' {
			_, err = c.Writer.Write([]byte{'\r'})
			if err != nil {
				return
			}
		}
		n++
	}
	return
}
