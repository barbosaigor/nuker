package bufwriter

import (
	"github.com/barbosaigor/nuker/internal/domain/repository"
)

type nopBufWriter struct{}

func New() repository.BufWriter {
	return &nopBufWriter{}
}

func (bw *nopBufWriter) Write(data []byte) (n int, err error) {
	return len(data), nil
}

func (bw *nopBufWriter) Location() string {
	return "nop"
}
