package log

// Logger represents the logger wrapper.
type Logger interface {
	Info(a ...interface{})
	Debug(a ...interface{})
	Warning(a ...interface{})
	Error(a ...interface{})

	SetLevel(level string)
	AddFields(fields map[string]interface{})
	CreateChild() Logger
}
