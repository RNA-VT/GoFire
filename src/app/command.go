package app

import (
	"net/http"

	"github.com/labstack/echo"
)

func (a *Application) addCommandRoutes() {
	a.Echo.GET("/cmd", a.getCommands)
}

func (a *Application) getCommands(c echo.Context) error {
	//TODO: De-sass this endpoint
	return c.JSON(http.StatusMethodNotAllowed, "You cannot control me")
}
