package publisher

import (
	"testing"

	"github.com/barbosaigor/nuker/internal/domain/service/publisher"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
)

func TestModule(t *testing.T) {
	app := fx.New(
		Module(),
		fx.Invoke(func(p publisher.Publisher) {
			assert.NotNil(t, p)
		}),
	)

	assert.Nil(t, app.Err())
}
