package cluster

import (
	"bytes"
	"encoding/json"
	mc "firecontroller/microcontroller"
	"log"
	"net/http"
	"net/url"

	"github.com/spf13/viper"
)

//******************************************************************************************************
//*******Slave Only Methods*****************************************************************************
//******************************************************************************************************

//ALifeOfServitude is all that awaits this microcontroller
func (c *Cluster) ALifeOfServitude() {
	me, err := mc.NewMicrocontroller(viper.GetString("GOFIRE_HOST"), viper.GetString("GOFIRE_PORT"))
	if err != nil {
		log.Println("Failed to Create New Microcontroller:", err.Error())
	}
	me.ID = c.generateUniqueID()
	Me = &me
	masterHostname := viper.GetString("GOFIRE_MASTER_HOST") + ":" + viper.GetString("GOFIRE_MASTER_PORT")
	//Try and Connect to the Master
	err = test(masterHostname)
	if err != nil {
		log.Println("Failed to Reach Master Microcontroller: PANIC")
		//TODO: Add Retry or failover maybe? panic for now
		panic(err)
	}
	err = c.JoinNetwork(masterHostname)
	if err != nil {
		log.Println("Failed to Join Network: PANIC")
		panic(err)
	}
}

// JoinNetwork checks if the master exists and joins the network
func (c *Cluster) JoinNetwork(URL string) error {
	parsedURL, err := url.Parse("http://" + URL + "/join_network")
	log.Println("Trying to Join: " + parsedURL.String())
	msg := JoinNetworkMessage{
		ImNewHere: *Me,
	}
	body, err := json.Marshal(msg)
	if err != nil {
		log.Println("Failed to create json message body")
		return err
	}
	resp, err := http.Post(parsedURL.String(), "application/json", bytes.NewBuffer(body))

	if err != nil {
		log.Println("[test] Couldn't connect to master.", Me.ID)
		log.Println(err)
		return err
	}
	log.Println("Connected to master. Sending message to peers.")

	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	var t PeerUpdateMessage
	err = decoder.Decode(&t)
	if err != nil {
		log.Println("Failed to decode response from Master Microcontroller")
		log.Println(err)
		return err
	}
	//Update self with data from the master
	c.LoadCluster(t.Cluster)

	return nil
}

//LoadCluster sets all Cluster values except for Me
func (c *Cluster) LoadCluster(cluster Cluster) {
	log.Println("Loading Updated Cluster Data...")
	c.Name = cluster.Name
	c.Master = cluster.Master
	c.SlaveMicrocontrolers = cluster.SlaveMicrocontrolers
	PrintClusterInfo(*c)
}
