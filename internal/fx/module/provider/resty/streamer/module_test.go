package streamer

import (
	"testing"

	"github.com/barbosaigor/nuker/internal/domain/repository"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
)

func TestModule(t *testing.T) {
	app := fx.New(
		Module(),
		fx.Invoke(func(s repository.Streamer) {
			assert.NotNil(t, s)
		}),
	)

	assert.Nil(t, app.Err())
}
