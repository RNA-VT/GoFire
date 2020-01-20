package app

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/labstack/echo"
)

func (a *Application) addInfoRoutes() {
	a.Echo.GET("/cluster_info", a.getClusterInfo)
	a.Echo.GET("/microcontroller", a.getMicrocontrollers)
	a.Echo.GET("/microcontroller/:id", a.getMicrocontrollers)
	a.Echo.GET("/component", a.getComponents)
	a.Echo.GET("/component/:id", a.getComponent)
	a.Echo.GET("/config", a.getComponentConfig)
}

func (a *Application) getClusterInfo(c echo.Context) error {
	return c.JSON(http.StatusOK, a.Cluster)
}

func (a *Application) getMicrocontrollers(c echo.Context) error {
	return c.JSON(http.StatusOK, a.Cluster.SlaveDevices)
}

func (a *Application) getComponents(c echo.Context) error {
	return c.JSON(http.StatusOK, a.Cluster.GetComponents())
}

func (a *Application) getComponent(c echo.Context) error {
	id := c.Param("id")
	components := a.Cluster.GetComponents()
	component, ok := components[id]
	if !ok {
		return c.JSON(http.StatusOK, "Component Not Found")
	}
	return c.JSON(http.StatusOK, component)
}

func (a *Application) getComponentConfig(c echo.Context) error {
	yamlFile, err := ioutil.ReadFile("./app/config/solenoids.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
		return err
	}
	return c.JSON(http.StatusOK, string(yamlFile))
}
