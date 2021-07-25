package runner

import (
	"context"

	"github.com/rs/zerolog/log"
)

func StartRunner(ctx context.Context, r Runner) error {
	if err := r.Run(ctx); err != nil {
		log.Ctx(ctx).Error().Err(err)
		return err
	}

	return nil
}
