package srv

import (
	"io"

	"github.com/TsvetanMilanov/go-gin-prometheus-middleware/middleware"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

// AppInitializer provides method for initializing
// the web application.
type AppInitializer interface {
	Initialize(loggerOut io.Writer) MetricsServerConfigurator
}

// MetricsServerConfigurator provides methods for configuring the metrics server.
type MetricsServerConfigurator interface {
	AppDependenciesRegisterer

	// This step is not required.
	EnableMetricsServer(options *middleware.Options) MetricsServerConfigurator
	// This step is not required.
	// If no registry is set the server will use the default prometheus registry.
	SetMetricsRegistry(registry *prometheus.Registry) MetricsServerConfigurator
	// This step is not required.
	// Defaults to :80
	SetMetricsServerAddr(addr string) MetricsServerConfigurator
}

// AppDependenciesRegisterer provides method for registering the app dependencies.
type AppDependenciesRegisterer interface {
	RegisterAppDependencies(registerer DIContainerUser) AppDependenciesResolver
}

// AppDependenciesResolver provides method for resolving all app dependencies.
type AppDependenciesResolver interface {
	ResolveAppDependencies() ReqDIConfiguratorRegisterer
}

// ReqDIConfiguratorRegisterer provides method for configuring the request di.
type ReqDIConfiguratorRegisterer interface {
	AppConfigurator
	RegisterReqDIConfigurator(configurator RequestDIConfiguratorFunc) AppConfigurator
}

// AppConfigurator provides method for configuring the application.
// This method provides access to the resolved app dependencies.
// Use it to set the log level for example.
type AppConfigurator interface {
	RouterRegisterer
	// This step is not required.
	ConfigureApp(configurator DIContainerUser) RouterRegisterer
}

// RouterRegisterer provides method for registering the app router.
type RouterRegisterer interface {
	RegisterRouter(router *gin.Engine) DefaultMiddlewaresConfigurator
}

// DefaultMiddlewaresConfigurator provides method which adds the default middlewares
// to the router.
type DefaultMiddlewaresConfigurator interface {
	RouterConfigurator
	// This step is not required.
	AddDefaultMiddlewares() RouterConfigurator
}

// RouterConfigurator provides method for configuring the router.
type RouterConfigurator interface {
	ConfigureRouter(configurator RouterConfiguratorFunc) AppBuilder
}

// AppBuilder provides method to build the web application.
type AppBuilder interface {
	BuildApp() (App, error)
}
