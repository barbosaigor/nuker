package streamer

import (
	"github.com/barbosaigor/nuker/internal/provider/resty/streamer"
	"github.com/go-resty/resty/v2"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Options(
		fx.Provide(
			func() *resty.Client {
				client := resty.New()
				client.SetRetryCount(3)
				client.SetLogger(nopLogger{})
				return client
			},
			streamer.New,
		),
	)
}

type nopLogger struct{}

func (nopLogger) Errorf(format string, v ...interface{}) {}
func (nopLogger) Warnf(format string, v ...interface{})  {}
func (nopLogger) Debugf(format string, v ...interface{}) {}
