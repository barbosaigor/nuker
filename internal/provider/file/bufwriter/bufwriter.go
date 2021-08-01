package bufwriter

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/barbosaigor/nuker/internal/domain/repository"
)

type bufWriter struct {
	fileName string
	buffer   *bufio.Writer
}

func New(fileName string) (repository.BufWriter, error) {
	if fileName == "" {
		fileName = fmt.Sprintf("nuker-%v.jsonl", time.Now().Unix())
	}
	file, err := os.Create(fileName)

	return &bufWriter{
		fileName: fileName,
		buffer:   bufio.NewWriter(file),
	}, err
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
	return bw.fileName
}
