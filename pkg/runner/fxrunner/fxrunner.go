package fxrunner

import (
	"context"

	"github.com/barbosaigor/nuker/pkg/runner"
	"go.uber.org/fx"
)

func Run(opt fx.Option) error {
	app := fx.New(
		opt,
		fx.Invoke(startRunner),
	)

	if err := app.Err(); err != nil {
		return err
	}

	return nil
}

func startRunner(ctx context.Context, r runner.Runner) error {
	return r.Run(ctx)
}
