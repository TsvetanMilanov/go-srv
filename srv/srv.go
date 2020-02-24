package srv

import (
	"fmt"
	"net"
	"net/http"

	"github.com/TsvetanMilanov/go-graceful-server-shutdown/gss"
	"github.com/TsvetanMilanov/go-srv/srv/log"
	"github.com/prometheus/client_golang/prometheus"
)

type app struct {
	appLogger       log.Logger
	router          http.Handler
	metricsRouter   http.Handler
	metricsRegistry *prometheus.Registry
	metricsAddr     string
}

func (a *app) Start(addr string) error {
	if a.metricsRouter != nil {
		err := a.setMetricsServerAddr()
		if err != nil {
			return err
		}

		srvChan := make(chan error)
		srvStopChan := make(chan bool)
		metricsSrvChan := make(chan error)
		metricsSrvStopChan := make(chan bool)

		go func() {
			settings := &gss.Settings{Addr: a.metricsAddr, ShutdownChannel: metricsSrvStopChan}
			err := gss.StartServerWithSettings(a.metricsRouter, settings)

			metricsSrvChan <- err
		}()

		go func() {
			settings := &gss.Settings{Addr: addr, ShutdownChannel: srvStopChan}
			err := gss.StartServerWithSettings(a.router, settings)

			srvChan <- err
		}()

		select {
		case err := <-metricsSrvChan:
			srvStopChan <- true
			return err
		case err := <-srvChan:
			metricsSrvStopChan <- true
			return err
		}
	} else {
		err := gss.StartServer(http.DefaultServeMux)

		return err
	}
}

func (a *app) GetRouter() http.Handler {
	return a.router
}

func (a *app) GetMetricsRouter() http.Handler {
	return a.metricsRouter
}

func (a *app) GetLogger() log.Logger {
	return a.appLogger
}

func (a *app) setMetricsServerAddr() error {
	addr, err := net.ResolveTCPAddr("tcp", a.metricsAddr)
	if err != nil {
		return err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		// Fallback to free port if running two instances of this module locally.
		addr, err = net.ResolveTCPAddr("tcp", ":0")
		if err != nil {
			return err
		}

		l, err = net.ListenTCP("tcp", addr)
		if err != nil {
			return err
		}

		a.metricsAddr = fmt.Sprintf(":%d", l.Addr().(*net.TCPAddr).Port)
	}

	l.Close()

	return nil
}

func newApp(ab *appBuilder) App {
	a := new(app)
	a.appLogger = ab.appLogger
	a.router = ab.router
	a.metricsRouter = ab.metricsRouter
	a.metricsRegistry = ab.metricsRegistry
	if len(ab.metricsAddr) == 0 {
		a.metricsAddr = ":80"
	} else {
		a.metricsAddr = ab.metricsAddr
	}

	return a
}
