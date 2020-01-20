package app

import (
	"log"
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// ConfigureRoutes will use Echo to start listening on the appropriate paths
func (a *Application) ConfigureRoutes(listenURL string) {

	// Middleware
	a.Echo.Use(middleware.Logger())
	a.Echo.Use(middleware.Recover())

	// CORS restricted
	// Allows requests from any `https://labstack.com` or `https://labstack.net` origin
	// wth GET, PUT, POST or DELETE method.
	a.Echo.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
	}))

	// Routes
	a.Echo.GET("/", a.defaultGet)

	a.addRegistrationRoutes()
	a.addInfoRoutes()
	a.addCommandRoutes()

	log.Println("Configure routes listening on " + listenURL)

	log.Println("***************************************")
	log.Println("~Rejoice~ GoFire Lives Again! ~Rejoice~")
	log.Println("***************************************")

	// Start server
	a.Echo.Logger.Fatal(a.Echo.Start(listenURL))
}

func (a *Application) defaultGet(c echo.Context) error {
	log.Println("Someone is touching me")

	return c.String(http.StatusOK, "Help Me! I'm trapped in the Server! You're the only one receiving this message.")
}
