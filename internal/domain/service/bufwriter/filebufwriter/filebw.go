package filebufwriter

import (
	"fmt"
	"os"
	"time"

	"github.com/barbosaigor/nuker/internal/domain/service/bufwriter"
)

// New creates a bufWriter with file provider
func New(fName string) (bufwriter.BufWriter, error) {
	if fName == "" {
		fName = fileName(fName)
	}

	file, err := os.Create(fName)
	if err != nil {
		return nil, err
	}

	return bufwriter.New(file, fName)
}

func fileName(fName string) string {
	return fmt.Sprintf("nuker-%v.jsonl", time.Now().Unix())
}
