package app

import (
	"encoding/json"
	"firecontroller/cluster"
	"log"
	"net/http"

	"github.com/labstack/echo"
)

func (a *Application) addRegistrationRoutes() {
	a.Echo.POST("/", a.peerUpdate)
	a.Echo.POST("/join_network", a.joinNetwork)
}

func (a *Application) joinNetwork(c echo.Context) error {
	log.Println("[master] Device asked to join cluster")

	body := c.Request().Body
	decoder := json.NewDecoder(body)
	var msg cluster.JoinNetworkMessage
	err := decoder.Decode(&msg)
	if err != nil {
		log.Println("Error decoding Request Body", err)
	}

	response, err := a.Cluster.AddDevice(msg.ImNewHere)
	if err != nil {
		log.Println("Error Joining Cluster")
	}

	return c.JSON(http.StatusOK, response)
}

//PeerUpdate receives new cluster info from the most recently registered peer
func (a *Application) peerUpdate(c echo.Context) error {
	log.Println("Receiving Update from New Peer")
	body := c.Request().Body

	var clustahUpdate cluster.PeerUpdateMessage
	err := json.NewDecoder(body).Decode(&clustahUpdate)
	if err != nil {
		log.Println("Failed to decode Cluster info from new peer")
		//TODO: Add Cluster Info Request to repair Cluster info
		return err
	}
	//TODO: Verify my presence in SlaveList
	//TODO: Verify my Master state
	//TODO: Inform Master of Bad Config

	//Update my cluster
	a.Cluster.LoadCluster(clustahUpdate.Cluster)

	log.Println("Peer Update Completed")
	return c.JSON(http.StatusOK, "Peer Update Successfully Received by : "+a.Cluster.Me.String())
}
