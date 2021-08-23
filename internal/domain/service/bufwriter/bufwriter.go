package bufwriter

import (
	"bufio"
	"errors"
	"io"
)

type bufWriter struct {
	location string
	buffer   *bufio.Writer
}

func New(writer io.Writer, location string) (BufWriter, error) {
	if location == "" {
		return nil, errors.New("bufWriter: nil location")
	}

	if writer == nil {
		return nil, errors.New("bufWriter: nil writer")
	}

	return &bufWriter{
		location: location,
		buffer:   bufio.NewWriter(writer),
	}, nil
}

func (bw *bufWriter) Write(data []byte) (n int, err error) {
	n, err = bw.buffer.Write([]byte(string(data) + "\n"))
	if err != nil {
		return
	}

	bw.buffer.Flush()

	return
}

func (bw *bufWriter) Location() string {
	return bw.location
}
