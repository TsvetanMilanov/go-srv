package log

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createLogger() (Logger, *bytes.Buffer) {
	b := new(bytes.Buffer)
	l := NewLogger(b)

	return l, b
}

func assertLogEntry(
	t *testing.T,
	writer *bytes.Buffer,
	expectedLogLevel,
	expectedMsg string,
	expectedAdditionalFields map[string]interface{},
	expectedMissingFields map[string]interface{},
) {
	entry := make(map[string]interface{})
	err := json.Unmarshal(writer.Bytes(), &entry)

	assert.NoError(t, err)
	assert.Equal(t, expectedLogLevel, entry["level"])
	assert.Equal(t, expectedMsg, entry["msg"])

	if len(expectedAdditionalFields) > 0 {
		for k, v := range expectedAdditionalFields {
			actual, ok := entry[k]
			assert.True(t, ok)
			assert.Equal(t, v, actual)
		}
	}

	if len(expectedMissingFields) > 0 {
		for k := range expectedMissingFields {
			_, ok := entry[k]
			assert.False(t, ok)
		}
	}
}

func TestLoggerInfo(t *testing.T) {
	l, b := createLogger()

	l.Info("test")

	assertLogEntry(t, b, "info", "test", nil, nil)
}

func TestLoggerDebug(t *testing.T) {
	l, b := createLogger()
	l.SetLevel("debug")

	l.Debug("test")

	assertLogEntry(t, b, "debug", "test", nil, nil)
}

func TestLoggerWarning(t *testing.T) {
	l, b := createLogger()

	l.Warning("test")

	assertLogEntry(t, b, "warning", "test", nil, nil)
}

func TestLoggerError(t *testing.T) {
	l, b := createLogger()

	l.Error("test")

	assertLogEntry(t, b, "error", "test", nil, nil)
}

func TestLoggerAddFields(t *testing.T) {
	l, b := createLogger()
	customFields := map[string]interface{}{"custom": "field"}
	l.AddFields(customFields)

	l.Info("test")

	assertLogEntry(t, b, "info", "test", customFields, nil)
}

func TestLoggerCreateChild(t *testing.T) {
	l, b := createLogger()
	customFields := map[string]interface{}{"custom": "field"}
	l.AddFields(customFields)

	child := l.CreateChild()
	childFields := map[string]interface{}{"child": "value"}
	child.AddFields(childFields)

	l.Info("test")

	assertLogEntry(t, b, "info", "test", customFields, childFields)

	b.Reset()

	child.Info("childTest")

	assertLogEntry(t, b, "info", "childTest", customFields, nil)
	assertLogEntry(t, b, "info", "childTest", childFields, nil)
}

func TestLoggerSetLevelPanic(t *testing.T) {
	l := NewLogger(os.Stdout)

	assert.Panics(t, func() { l.SetLevel("invalid-level") })
}
