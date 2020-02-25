package srv

import (
	"net"
	"os"
	"testing"

	"github.com/TsvetanMilanov/go-graceful-server-shutdown/gss"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestStart(t *testing.T) {
	app, _, _, err := createApp(nil)

	assert.NoError(t, err)

	srvChan := make(chan error)
	srvStopChan := make(chan bool)

	go func() {
		err := app.Start(&gss.Settings{ShutdownChannel: srvStopChan, Addr: getFreeAddr()}, nil)
		srvChan <- err
	}()

	srvStopChan <- true
	err = <-srvChan

	assert.NoError(t, err)
}

func TestStartStopMetricsServer(t *testing.T) {
	app, _, _, err := createApp(nil)

	assert.NoError(t, err)

	srvChan := make(chan error)
	metricsSrvStopChan := make(chan bool)

	go func() {
		err := app.Start(&gss.Settings{Addr: getFreeAddr()}, &gss.Settings{ShutdownChannel: metricsSrvStopChan})
		srvChan <- err
	}()

	metricsSrvStopChan <- true
	err = <-srvChan

	assert.NoError(t, err)
}

func TestStartNilSrvSettings(t *testing.T) {
	app, _, _, err := createApp(nil)

	assert.NoError(t, err)

	err = app.Start(nil, nil)

	assert.EqualError(t, err, "srv: Start: srvSettings should not be nil")
}

func TestStartNoSrvAddr(t *testing.T) {
	app, _, _, err := createApp(nil)

	assert.NoError(t, err)

	err = app.Start(&gss.Settings{}, nil)

	assert.EqualError(t, err, "srv: Start: srvSettings.Addr should be a valid address")
}

func TestStartInvalidMetricsSrvAddr(t *testing.T) {
	app, _, _, err := createApp(nil)

	assert.NoError(t, err)

	err = app.Start(&gss.Settings{Addr: getFreeAddr()}, &gss.Settings{Addr: "invalid"})

	assert.EqualError(t, err, "srv: Start: unable to set the metrics server addr address invalid: missing port in address")
}

func TestStartWithoutMetricsServer(t *testing.T) {
	app, err := NewAppBuilder().
		Initialize(os.Stdout).
		RegisterAppDependencies(registerAppDependencies).
		ResolveAppDependencies().
		RegisterRouter(gin.New()).
		ConfigureRouter(configureRouter).
		BuildApp()

	assert.NoError(t, err)

	srvChan := make(chan error)
	srvStopChan := make(chan bool)

	go func() {
		err := app.Start(&gss.Settings{ShutdownChannel: srvStopChan, Addr: getFreeAddr()}, nil)
		srvChan <- err
	}()

	srvStopChan <- true
	err = <-srvChan

	assert.NoError(t, err)
}

func TestStartFallbackMetricsAddr(t *testing.T) {
	app, _, _, err := createApp(nil)

	assert.NoError(t, err)

	srvChan := make(chan error)
	srvStopChan := make(chan bool)

	metricsSrvAddr := getFreeAddr()
	addr, _ := net.ResolveTCPAddr("tcp", metricsSrvAddr)
	l, _ := net.ListenTCP("tcp", addr)
	defer l.Close()

	go func() {
		err := app.Start(&gss.Settings{ShutdownChannel: srvStopChan, Addr: getFreeAddr()}, &gss.Settings{Addr: metricsSrvAddr})
		srvChan <- err
	}()

	srvStopChan <- true
	err = <-srvChan

	assert.NoError(t, err)
}
