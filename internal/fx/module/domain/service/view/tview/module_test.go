package tview

import (
	"testing"

	"github.com/barbosaigor/nuker/internal/domain/service/view"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
)

func TestModuleWithNoLogFile(t *testing.T) {
	app := fx.New(
		Module(),
		fx.Invoke(func(v view.View) {
			assert.NotNil(t, v)
		}),
	)

	assert.Nil(t, app.Err())
}
