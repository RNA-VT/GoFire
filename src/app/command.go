package app

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo"
)

func (a *Application) addCommandRoutes() {
	a.Echo.GET("/cmd", a.getCommands)
	a.Echo.GET("/component/:id/fire/:duration", a.fireSolenoid)
}

func (a *Application) getCommands(c echo.Context) error {
	//TODO: De-sass this endpoint
	return c.JSON(http.StatusMethodNotAllowed, "You cannot control me")
}

func (a *Application) fireSolenoid(c echo.Context) error {
	components := a.Cluster.GetComponents()
	componentID := c.Param("id")
	duration, err := strconv.Atoi(c.Param("duration"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Duration invalid.")
	}

	component, ok := components[componentID]
	if !ok {
		return c.JSON(http.StatusBadRequest, "Component Not Found.")
	}
	component.Open(duration)
	return c.JSON(http.StatusOK, a.Cluster.SlaveDevices)
}
func (a *Application) closeSolenoid(c echo.Context) error {
	return nil
}
func (a *Application) disableSolenoid(c echo.Context) error {
	return nil
}
