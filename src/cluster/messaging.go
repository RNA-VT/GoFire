package cluster

import (
	"firecontroller/microcontroller"
	mc "firecontroller/microcontroller"
)

//JoinNetworkMessage is the registration request
type JoinNetworkMessage struct {
	ImNewHere mc.Microcontroller
}

//PeerUpdateMessage contains a source and cluster info
type PeerUpdateMessage struct {
	Source  mc.Microcontroller
	Cluster Cluster
}

//PeerErrorMessage -
type PeerErrorMessage struct {
	Panic bool
	PeerInfoMessage
}

//PeerInfoMessage -
type PeerInfoMessage struct {
	Source   mc.Microcontroller
	Messages []string
}

//ClusterError - Log the errors, warn the others and then panic.
func (c *Cluster) ClusterError(panicAfterWarning bool, panicCluster bool, notGoodThings ...string) {
	//Errors that render this microcontroller unusable, but do not effect the rest of the cluster
	if panicCluster {
		c.EverybodyPanic(notGoodThings...)
	} else {
		c.WarnTheOthers(notGoodThings...)
	}
	if panicAfterWarning {
		panic(notGoodThings)
	}
}

//WarnTheOthers - POST Error(s) to cluster.
func (c *Cluster) WarnTheOthers(msgs ...string) {
	//This path should be used for errors that make this instance of GoFire unavailable
	message := PeerErrorMessage{
		Panic: false,
	}
	message.Source = c.Me
	message.Messages = msgs
	c.tellTheOthers("/errors/warn", message)
}

//EverybodyPanic - Meant for Errors that should stop the entire cluster
func (c *Cluster) EverybodyPanic(notGoodThings ...string) {
	message := PeerErrorMessage{
		Panic: true,
	}
	message.Source = c.Me
	message.Messages = notGoodThings
	c.tellTheOthers("errors/panic", message)
}

func (c *Cluster) tellTheOthers(path string, msg interface{}) {
	c.UpdatePeers(
		path,
		msg,
		[]microcontroller.Microcontroller{c.Me})
}
