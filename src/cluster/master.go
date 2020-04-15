package cluster

import (
	"errors"
	mc "firecontroller/microcontroller"
	"log"

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
	//The master also serves
	c.SlaveMicrocontrollers = append(c.SlaveMicrocontrollers, me)
	c.Master = c.SlaveMicrocontrollers[len(c.SlaveMicrocontrollers)-1]
	c.Me = &c.Master

	//The Master waits ...
}

//AddMicrocontroller attempts to add a microcontroller to the cluster and returns the response data. This should only be run by the master.
func (c *Cluster) AddMicrocontroller(newMC mc.Config) (response PeerUpdateMessage, err error) {
	var newGuy mc.Microcontroller
	newGuy.Load(newMC)
	newGuy.ID = c.generateUniqueID()
	if viper.GetString("ENV") == "production" {
		for _, micro := range c.SlaveMicrocontrollers {
			if micro.Host == newGuy.Host {
				//This guy ain't so new!
				return PeerUpdateMessage{}, errors.New("Requesting instance is running on a microcontroller already registered to this cluster")
			}
		}
	}

	c.SlaveMicrocontrollers = append(c.SlaveMicrocontrollers, newGuy)
	//FUCK: OK It's definitely wrong here, the newGuy is all 0s

	log.Println("DEFINITELY BROKEN HERE", newGuy, c.GetConfig())
	PrintClusterInfo(*c)
	response = PeerUpdateMessage{
		Cluster: c.GetConfig(),
		Header:  c.GetHeader(),
	}

	exclusions := []mc.Microcontroller{newGuy, *c.Me}
	err = c.UpdatePeers("/", response, exclusions)
	if err != nil {
		log.Println("Unexpected Error during attempt to contact all peers: ", err)
		return PeerUpdateMessage{}, err
	}

	return response, nil
}

//RemoveMicrocontroller -
func (c *Cluster) RemoveMicrocontroller(ImDoneHere mc.Microcontroller) {
	for index, mc := range c.SlaveMicrocontrollers {
		if mc.ID == ImDoneHere.ID {
			s := c.SlaveMicrocontrollers
			count := len(c.SlaveMicrocontrollers)
			s[count-1], s[index] = s[index], s[count-1]
			c.SlaveMicrocontrollers = s[:len(s)-1]
			return
		}
	}
}
