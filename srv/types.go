package srv

import (
	"net/http"

	"github.com/TsvetanMilanov/go-graceful-server-shutdown/gss"
	"github.com/TsvetanMilanov/go-simple-di/di"
	"github.com/TsvetanMilanov/go-srv/srv/log"
	"github.com/gin-gonic/gin"
)

// App provides methods for working with the web app.
type App interface {
	Start(srvSettings, metricsSrvSettings *gss.Settings) error
	GetRouter() http.Handler
	GetMetricsRouter() http.Handler
	GetLogger() log.Logger
}

// TraceIDProvider provides methods for acquiring trace id.
type TraceIDProvider interface {
	GetTraceID() string
}

// DIContainerUser is a function which receives a di container and
// provides implementation which uses it.
type DIContainerUser = func(container *di.Container) error

// RequestDIConfiguratorFunc function which can be used to configure the request di.
type RequestDIConfiguratorFunc = func(req *http.Request, reqDI *di.Container) error

// RouterConfiguratorFunc function which can be used to configure the router.
type RouterConfiguratorFunc = func(router *gin.Engine, appDI *di.Container) error
