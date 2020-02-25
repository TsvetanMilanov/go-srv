package srv

import (
	"context"
	"errors"

	"github.com/TsvetanMilanov/go-simple-di/di"
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
