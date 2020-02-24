package log

import (
	"github.com/sirupsen/logrus"
)

type logger struct {
	loggerInstance *logrus.Logger
	fields         map[string]interface{}
}

func (l *logger) Info(a ...interface{}) {
	l.log(logrus.InfoLevel, a...)
}

func (l *logger) Debug(a ...interface{}) {
	l.log(logrus.DebugLevel, a...)
}

func (l *logger) Warning(a ...interface{}) {
	l.log(logrus.WarnLevel, a...)
}

func (l *logger) Error(a ...interface{}) {
	l.log(logrus.ErrorLevel, a...)
}

func (l *logger) SetLevel(level string) {
	parsedLevel, err := logrus.ParseLevel(level)
	if err != nil {
		panic(err)
	}

	l.loggerInstance.SetLevel(parsedLevel)
}

func (l *logger) AddFields(fields map[string]interface{}) {
	if l.fields == nil {
		l.fields = make(map[string]interface{})
	}

	for k, v := range fields {
		l.fields[k] = v
	}
}

func (l *logger) CreateChild() Logger {
	child := l.clone()

	return child
}

func (l *logger) log(level logrus.Level, a ...interface{}) {
	l.loggerInstance.WithFields(l.fields).Log(level, a...)
}

func (l *logger) clone() *logger {
	clone := new(logger)
	clone.fields = make(map[string]interface{})
	for k, v := range l.fields {
		clone.fields[k] = v
	}

	clone.loggerInstance = logrus.New()
	clone.loggerInstance.ExitFunc = l.loggerInstance.ExitFunc
	clone.loggerInstance.Formatter = l.loggerInstance.Formatter
	clone.loggerInstance.Hooks = l.loggerInstance.Hooks
	clone.loggerInstance.Level = l.loggerInstance.Level
	clone.loggerInstance.Out = l.loggerInstance.Out
	clone.loggerInstance.ReportCaller = l.loggerInstance.ReportCaller

	return clone
}
