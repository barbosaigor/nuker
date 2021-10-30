package worker

import (
	"testing"

	"github.com/barbosaigor/nuker/internal/domain/service/worker"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
)

func TestModule(t *testing.T) {
	app := fx.New(
		Module("w-test", "http://master.io", 1),
		fx.Invoke(func(w worker.Worker) {
			assert.NotNil(t, w)
		}),
	)

	assert.Nil(t, app.Err())
}
