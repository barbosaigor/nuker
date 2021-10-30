package orchestrator

import (
	"testing"

	"github.com/barbosaigor/nuker/internal/domain/service/orchestrator"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
)

func TestModule(t *testing.T) {
	app := fx.New(
		Module(),
		fx.Invoke(func(s orchestrator.Orchestrator) {
			assert.NotNil(t, s)
		}),
	)

	assert.Nil(t, app.Err())
}
