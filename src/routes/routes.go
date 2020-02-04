package routes

import (
	"firecontroller/cluster"
	"log"
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

//APIService -
type APIService struct {
	Cluster *cluster.Cluster
}

//API - Container object for API worker methods
var API APIService

const apiVersion = "1"

// ConfigureRoutes will use Echo to start listening on the appropriate paths
func ConfigureRoutes(listenURL string, e *echo.Echo) {

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// CORS restricted
	// Allows requests from any `https://labstack.com` or `https://labstack.net` origin
	// wth GET, PUT, POST or DELETE method.
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
	}))
	// Routes
	e.GET("/", API.defaultGet)
	e.GET("/v"+apiVersion, API.defaultGet)

	API.addRegistrationRoutes(e, apiVersion)
	API.addInfoRoutes(e, apiVersion)
	API.addCommandRoutes(e, apiVersion)

	log.Println("Configure routes listening on " + listenURL)

	log.Println("***************************************")
	log.Println("~Rejoice~ GoFire Lives Again! ~Rejoice~")
	log.Println("***************************************")

	// Start server
	e.Logger.Fatal(e.Start(listenURL))
}

func (a APIService) defaultGet(c echo.Context) error {
	log.Println("Someone is touching me")
	return c.String(http.StatusOK, "Help Me! I'm trapped in the Server! You're the only one receiving this message.")
}
