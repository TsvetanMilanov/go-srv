package log

import (
	"io"

	"github.com/sirupsen/logrus"
)

// NewLogger returns new instance of the default logger.
func NewLogger(out io.Writer) Logger {
	l := new(logger)

	logrusLogger := logrus.New()
	logrusLogger.SetFormatter(new(logrus.JSONFormatter))
	logrusLogger.SetLevel(logrus.InfoLevel)
	if out != nil {
		logrusLogger.SetOutput(out)
	}

	l.loggerInstance = logrusLogger

	return l
}
