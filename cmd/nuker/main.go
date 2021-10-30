package main

import (
	"github.com/barbosaigor/nuker/internal/fx/module/runner"
	"github.com/barbosaigor/nuker/pkg/runner/fxrunner"
	"github.com/sirupsen/logrus"
)

func main() {
	err := fxrunner.Run(runner.Module())
	if err != nil {
		logrus.Fatal(err)
	}
}
