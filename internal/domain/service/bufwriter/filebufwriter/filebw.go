package filebufwriter

import (
	"fmt"
	"os"
	"time"

	"github.com/barbosaigor/nuker/internal/domain/service/bufwriter"
)

// New creates a bufWriter with file provider
func New(fileName string) (bufwriter.BufWriter, error) {
	if fileName == "" {
		fileName = fmt.Sprintf("nuker-%v.jsonl", time.Now().Unix())
	}

	file, err := os.Create(fileName)
	if err != nil {
		return nil, err
	}

	return bufwriter.New(file, fileName)
}
