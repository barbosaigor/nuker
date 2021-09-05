package runner

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
)

func TestModule(t *testing.T) {
	app := fx.New(Module())
	assert.Nil(t, app.Err())
}
