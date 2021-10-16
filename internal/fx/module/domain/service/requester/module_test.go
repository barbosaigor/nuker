package requester

import (
	"testing"

	"github.com/barbosaigor/nuker/internal/domain/service/requester"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
)

func TestModule(t *testing.T) {
	app := fx.New(
		Module(),
		fx.Invoke(func(r requester.Requester) {
			assert.NotNil(t, r)
		}),
	)

	assert.Nil(t, app.Err())
}
