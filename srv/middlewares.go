package srv

import (
	"fmt"
	"net/http"
	"time"

	"github.com/TsvetanMilanov/go-simple-di/di"
	"github.com/TsvetanMilanov/go-srv/srv/log"
	"github.com/gin-gonic/gin"
)

func recoverMiddlewareFactory(appLogger log.Logger) gin.HandlerFunc {
	start := time.Now()
	return func(c *gin.Context) {
		defer func() {
			err := recover()
			if err != nil {
				logger := appLogger
				reqDI, err := GetReqDI(c)
				if err == nil {
					reqLogger := new(log.Logger)
					err = reqDI.ResolveByName(ReqLoggerName, reqLogger)
					if err == nil {
						logger = *reqLogger
					}
				}

				logger.Error(fmt.Sprintf("srv: recoverMiddleware: recovered %s", err))

				c.JSON(http.StatusInternalServerError, map[string]string{"message": "Internal server error"})

				logger.Info(createLogEntryFields(start, c))
			}
		}()

		c.Next()
	}
}

func contextPropertiesMiddlewareFactory(ab *appBuilder) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(AppLoggerName, ab.appLogger)
		c.Set(AppDIName, ab.appDI)
	}
}

func requestDIConfiguratorMiddlewareFactory(ab *appBuilder) gin.HandlerFunc {
	return func(c *gin.Context) {
		// We need new container for each request.
		reqDI := di.NewContainer()
		c.Set(ReqDIName, reqDI)

		reqLogger := ab.appLogger.CreateChild()
		traceIDProvider := newRequestTraceIDProvider(c.Request)
		reqLogger.AddFields(map[string]interface{}{TraceIDName: traceIDProvider.GetTraceID()})

		err := reqDI.Register(
			&di.Dependency{Name: ReqLoggerName, Value: reqLogger},
			&di.Dependency{Value: traceIDProvider},
			&di.Dependency{Value: reqDI},
		)

		if err != nil {
			panic(fmt.Errorf("srv: requestDIConfiguratorMiddleware: unable to register the default req di dependencies %s", err))
		}

		err = ab.reqDIConfigurator(c.Request, reqDI)
		if err != nil {
			panic(fmt.Errorf("srv: requestDIConfiguratorMiddleware: unable to configure the request di %s", err))
		}
	}
}

func responseLoggerMiddlewareFactory(ab *appBuilder) gin.HandlerFunc {
	return func(c *gin.Context) {
		reqDI, err := GetReqDI(c)
		if err != nil {
			panic(fmt.Errorf("srv: responseLoggerMiddleware: unable to get the req di %s", err))
		}

		reqLogger := new(log.Logger)
		err = reqDI.ResolveByName(ReqLoggerName, reqLogger)
		if err != nil {
			panic(fmt.Errorf("srv: responseLoggerMiddleware: unable to resolve the request logger %s", err))
		}

		start := time.Now()
		c.Next()

		(*reqLogger).AddFields(createLogEntryFields(start, c))

		(*reqLogger).Info("request complete")
	}
}

func createLogEntryFields(start time.Time, c *gin.Context) map[string]interface{} {
	durationSeconds := float64(time.Since(start)) / float64(time.Second)

	req := c.Request

	entry := map[string]interface{}{
		"res": map[string]interface{}{
			"statusCode": c.Writer.Status(),
			"latency":    durationSeconds,
			"headers":    c.Writer.Header(),
		},
		"req": map[string]interface{}{
			"method":        req.Method,
			"path":          req.URL.Path,
			"remoteAddress": req.RemoteAddr,
			"clientIp":      c.ClientIP(),
			"headers":       req.Header,
		},
	}

	return entry
}
