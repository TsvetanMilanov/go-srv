package srv

import (
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

type testController struct {
	ReqDI *di.Container `di:"name=reqDI"`
}

func (tc *testController) Get(c *gin.Context) {
	dataSvc := new(testDataService)
	err := tc.ReqDI.Resolve(dataSvc)
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, dataSvc.GetData())
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

func configureRouter(router *gin.Engine, appDI, reqDI *di.Container) error {
	tc := new(testController)
	err := appDI.Resolve(tc)
	if err != nil {
		return err
	}

	router.GET("/data", tc.Get)

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
