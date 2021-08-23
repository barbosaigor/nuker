package nopbufwriter

import (
	"github.com/barbosaigor/nuker/internal/domain/service/bufwriter"
)

type nopBufWriter struct{}

// New creates a dumb bufwriter implementation
func New() bufwriter.BufWriter {
	return &nopBufWriter{}
}

func (bw *nopBufWriter) Write(data []byte) (n int, err error) {
	return len(data), nil
}

func (bw *nopBufWriter) Location() string {
	return "nop"
}
