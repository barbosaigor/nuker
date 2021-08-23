package runner

import (
	log "github.com/sirupsen/logrus"
	"go.uber.org/fx"
)

type fxLogger struct {
	logger *log.Logger
}

func newFxLooger(logger *log.Logger) fx.Printer {
	return &fxLogger{
		logger: logger,
	}
}

func (l fxLogger) Printf(format string, v ...interface{}) {
	l.logger.Tracef(format, v...)
}
