package master

import (
	"testing"

	"github.com/barbosaigor/nuker/internal/cli"
	"github.com/barbosaigor/nuker/internal/domain/service/master"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
)

func TestModule(t *testing.T) {
	cli.NoLogFile = true

	app := fx.New(
		Module(),
		fx.Invoke(func(m master.Master) {
			assert.NotNil(t, m)
		}),
	)

	assert.Nil(t, app.Err())
}
