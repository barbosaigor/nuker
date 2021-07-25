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

func New() (repository.BufWriter, error) {
	fName := fmt.Sprintf("nuker-%v.jsonl", time.Now().Unix())
	file, err := os.Create(fName)

	return &bufWriter{
		fileName: fName,
		buffer:   bufio.NewWriter(file),
	}, err
}

const newLine = byte('\n')

func (bw *bufWriter) Write(data []byte) (n int, err error) {
	tmpData := make([]byte, len(data)+1)
	copy(tmpData, data)
	return bw.buffer.Write(append(tmpData, newLine))
}

func (bw *bufWriter) Location() string {
	return bw.fileName
}
