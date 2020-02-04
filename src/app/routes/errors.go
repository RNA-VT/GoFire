package routes

import (
	"net/http"

	"github.com/labstack/echo"
)

func (a APIService) addErrorRoutes(e *echo.Echo) {
	api := e.Group("/v1")
	api.GET("/errors", a.getErrors)
}

func (a APIService) getErrors(c echo.Context) error {
	//TODO: De-sass this endpoint
	return c.JSON(http.StatusMethodNotAllowed, "I've never made a mistake")
}
