package srv

import (
	"fmt"
	"io"
	"net/http"

	"github.com/TsvetanMilanov/go-gin-prometheus-middleware/middleware"
	"github.com/TsvetanMilanov/go-simple-di/di"
	"github.com/TsvetanMilanov/go-srv/srv/log"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type appBuilder struct {
	router *gin.Engine

	metricsRouter            http.Handler
	metricsMiddlewareOptions *middleware.Options
	metricsGatherer          prometheus.Gatherer

	appDI             *di.Container
	reqDIConfigurator RequestDIConfiguratorFunc
	appLogger         log.Logger
	err               error
}

func (ab *appBuilder) Initialize(loggerOut io.Writer) MetricsServerConfigurator {
	ab.appDI = di.NewContainer()
	ab.appLogger = log.NewLogger(loggerOut)

	return ab
}

func (ab *appBuilder) EnableMetricsServer(gatherer prometheus.Gatherer, options *middleware.Options) MetricsServerConfigurator {
	ab.metricsGatherer = gatherer
	if ab.metricsGatherer == nil {
		ab.metricsGatherer = prometheus.DefaultGatherer
	}

	ab.metricsMiddlewareOptions = options
	if ab.metricsMiddlewareOptions == nil {
		ab.metricsMiddlewareOptions = new(middleware.Options)
	}

	if gatherer != nil && gatherer != prometheus.DefaultGatherer && ab.metricsMiddlewareOptions == nil {
		return ab.setErr("EnableMetricsServer", "no custom registerer set in the middleware options but custom gatherer was provided", nil)
	}

	ab.metricsRouter = promhttp.HandlerFor(ab.metricsGatherer, promhttp.HandlerOpts{})
	ab.metricsMiddlewareOptions = options

	return ab
}

func (ab *appBuilder) RegisterAppDependencies(registerer DIContainerUser) AppDependenciesResolver {
	err := ab.appDI.Register(
		&di.Dependency{Name: AppLoggerName, Value: ab.appLogger},
		&di.Dependency{Name: AppDIName, Value: ab.appDI},
	)

	if err != nil {
		return ab.setErr("RegisterAppDependencies", "unable to register the default app dependencies", err)
	}

	err = registerer(ab.appDI)
	if err != nil {
		ab.setErr("RegisterAppDependencies", "unable to register the app dependencies", err)
	}

	return ab
}

func (ab *appBuilder) ResolveAppDependencies() ReqDIConfiguratorRegisterer {
	if ab.err != nil {
		return ab
	}

	err := ab.appDI.ResolveAll()
	if err != nil {
		ab.setErr("ResolveAppDependencies", "unable to resolve all app dependencies", err)
	}

	return ab
}

func (ab *appBuilder) RegisterReqDIConfigurator(configurator RequestDIConfiguratorFunc) AppConfigurator {
	if ab.err != nil {
		return ab
	}

	ab.reqDIConfigurator = configurator

	return ab
}

func (ab *appBuilder) ConfigureApp(configurator DIContainerUser) RouterRegisterer {
	if ab.err != nil {
		return ab
	}

	err := configurator(ab.appDI)
	if err != nil {
		ab.setErr("ResolveAppDependencies", "unable to configure the application", err)
	}

	return ab
}

func (ab *appBuilder) RegisterRouter(router *gin.Engine) DefaultMiddlewaresConfigurator {
	if ab.err != nil {
		return ab
	}

	if router == nil {
		return ab.setErr("RegisterRouter", "router should not be nil", nil)
	}

	ab.router = router

	return ab
}

func (ab *appBuilder) AddDefaultMiddlewares() RouterConfigurator {
	if ab.err != nil {
		return ab
	}

	ab.router.Use(recoverMiddlewareFactory(ab.appLogger))
	ab.router.Use(contextPropertiesMiddlewareFactory(ab))

	if ab.reqDIConfigurator != nil {
		ab.router.Use(requestDIConfiguratorMiddlewareFactory(ab))
	}

	if ab.metricsRouter != nil {
		opts := ab.metricsMiddlewareOptions
		if opts == nil {
			opts = new(middleware.Options)
		}

		ab.router.Use(responseLoggerMiddlewareFactory(ab))
		ab.router.Use(middleware.NewWithOptions(opts))
	}

	return ab
}

func (ab *appBuilder) ConfigureRouter(configurator RouterConfiguratorFunc) AppBuilder {
	if ab.err != nil {
		return ab
	}

	err := configurator(ab.router, ab.appDI)
	if err != nil {
		ab.setErr("ConfigureRouter", "unable to configure the router", err)
	}

	return ab
}

func (ab *appBuilder) BuildApp() (App, error) {
	if ab.err != nil {
		return nil, ab.err
	}

	a := newApp(ab)

	return a, nil
}

func (ab *appBuilder) setErr(method, msg string, err error) *appBuilder {
	ab.err = fmt.Errorf("srv: appBuilder: %s: %s", method, msg)

	if err != nil {
		ab.err = fmt.Errorf("%s %s", ab.err, err)
	}

	return ab
}

// NewAppBuilder creates new web application builder.
func NewAppBuilder() AppInitializer {
	ab := new(appBuilder)

	return ab
}
