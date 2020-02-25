package srv_test

import (
	"net/http"
	"os"

	"github.com/TsvetanMilanov/go-graceful-server-shutdown/gss"
	"github.com/TsvetanMilanov/go-simple-di/di"
	"github.com/TsvetanMilanov/go-srv/srv"
	"github.com/TsvetanMilanov/go-srv/srv/log"
	"github.com/gin-gonic/gin"
)

type dbClient struct{}

func (db *dbClient) GetData() map[string]interface{} {
	return map[string]interface{}{"key": "value"}
}

type dataService struct {
	DBClient *dbClient `di:""`

	// The reqLogger dependency will be included in the reqDI container.
	RequestLogger log.Logger `di:"name=reqLogger"`
	// A srv.TraceIDProvider dependency will be included in the reqDI container.
	TraceIDProvider srv.TraceIDProvider `di:""`
}

func (s *dataService) GetData() map[string]interface{} {
	s.RequestLogger.Info("this is log from the request logger")
	s.RequestLogger.Info(s.TraceIDProvider.GetTraceID())

	// Use the resolved DBClient.
	return s.DBClient.GetData()
}

type controller struct{}

func (controller *controller) Get(c *gin.Context) {
	dataSvc := new(dataService)
	reqDI, _ := srv.GetReqDI(c)

	err := reqDI.Resolve(dataSvc)
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, dataSvc.GetData())
}

func Example() {
	// Function which registers all application dependencies.
	// This function will be invoked only once.
	var registerAppDependencies = func(container *di.Container) error {
		return container.Register(
			&di.Dependency{Value: new(controller)},
		)
	}

	// This function will be invoked before each request.
	// You can register all dependencies which needs to be resolved for each request here.
	// The container instance will be new for each request.
	var reqDIConfigurator = func(req *http.Request, container *di.Container) error {
		return container.Register(
			&di.Dependency{Value: new(dataService)},
			&di.Dependency{Value: new(dbClient)},
		)
	}

	// Function which receives a router object which has some predefined configurations
	// and can be used to finish the router configuration with the resolved app dependencies.
	var configureRouter = func(router *gin.Engine, appDI *di.Container) error {
		// Set the env to production.
		if env, ok := os.LookupEnv("GO_ENV"); ok && env == "production" {
			gin.SetMode(gin.ReleaseMode)
		}

		tc := new(controller)
		// Use the appDI container to resolve the registered app dependencies (e.g. controllers)
		err := appDI.Resolve(tc)
		if err != nil {
			return err
		}

		// Register all routes.
		router.GET("/data", tc.Get)

		return nil
	}

	app, err := srv.NewAppBuilder().
		Initialize(os.Stdout).
		EnableMetricsServer(nil, nil). // You can pass custom prometheus gatherer or registerer here.
		RegisterAppDependencies(registerAppDependencies).
		ResolveAppDependencies().
		RegisterReqDIConfigurator(reqDIConfigurator).
		ConfigureApp(func(appDI *di.Container) error { return nil }).
		RegisterRouter(gin.New()).
		AddDefaultMiddlewares().
		ConfigureRouter(configureRouter).
		BuildApp()

	if err != nil {
		panic(err)
	}

	err = app.Start(&gss.Settings{Addr: ":80"}, nil)
}
