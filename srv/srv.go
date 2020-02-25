package srv

import (
	"fmt"
	"net"
	"net/http"

	"github.com/TsvetanMilanov/go-graceful-server-shutdown/gss"
	"github.com/TsvetanMilanov/go-srv/srv/log"
)

type app struct {
	appLogger     log.Logger
	router        http.Handler
	metricsRouter http.Handler
}

func (a *app) Start(srvSettings, metricsSrvSettings *gss.Settings) error {
	a.setSrvSettings(srvSettings, make(chan bool))
	a.setSrvSettings(metricsSrvSettings, make(chan bool))

	if len(srvSettings.Addr) == 0 {
		return fmt.Errorf("srv: Start: srvSettings.Addr should be a valid address")
	}

	if a.metricsRouter != nil {
		err := a.setMetricsServerAddr(metricsSrvSettings)
		if err != nil {
			return err
		}

		srvChan := make(chan error)
		metricsSrvChan := make(chan error)

		go func() {
			err := gss.StartServerWithSettings(a.metricsRouter, metricsSrvSettings)

			metricsSrvChan <- err
		}()

		go func() {
			err := gss.StartServerWithSettings(a.router, srvSettings)

			srvChan <- err
		}()

		select {
		case err := <-metricsSrvChan:
			srvSettings.ShutdownChannel <- true
			return err
		case err := <-srvChan:
			metricsSrvSettings.ShutdownChannel <- true
			return err
		}
	} else {
		err := gss.StartServerWithSettings(a.router, srvSettings)

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

func (a *app) setMetricsServerAddr(metricsSrvSettings *gss.Settings) error {
	if len(metricsSrvSettings.Addr) == 0 {
		metricsSrvSettings.Addr = defaultMetricsServerAddr
	}

	addr, err := net.ResolveTCPAddr("tcp", metricsSrvSettings.Addr)
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

		metricsSrvSettings.Addr = fmt.Sprintf(":%d", l.Addr().(*net.TCPAddr).Port)
	}

	l.Close()

	return nil
}

func (a *app) setSrvSettings(settings *gss.Settings, shutdownChannel chan bool) {
	if settings == nil {
		settings = new(gss.Settings)
	}

	if settings.ShutdownChannel == nil {
		settings.ShutdownChannel = shutdownChannel
	}
}

func newApp(ab *appBuilder) App {
	a := new(app)
	a.appLogger = ab.appLogger
	a.router = ab.router
	a.metricsRouter = ab.metricsRouter

	return a
}
