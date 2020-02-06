package cluster

import (
	"errors"
	"firecontroller/component"
	"firecontroller/microcontroller"
	mc "firecontroller/microcontroller"
	"firecontroller/utilities"
	"log"
	"strconv"

	"github.com/spf13/viper"
)

//Cluster - This object defines an array of microcontrollers
type Cluster struct {
	Name                 string
	SlaveMicrocontrolers []mc.Microcontroller
	Master               mc.Microcontroller
}

//Me - a reference to this micros instance in the slave list
var Me *mc.Microcontroller

func (c *Cluster) String() string {
	cluster, err := utilities.StringJSON(c)
	if err != nil {
		return ""
	}
	return cluster
}

//Start registers this microcontroller, retrieves cluster config, loads local components and verifies peers
func (c *Cluster) Start() {
	//Set global ref to cluster
	gofireMaster := viper.GetBool("GOFIRE_MASTER")
	if gofireMaster {
		log.Println("Master Mode Enabled!")
		c.KingMe()
	} else {
		log.Println("Slave Mode Enabled.")
		c.ALifeOfServitude()
	}
}

//GetMicrocontrollers returns a map[microcontrollerID]microcontroller of all Microcontrollers in the cluster
func (c *Cluster) GetMicrocontrollers() map[int]microcontroller.Microcontroller {
	micros := make(map[int]microcontroller.Microcontroller)
	for i := 0; i < len(c.SlaveMicrocontrolers); i++ {
		micros[c.SlaveMicrocontrolers[i].ID] = c.SlaveMicrocontrolers[i]
	}
	return micros
}

//GetComponent - gets a component by its id
func (c *Cluster) GetComponent(id string) (sol component.Solenoid, err error) {
	components := c.GetComponents()
	sol, ok := components[id]
	if !ok {
		return sol, errors.New("Component Not Found")
	}
	return sol, nil
}

//GetComponents builds a map of all the components in the cluster by a cluster wide unique key
func (c *Cluster) GetComponents() map[string]component.Solenoid {
	components := make(map[string]component.Solenoid, c.countComponents())
	for i := 0; i < len(c.SlaveMicrocontrolers); i++ {
		for j := 0; j < len(c.SlaveMicrocontrolers[i].Solenoids); j++ {
			key := strconv.Itoa(c.SlaveMicrocontrolers[i].Solenoids[j].UID)
			components[key] = c.SlaveMicrocontrolers[i].Solenoids[j]
		}
	}
	return components
}
func (c *Cluster) countComponents() int {
	count := 0
	for i := 0; i < len(c.SlaveMicrocontrolers); i++ {
		count += len(c.SlaveMicrocontrolers[i].Solenoids)
	}

	return count
}
