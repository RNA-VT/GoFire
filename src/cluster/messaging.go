package cluster

import (
	mc "firecontroller/microcontroller"
	"log"
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
	BaseMessage
}

//JoinNetworkMessage is the registration request
type JoinNetworkMessage struct {
	ImNewHere mc.Microcontroller
	BaseMessage
}

//PeerUpdateMessage contains a source and cluster info
type PeerUpdateMessage struct {
	Cluster Cluster
	BaseMessage
}

//BaseMessage -
type BaseMessage struct {
	Source  mc.Microcontroller
	Created time.Time
}

//EverybodyHasToKnow - Meant for Errors that should stop the entire cluster
func (c *Cluster) EverybodyHasToKnow(panicAfterWarning bool, panicCluster bool, MicrocontrollerToRemove mc.Microcontroller, notGoodThings ...string) {
	var message PeerErrorMessage
	message.Source = c.Me
	message.Messages = notGoodThings
	message.Panic = panicCluster
	message.DeregisterMe = MicrocontrollerToRemove
	c.UpdatePeers("errors", message, []mc.Microcontroller{c.Me})
	if panicAfterWarning {
		panic(notGoodThings)
	}
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
		log.Println("Deregistering Cluster: ", msg.DeregisterMe.String())
		for index, mc := range c.SlaveMicrocontrolers {
			if mc.ID == msg.DeregisterMe.ID {
				c.SlaveMicrocontrolers = RemoveMicrocontroller(c.SlaveMicrocontrolers, index)
			}
		}
	}
}
