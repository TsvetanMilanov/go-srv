package srv

import (
	"fmt"
	"io"

	"github.com/TsvetanMilanov/go-simple-di/di"
	"github.com/TsvetanMilanov/go-srv/srv/log"
	"github.com/gin-gonic/gin"
)

type appBuilder struct {
	router            *gin.Engine
	appDI             *di.Container
	reqDI             *di.Container
	reqDIConfigurator RequestDIConfiguratorFunc
	appLogger         log.Logger
	err               error
}

func (ab *appBuilder) Initialize(loggerOut io.Writer) AppDependenciesRegisterer {
	ab.appDI = di.NewContainer()
	ab.reqDI = di.NewContainer()
	ab.appLogger = log.NewLogger(loggerOut)

	return ab
}

func (ab *appBuilder) RegisterAppDependencies(registerer DIContainerUser) AppDependenciesResolver {
	err := ab.appDI.Register(&di.Dependency{Name: AppLoggerName, Value: ab.appLogger})
	if err != nil {
		return ab.setErr("RegisterAppDependencies", "unable to register the app logger", err)
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

	return ab
}

func (ab *appBuilder) ConfigureRouter(configurator RouterConfiguratorFunc) AppBuilder {
	if ab.err != nil {
		return ab
	}

	err := configurator(ab.router, ab.appDI, ab.reqDI)
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
