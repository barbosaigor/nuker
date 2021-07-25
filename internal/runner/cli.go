package runner

import (
	"log"
	"os"
)

var cfgFileName string

func init() {
	if len(os.Args) < 2 {
		log.Fatalln("missing plan file path")
	}

	cfgFileName = os.Args[1]
}
