package srv

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/TsvetanMilanov/go-simple-di/di"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestBuilderAllSteps(t *testing.T) {
	b := new(bytes.Buffer)
	app, err := NewAppBuilder().
		Initialize(b).
		RegisterAppDependencies(func(container *di.Container) error {
			return nil
		}).
		ResolveAppDependencies().
		RegisterReqDIConfigurator(func(req *http.Request, reqDI *di.Container) error {
			return nil
		}).
		ConfigureApp(func(appDI *di.Container) error {
			return nil
		}).
		RegisterRouter(gin.New()).
		AddDefaultMiddlewares().
		ConfigureRouter(func(router *gin.Engine, abDI, reqDI *di.Container) error {
			return nil
		}).
		BuildApp()

	assert.NoError(t, err)
	assert.NotNil(t, app)
}
