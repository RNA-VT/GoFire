package cluster

import (
	"bytes"
	"encoding/json"
	mc "firecontroller/microcontroller"
	"firecontroller/utilities"
	"log"
	"net/http"

	"github.com/spf13/viper"
)

//******************************************************************************************************
//*******Master Only Methods****************************************************************************
//******************************************************************************************************

//KingMe makes this microcontroller the master
func (c *Cluster) KingMe() {
	me, err := mc.NewMicrocontroller(viper.GetString("GOFIRE_MASTER_HOST"), viper.GetString("GOFIRE_MASTER_PORT"))
	if err != nil {
		log.Println("Failed to Create New Microcontroller:", err.Error())
	}
	me.ID = c.generateUniqueID()
	Me = &me
	c.Master = me
	//The master also serves
	c.SlaveMicrocontrolers = append(c.SlaveMicrocontrolers, me)
	//The Master waits ...
}

// UpdatePeers will take a byte slice and POST it to each microcontroller
func (c *Cluster) UpdatePeers(urlPath string, message interface{}, exclude []mc.Microcontroller) error {
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

//AddMicrocontroller attempts to add a microcontroller to the cluster and returns the response data. This should only be run by the master.
func (c *Cluster) AddMicrocontroller(newMC mc.Microcontroller) (response PeerUpdateMessage, err error) {
	newMC.ID = c.generateUniqueID()
	c.SlaveMicrocontrolers = append(c.SlaveMicrocontrolers, newMC)
	PrintClusterInfo(*c)

	response = PeerUpdateMessage{
		Cluster: *c,
	}
	response.Source = *Me

	exclusions := []mc.Microcontroller{newMC, *Me}
	err = c.UpdatePeers("/", response, exclusions)
	if err != nil {
		log.Println("Unexpected Error during attempt to contact all peers: ", err)
		return PeerUpdateMessage{}, err
	}

	return response, nil
}

//RemoveMicrocontroller -
func (c *Cluster) RemoveMicrocontroller(ImDoneHere mc.Microcontroller) {
	for index, mc := range c.SlaveMicrocontrolers {
		if mc.ID == ImDoneHere.ID {
			s := c.SlaveMicrocontrolers
			count := len(c.SlaveMicrocontrolers)
			s[count-1], s[index] = s[index], s[count-1]
			c.SlaveMicrocontrolers = s[:len(s)-1]
			return
		}
	}
}
