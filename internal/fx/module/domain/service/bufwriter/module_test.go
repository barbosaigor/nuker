package bufwriter

import (
	"testing"

	"github.com/barbosaigor/nuker/internal/cli"
	"github.com/barbosaigor/nuker/internal/domain/service/bufwriter"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
)

func TestModuleWithNoLogFile(t *testing.T) {
	cli.NoLogFile = true

	app := fx.New(
		Module(),
		fx.Invoke(func(b bufwriter.BufWriter) {
			assert.NotNil(t, b)
		}),
	)

	assert.Nil(t, app.Err())
}
