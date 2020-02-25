package srv

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/TsvetanMilanov/go-gin-prometheus-middleware/middleware"
	"github.com/TsvetanMilanov/go-simple-di/di"
	"github.com/TsvetanMilanov/go-srv/srv/log"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestBuilderAllSteps(t *testing.T) {
	app, b, _, err := createApp(nil)
	assert.NoError(t, err)
	assert.NotNil(t, app)

	res := performRequest(app.GetRouter(), http.MethodGet, "/data")

	assert.Equal(t, http.StatusOK, res.Code)
	expected, _ := json.Marshal(dbData)
	assert.Equal(t, strings.TrimSpace(string(expected)), strings.TrimSpace(string(res.Body.Bytes())))
	assertLogMessage(t, b, []map[string]interface{}{{"level": "info", "msg": "request complete"}}, [][]string{{"time", "traceId"}})

	metricsRes := performRequest(app.GetMetricsRouter(), http.MethodGet, "/metrics")
	assert.Contains(t, string(metricsRes.Body.Bytes()), "http_request_duration_seconds_count{method=\"GET\",path=\"/data\",status_code=\"200\"} 1")
}

func TestBuilderRegistererGathererMismatch(t *testing.T) {
	_, _, _, err := createApp(new(middleware.Options))
	assert.EqualError(t, err, "srv: appBuilder: EnableMetricsServer: no custom registerer set in the middleware options but custom gatherer was provided")
}

func TestBuilderDefaultMetricsGatherer(t *testing.T) {
	b := new(bytes.Buffer)

	app, err := configureMetricsServerConfigurator(NewAppBuilder().
		Initialize(b).
		EnableMetricsServer(nil, nil))

	assert.NoError(t, err)

	metricsRes := performRequest(app.GetMetricsRouter(), http.MethodGet, "/metrics")
	assert.Contains(t, string(metricsRes.Body.Bytes()), "# HELP go_gc_duration_seconds A summary of the pause duration of garbage collection cycles.")
}

func TestBuilderAppDependenciesRegistreationError(t *testing.T) {
	_, err := NewAppBuilder().
		Initialize(os.Stdout).
		RegisterAppDependencies(func(appDI *di.Container) error { return errors.New("test err") }).
		ResolveAppDependencies().
		RegisterRouter(gin.New()).
		ConfigureRouter(configureRouter).
		BuildApp()

	assert.EqualError(t, err, "srv: appBuilder: RegisterAppDependencies: unable to register the app dependencies test err")
}

func TestBuilderResolveAppDependencies(t *testing.T) {
	_, err := NewAppBuilder().
		Initialize(os.Stdout).
		RegisterAppDependencies(func(appDI *di.Container) error {
			type unresolvable struct {
				Logger log.Logger `di:"name=unresolvable"`
			}

			return appDI.Register(&di.Dependency{Value: new(unresolvable)})
		}).
		ResolveAppDependencies().
		RegisterRouter(gin.New()).
		ConfigureRouter(configureRouter).
		BuildApp()

	assert.EqualError(t, err, "srv: appBuilder: ResolveAppDependencies: unable to resolve all app dependencies [*srv.unresolvable] unable to find registered dependency: Logger")
}

func TestBuilderAppConfigurationError(t *testing.T) {
	_, err := NewAppBuilder().
		Initialize(os.Stdout).
		RegisterAppDependencies(registerAppDependencies).
		ResolveAppDependencies().
		ConfigureApp(func(appDI *di.Container) error { return errors.New("test err") }).
		RegisterRouter(gin.New()).
		ConfigureRouter(configureRouter).
		BuildApp()

	assert.EqualError(t, err, "srv: appBuilder: ConfigureApp: unable to configure the application test err")
}

func TestBuilderRegisterRouterError(t *testing.T) {
	_, err := NewAppBuilder().
		Initialize(os.Stdout).
		RegisterAppDependencies(registerAppDependencies).
		ResolveAppDependencies().
		RegisterRouter(nil).
		ConfigureRouter(configureRouter).
		BuildApp()

	assert.EqualError(t, err, "srv: appBuilder: RegisterRouter: router should not be nil")
}

func TestBuilderConfigureRouterError(t *testing.T) {
	_, err := NewAppBuilder().
		Initialize(os.Stdout).
		RegisterAppDependencies(registerAppDependencies).
		ResolveAppDependencies().
		RegisterRouter(gin.New()).
		ConfigureRouter(func(r *gin.Engine, appDI *di.Container) error { return errors.New("test err") }).
		BuildApp()

	assert.EqualError(t, err, "srv: appBuilder: ConfigureRouter: unable to configure the router test err")
}
