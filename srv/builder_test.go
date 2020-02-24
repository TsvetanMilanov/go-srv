package srv

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/TsvetanMilanov/go-simple-di/di"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
)

func TestBuilderAllSteps(t *testing.T) {
	b := new(bytes.Buffer)
	registry := prometheus.NewRegistry()
	app, err := NewAppBuilder().
		Initialize(b).
		EnableMetricsServer(nil).
		SetMetricsRegistry(registry).
		SetMetricsServerAddr(customMetricsAddr).
		RegisterAppDependencies(registerAppDependencies).
		ResolveAppDependencies().
		RegisterReqDIConfigurator(registerReqDependencies).
		ConfigureApp(func(appDI *di.Container) error { return nil }).
		RegisterRouter(gin.New()).
		AddDefaultMiddlewares().
		ConfigureRouter(configureRouter).
		BuildApp()

	assert.NoError(t, err)
	assert.NotNil(t, app)

	res := performRequest(app.GetRouter(), http.MethodGet, "/data")

	expected, _ := json.Marshal(dbData)
	assert.Equal(t, strings.TrimSpace(string(expected)), strings.TrimSpace(string(res.Body.Bytes())))
}
