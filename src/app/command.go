package app

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo"
)

func (a *Application) addCommandRoutes() {
	a.Echo.GET("/cmd", a.getCommands)
	a.Echo.GET("/component/:id/fire/:duration", a.fireSolenoid)
	a.Echo.GET("/component/:id/open", a.openSolenoid)
	a.Echo.GET("/component/:id/close", a.closeSolenoid)
	a.Echo.GET("/component/:id/enable", a.enableSolenoid)
	a.Echo.GET("/component/:id/disable", a.disableSolenoid)
}

func (a *Application) getCommands(c echo.Context) error {
	//TODO: De-sass this endpoint
	return c.JSON(http.StatusMethodNotAllowed, "You cannot control me")
}

func (a *Application) openSolenoid(c echo.Context) error {
	component, err := a.Cluster.GetComponent(c.Param("id"))
	if err != nil {
		return err
	}
	component.Open()
	return c.JSON(http.StatusOK, component.State())
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
	component.OpenFor(duration)
	return c.JSON(http.StatusOK, a.Cluster.SlaveMicrocontrolers)
}
func (a *Application) closeSolenoid(c echo.Context) error {
	component, err := a.Cluster.GetComponent(c.Param("id"))
	if err != nil {
		return err
	}
	component.Close(0)
	return c.JSON(http.StatusOK, component.State())
}
func (a *Application) disableSolenoid(c echo.Context) error {
	component, err := a.Cluster.GetComponent(c.Param("id"))
	if err != nil {
		return err
	}
	component.Disable()
	return c.JSON(http.StatusOK, component.State())
}

func (a *Application) enableSolenoid(c echo.Context) error {
	component, err := a.Cluster.GetComponent(c.Param("id"))
	if err != nil {
		return err
	}
	//I Guess Always
	component.Enable(true)
	return c.JSON(http.StatusOK, component.State())
}
