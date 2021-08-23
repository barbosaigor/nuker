package bufwriter

import "io"

type BufWriter interface {
	io.Writer
	Location() string
}
