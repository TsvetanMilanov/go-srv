package srv

import (
	"net/http"

	"github.com/TsvetanMilanov/go-graceful-server-shutdown/gss"
	"github.com/TsvetanMilanov/go-srv/srv/log"
)

type app struct {
	appLogger log.Logger
	router    http.Handler
}

func (a *app) Start() error {
	err := gss.StartServer(http.DefaultServeMux)

	return err
}

func (a *app) GetRouter() http.Handler {
	return a.router
}

func (a *app) GetLogger() log.Logger {
	return a.appLogger
}

func newApp(ab *appBuilder) App {
	a := new(app)
	a.appLogger = ab.appLogger
	a.router = ab.router

	return a
}
