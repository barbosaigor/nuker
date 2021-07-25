package main

import (
	"log"

	"github.com/barbosaigor/nuker/internal/fx/module/runner"
	"github.com/barbosaigor/nuker/pkg/runner/fxrunner"
)

func main() {
	err := fxrunner.Run(runner.Module())
	if err != nil {
		log.Fatal(err)
	}
}
