package repository

import "io"

type BufWriter interface {
	io.Writer
	Location() string
}
