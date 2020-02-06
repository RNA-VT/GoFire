package cluster

import (
	"bytes"
	"encoding/json"
	mc "firecontroller/microcontroller"
	"firecontroller/utilities"
	"log"
	"net/http"
	"time"
)

//PeerErrorMessage -
type PeerErrorMessage struct {
	Panic        bool
	DeregisterMe mc.Microcontroller
	PeerInfoMessage
}

//PeerInfoMessage -
type PeerInfoMessage struct {
	Messages []string
	Header   GoFireHeader
}

//JoinNetworkMessage is the registration request
type JoinNetworkMessage struct {
	ImNewHere mc.Microcontroller
	Header    GoFireHeader
}

//PeerUpdateMessage contains a source and cluster info
type PeerUpdateMessage struct {
	Cluster Cluster
	Header  GoFireHeader
}

//GoFireHeader -
type GoFireHeader struct {
	Source  mc.Microcontroller
	Created time.Time
}

//GetHeader -
func GetHeader() GoFireHeader {
	return GoFireHeader{
		Source:  *Me,
		Created: time.Now(),
	}
}

//EverybodyHasToKnow - Meant for Errors that should stop the entire cluster
func (c Cluster) EverybodyHasToKnow(panicAfterWarning bool, panicCluster bool, MicrocontrollerToRemove mc.Microcontroller, notGoodThings ...string) {
	var message PeerErrorMessage
	message.Messages = notGoodThings
	message.Panic = panicCluster
	message.DeregisterMe = MicrocontrollerToRemove
	message.Header = GetHeader()
	c.UpdatePeers("errors", message, []mc.Microcontroller{*Me})
	if panicAfterWarning {
		panic(notGoodThings)
	}
}

// UpdatePeers will take a byte slice and POST it to each microcontroller
func (c Cluster) UpdatePeers(urlPath string, message interface{}, exclude []mc.Microcontroller) error {
	for i := 0; i < len(c.SlaveMicrocontrolers); i++ {
		if !isExcluded(c.SlaveMicrocontrolers[i], exclude) {
			body, err := utilities.JSON(message)
			if err != nil {
				log.Println("Failed to convert cluster to json: ", c)
				return err
			}
			currURL := "http://" + c.SlaveMicrocontrolers[i].ToFullAddress() + urlPath

			resp, err := http.Post(currURL, "application/json", bytes.NewBuffer(body))
			if err != nil {
				log.Println("WARNING: Failed to POST to Peer: ", c.SlaveMicrocontrolers[i].String(), currURL)
				log.Println(err)
			} else {
				defer resp.Body.Close()
				var result string
				decoder := json.NewDecoder(resp.Body)
				decoder.Decode(&result)
				log.Println("Result:", result)
			}
		}
	}
	return nil
}

//ReceiveError -
func (c *Cluster) ReceiveError(msg PeerErrorMessage) {
	//log msgs to console
	for msg := range msg.Messages {
		log.Println(msg)
	}
	if msg.Panic {
		panic(map[string]interface{}{
			"Cluster": c,
			"Message": msg,
		})
	}
	// TODO do better with this check
	if msg.DeregisterMe.Host != "" {
		//Deregister Microcontroller
		log.Println("Deregistering Microcontroller From Cluster: ", msg.DeregisterMe.String())
		c.RemoveMicrocontroller(msg.DeregisterMe)
	}
}
