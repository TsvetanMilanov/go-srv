package srv

import (
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/TsvetanMilanov/go-simple-di/di"
	"github.com/gin-gonic/gin"
)

const (
	customMetricsAddr = ":33333"
)

var (
	dbData = map[string]interface{}{"key": "value"}
)

type testDBClient struct{}

func (tcl *testDBClient) GetData() map[string]interface{} {
	return dbData
}

type testDataService struct {
	DBClient *testDBClient `di:""`
}

func (ts *testDataService) GetData() map[string]interface{} {
	data := ts.DBClient.GetData()

	return data
}

type testController struct{}

func (tc *testController) Get(c *gin.Context) {
	dataSvc := new(testDataService)
	reqDI, _ := c.Get("reqDi")

	err := reqDI.(*di.Container).Resolve(dataSvc)
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, dataSvc.GetData())
}

func (tc *testController) Panic(c *gin.Context) {
	panic(errors.New("controller panic"))
}

func registerAppDependencies(container *di.Container) error {
	err := container.Register(
		&di.Dependency{Value: new(testController)},
	)

	return err
}

func registerReqDependencies(req *http.Request, container *di.Container) error {
	err := container.Register(
		&di.Dependency{Value: new(testDataService)},
		&di.Dependency{Value: new(testDBClient)},
	)

	return err
}

func configureRouter(router *gin.Engine, appDI *di.Container) error {
	tc := new(testController)
	err := appDI.Resolve(tc)
	if err != nil {
		return err
	}

	router.GET("/data", tc.Get)
	router.GET("/panic", tc.Panic)

	return nil
}

func performRequest(r http.Handler, method, path string) *httptest.ResponseRecorder {
	req, err := http.NewRequest(method, path, nil)
	if err != nil {
		panic(err)
	}

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	return w
}
