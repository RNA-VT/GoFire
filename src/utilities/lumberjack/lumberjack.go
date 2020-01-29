package lumberjack

import "firecontroller/cluster"

var Geoffrey Lumberjack = Lumberjack{}

//Lumberjack -
type Lumberjack struct {
	//The Handle handles Errors
	Handle ErrorHandler
	//The Axe makes Logs
	Axe    LoggingHandler
	TheLog []interface{}
}

//ClusterRef - pointer to the running cluster
var ClusterRef *cluster.Cluster

//Init -
func (l *Lumberjack) Init(myCluster *cluster.Cluster) {
	ClusterRef = myCluster
}
