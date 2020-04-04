package cluster

import (
	"firecontroller/microcontroller"
	mc "firecontroller/microcontroller"
	"firecontroller/utilities"
	"log"

	"github.com/spf13/viper"
)

//Cluster - This object defines an array of microcontrollers
type Cluster struct {
	Name                  string
	SlaveMicrocontrollers []mc.Microcontroller
	Master                *mc.Microcontroller
	Me                    *mc.Microcontroller
}

//Config -
type Config struct {
	Name                  string `yaml:"Name"`
	SlaveMicrocontrollers []mc.Config
	Master                mc.Config
}

//GetConfig -
func (c Cluster) GetConfig() (config Config) {
	config.Name = c.Name
	config.Master = c.Master.GetConfig()
	config.SlaveMicrocontrollers = make([]mc.Config, len(c.SlaveMicrocontrollers))
	for i, micro := range c.SlaveMicrocontrollers {
		config.SlaveMicrocontrollers[i] = micro.GetConfig()
	}
	return
}

//Load -
func (c *Cluster) Load(config Config) {
	c.Name = config.Name
	c.Master.Load(config.Master)
	c.SlaveMicrocontrollers = make([]mc.Microcontroller, len(config.SlaveMicrocontrollers))
	for i, micro := range config.SlaveMicrocontrollers {
		c.SlaveMicrocontrollers[i].Load(micro)
	}
}

func (c Cluster) String() string {
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
func (c Cluster) GetMicrocontrollers() map[int]microcontroller.Microcontroller {
	micros := make(map[int]microcontroller.Microcontroller)
	for i := 0; i < len(c.SlaveMicrocontrollers); i++ {
		micros[c.SlaveMicrocontrollers[i].ID] = c.SlaveMicrocontrollers[i]
	}
	return micros
}
