package cluster

import (
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
	Me = &me
	c.Master = me
	//The master also serves
	c.SlaveMicrocontrolers = append(c.SlaveMicrocontrolers, me)
	//The Master waits ...
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
