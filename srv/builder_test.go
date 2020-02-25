package srv

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/TsvetanMilanov/go-gin-prometheus-middleware/middleware"
	"github.com/TsvetanMilanov/go-simple-di/di"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
)

func createApp() (App, *bytes.Buffer, *prometheus.Registry, error) {
	b := new(bytes.Buffer)
	registry := prometheus.NewRegistry()
	app, err := NewAppBuilder().
		Initialize(b).
		EnableMetricsServer(registry, &middleware.Options{Registry: registry}).
		RegisterAppDependencies(registerAppDependencies).
		ResolveAppDependencies().
		RegisterReqDIConfigurator(registerReqDependencies).
		ConfigureApp(func(appDI *di.Container) error { return nil }).
		RegisterRouter(gin.New()).
		AddDefaultMiddlewares().
		ConfigureRouter(configureRouter).
		BuildApp()

	return app, b, registry, err
}

func assertLogMessage(t *testing.T, b *bytes.Buffer, expected map[string]interface{}, expectedToContain []string) {
	t.Helper()

	logEntry := make(map[string]interface{})
	err := json.Unmarshal(b.Bytes(), &logEntry)
	assert.NoError(t, err)

	for k, v := range expected {
		assert.Equal(t, v, logEntry[k])
	}

	for _, k := range expectedToContain {
		_, ok := logEntry[k]

		assert.True(t, ok)
	}
}

func TestBuilderAllSteps(t *testing.T) {
	app, b, _, err := createApp()
	assert.NoError(t, err)
	assert.NotNil(t, app)

	res := performRequest(app.GetRouter(), http.MethodGet, "/data")

	expected, _ := json.Marshal(dbData)
	assert.Equal(t, strings.TrimSpace(string(expected)), strings.TrimSpace(string(res.Body.Bytes())))
	assertLogMessage(t, b, map[string]interface{}{"level": "info", "msg": "request complete"}, []string{"time", "traceId"})

	metricsRes := performRequest(app.GetMetricsRouter(), http.MethodGet, "/metrics")
	assert.Contains(t, string(metricsRes.Body.Bytes()), "http_request_duration_seconds_count{method=\"GET\",path=\"/data\",status_code=\"200\"} 1")
}
