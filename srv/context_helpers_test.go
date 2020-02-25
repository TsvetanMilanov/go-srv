package srv

import (
	"testing"

	"github.com/TsvetanMilanov/go-simple-di/di"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetReqDI(t *testing.T) {
	ctx := new(gin.Context)

	expectedContainer := new(di.Container)
	ctx.Set("reqDi", expectedContainer)

	actual, err := GetReqDI(ctx)

	assert.NoError(t, err)
	assert.Same(t, expectedContainer, actual)
}

func TestGetReqDINotRegistered(t *testing.T) {
	ctx := new(gin.Context)

	_, err := GetReqDI(ctx)

	assert.EqualError(t, err, "srv: GetReqDI: unable to get the request di from the context")
}

func TestGetReqDINotPointerToDIContainer(t *testing.T) {
	ctx := new(gin.Context)

	ctx.Set("reqDi", di.Container{})

	_, err := GetReqDI(ctx)

	assert.EqualError(t, err, "srv: GetReqDI: the registered req di doesn't have the correct type")
}
