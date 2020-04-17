package cluster

import (
	"errors"
	mc "firecontroller/microcontroller"
	"log"
	"net/http"
	"time"

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
	me.Master = true
	//The master also serves
	c.Microcontrollers = append(c.Microcontrollers, me)
	c.Me = c.Master()

	//The Master pulls out its stethoscope ...
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		for {
			select {
			case t := <-ticker.C:
				log.Println("Begin Heartbeat Check", t)
				for _, m := range c.Microcontrollers {
					if !m.Master {
						log.Println("Checking Peer:", m.Name, m.ToFullAddress())
						url := "http://" + m.ToFullAddress() + "/v1/heartbeat"
						resp, err := http.Get(url)
						if err != nil || resp.StatusCode != 200 {
							log.Println(m.Name + " @" + m.ToFullAddress() + " is NOT ok")
							log.Println("Deregistering Microcontroller...")
							c.RemoveMicrocontroller(m)
						} else {
							log.Println(m.Name + " @" + m.ToFullAddress() + " is ok")
						}
					}
				}
			}
		}
	}()
}

//AddMicrocontroller attempts to add a microcontroller to the cluster and returns the response data. This should only be run by the master.
func (c *Cluster) AddMicrocontroller(newMC mc.Config) (response PeerUpdateMessage, err error) {
	var newGuy mc.Microcontroller
	newGuy.Load(newMC)
	newGuy.ID = c.generateUniqueID()
	if viper.GetString("ENV") == "production" {
		for _, micro := range c.Microcontrollers {
			if micro.Host == newGuy.Host {
				//This guy ain't so new!
				return PeerUpdateMessage{}, errors.New("Requesting instance is running on a microcontroller already registered to this cluster")
			}
		}
	}

	c.Microcontrollers = append(c.Microcontrollers, newGuy)

	PrintClusterInfo(*c)
	response = PeerUpdateMessage{
		Cluster: c.GetConfig(),
		Header:  c.GetHeader(),
	}

	exclusions := []mc.Microcontroller{newGuy, *c.Me}
	err = c.UpdatePeers("", response, exclusions)
	if err != nil {
		log.Println("Unexpected Error during attempt to contact all peers: ", err)
		return PeerUpdateMessage{}, err
	}

	return
}

//RemoveMicrocontroller -
func (c *Cluster) RemoveMicrocontroller(ImDoneHere mc.Microcontroller) (response PeerUpdateMessage, err error) {
	for index, mc := range c.Microcontrollers {
		if mc.ID == ImDoneHere.ID {
			s := c.Microcontrollers
			count := len(c.Microcontrollers)
			s[count-1], s[index] = s[index], s[count-1]
			c.Microcontrollers = s[:len(s)-1]
			return
		}
	}
	response = PeerUpdateMessage{
		Cluster: c.GetConfig(),
		Header:  c.GetHeader(),
	}

	exclusions := []mc.Microcontroller{*c.Me}
	err = c.UpdatePeers("", response, exclusions)
	if err != nil {
		log.Println("Unexpected Error during attempt to contact all peers: ", err)
		return PeerUpdateMessage{}, err
	}
	return
}
