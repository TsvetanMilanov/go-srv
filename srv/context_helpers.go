package srv

import (
	"context"
	"errors"

	"github.com/TsvetanMilanov/go-simple-di/di"
	"github.com/TsvetanMilanov/go-srv/srv/log"
)

// GetReqDI tries to get the request di from the context.
func GetReqDI(c context.Context) (*di.Container, error) {
	reqDIInstance := c.Value(ReqDIName)

	if reqDIInstance == nil {
		return nil, errors.New("srv: GetReqDI: unable to get the request di from the context")
	}

	reqDI, ok := reqDIInstance.(*di.Container)
	if !ok {
		return nil, errors.New("srv: GetReqDI: the registered req di doesn't have the correct type")
	}

	return reqDI, nil
}

// GetRequestLoggerOrDefaultChild returns the request logger or the child of the default logger.
func GetRequestLoggerOrDefaultChild(c context.Context, defaultLogger log.Logger) log.Logger {
	var logger log.Logger
	reqDI, err := GetReqDI(c)
	if err == nil {
		reqLogger := new(log.Logger)
		err = reqDI.ResolveByName(ReqLoggerName, reqLogger)
		if err == nil {
			logger = *reqLogger
		}
	}

	if logger == nil {
		logger = defaultLogger.CreateChild()
	}

	return logger
}
