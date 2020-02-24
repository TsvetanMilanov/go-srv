package srv

import (
	"fmt"
	"net/http"

	"github.com/TsvetanMilanov/go-srv/srv/log"
	"github.com/gin-gonic/gin"
)

func recoverMiddlewareFactory(appLogger log.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			err := recover()
			if err != nil {
				appLogger.Error(fmt.Sprintf("srv: recoverMiddleware: recovered %s", err))

				c.JSON(http.StatusInternalServerError, map[string]string{"message": "Internal server error"})
			}
		}()

		c.Next()
	}
}

func contextPropertiesMiddlewareFactory(ab *appBuilder) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(AppLoggerName, ab.appLogger)
		c.Set(AppDIName, ab.appDI)
		c.Set(ReqDIName, ab.reqDI)

		c.Next()
	}
}

func requestDIConfiguratorMiddlewareFactory(ab *appBuilder) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := ab.reqDIConfigurator(c.Request, ab.reqDI)
		if err != nil {
			panic(fmt.Errorf("srv: requestDIConfiguratorMiddleware: unable to configure the request di %s", err))
		}

		c.Next()
	}
}
