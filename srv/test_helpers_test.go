package srv

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/TsvetanMilanov/go-gin-prometheus-middleware/middleware"
	"github.com/TsvetanMilanov/go-simple-di/di"
	"github.com/TsvetanMilanov/go-srv/srv/log"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
)

const (
	customMetricsAddr = ":33333"
)

var (
	dbData = map[string]interface{}{"key": "value"}
)

func getFreeAddr() string {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		panic(err)
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		panic(err)
	}

	defer l.Close()

	port := l.Addr().(*net.TCPAddr).Port

	return fmt.Sprintf(":%d", port)
}

type testDBClient struct{}

func (tcl *testDBClient) GetData() map[string]interface{} {
	return dbData
}

type testDataService struct {
	DBClient        *testDBClient   `di:""`
	RequestLogger   log.Logger      `di:"name=reqLogger"`
	TraceIDProvider TraceIDProvider `di:""`
}

func (ts *testDataService) GetData() map[string]interface{} {
	data := ts.DBClient.GetData()

	return data
}

type testController struct{}

func (tc *testController) Get(c *gin.Context) {
	dataSvc := new(testDataService)
	reqDI, _ := c.Get("reqDi")

	err := reqDI.(*di.Container).Resolve(dataSvc)
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, dataSvc.GetData())
}

func (tc *testController) Panic(c *gin.Context) {
	panic(errors.New("controller panic"))
}

func registerAppDependencies(container *di.Container) error {
	err := container.Register(
		&di.Dependency{Value: new(testController)},
	)

	return err
}

func registerReqDependencies(req *http.Request, container *di.Container) error {
	err := container.Register(
		&di.Dependency{Value: new(testDataService)},
		&di.Dependency{Value: new(testDBClient)},
	)

	return err
}

func configureRouter(router *gin.Engine, appDI *di.Container) error {
	tc := new(testController)
	err := appDI.Resolve(tc)
	if err != nil {
		return err
	}

	router.GET("/data", tc.Get)
	router.GET("/panic", tc.Panic)

	return nil
}

func performRequest(r http.Handler, method, path string) *httptest.ResponseRecorder {
	req, err := http.NewRequest(method, path, nil)
	if err != nil {
		panic(err)
	}

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	return w
}

func createApp(promMiddlewareOpts *middleware.Options) (App, *bytes.Buffer, *prometheus.Registry, error) {
	registry := prometheus.NewRegistry()

	if promMiddlewareOpts == nil {
		promMiddlewareOpts = &middleware.Options{Registry: registry}
	}

	b := new(bytes.Buffer)
	msc := NewAppBuilder().
		Initialize(b).
		EnableMetricsServer(registry, promMiddlewareOpts)

	app, err := configureMetricsServerConfigurator(msc)

	return app, b, registry, err
}

func configureMetricsServerConfigurator(msc MetricsServerConfigurator) (App, error) {
	app, err := msc.
		RegisterAppDependencies(registerAppDependencies).
		ResolveAppDependencies().
		RegisterReqDIConfigurator(registerReqDependencies).
		ConfigureApp(func(appDI *di.Container) error { return nil }).
		RegisterRouter(gin.New()).
		AddDefaultMiddlewares().
		ConfigureRouter(configureRouter).
		BuildApp()

	return app, err
}

func assertLogMessage(t *testing.T, b *bytes.Buffer, expected []map[string]interface{}, expectedToContain [][]string) {
	t.Helper()

	allLines := strings.Split(strings.TrimSpace(string(b.Bytes())), "\n")
	for i, line := range allLines {
		logEntry := make(map[string]interface{})

		err := json.Unmarshal([]byte(line), &logEntry)
		assert.NoError(t, err)

		for k, v := range expected[i] {
			assert.Equal(t, v, logEntry[k])
		}

		for _, k := range expectedToContain[i] {
			_, ok := logEntry[k]

			assert.True(t, ok)
		}
	}
}
