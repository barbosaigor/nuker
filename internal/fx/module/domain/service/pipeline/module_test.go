package pipeline

import (
	"testing"

	"github.com/barbosaigor/nuker/internal/cli"
	"github.com/barbosaigor/nuker/internal/domain/service/pipeline"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
)

func TestModule(t *testing.T) {
	cli.NoLogFile = true

	app := fx.New(
		Module(),
		fx.Invoke(func(p pipeline.Pipeline) {
			assert.NotNil(t, p)
		}),
	)

	assert.Nil(t, app.Err())
}
