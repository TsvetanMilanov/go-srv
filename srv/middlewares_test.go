package srv

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRecoverMiddleware(t *testing.T) {
	app, b, _, err := createApp(nil)
	assert.NoError(t, err)
	assert.NotNil(t, app)

	res := performRequest(app.GetRouter(), http.MethodGet, "/panic")

	assert.Equal(t, http.StatusInternalServerError, res.Code)
	expected, _ := json.Marshal(map[string]string{"message": "Internal server error"})
	assert.Equal(t, strings.TrimSpace(string(expected)), strings.TrimSpace(string(res.Body.Bytes())))
	assertLogMessage(t, b, []map[string]interface{}{
		{"level": "error"},
		{"level": "info", "msg": "request complete"},
	}, [][]string{
		{"time", "traceId"},
		{"time", "traceId"},
	})

	metricsRes := performRequest(app.GetMetricsRouter(), http.MethodGet, "/metrics")
	assert.Contains(t, string(metricsRes.Body.Bytes()), "http_request_duration_seconds_count{method=\"GET\",path=\"/panic\",status_code=\"500\"} 1")
}
