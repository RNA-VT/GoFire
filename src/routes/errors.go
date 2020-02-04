package routes

import "github.com/labstack/echo"

func (a APIService) addErrorRoutes(e *echo.Echo, version string) {
	api := e.Group("/v" + version + "/errors")
	api.GET("/", a.getErrors)
	api.POST("/warn", a.handleWarning)
	api.POST("/panic", a.handlePanic)
}

func (a APIService) getErrors(c echo.Context) error {
	return nil
}

func (a APIService) handleWarning(c echo.Context) error {
	//a.Cluster.ReceiveWarning()
	return nil
}

func (a APIService) handlePanic(c echo.Context) error {
	//a.Cluster.ReceivePanic()
	return nil
}
