package fxrunner

import (
	"github.com/barbosaigor/nuker/pkg/runner"
	"go.uber.org/fx"
)

func Run(opt fx.Option) error {

	app := fx.New(
		opt,
		fx.Invoke(runner.StartRunner),
	)

	if err := app.Err(); err != nil {
		return err
	}

	return nil
}
