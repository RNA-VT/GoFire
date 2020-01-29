package lumberjack

import (
	"firecontroller/microcontroller"
	"firecontroller/utilities"
	"log"
	"time"
)

//ErrorHandler -
type ErrorHandler []GoFireError

//GoFireError -
type GoFireError struct {
	Code           int `yaml:"code"`
	GoError        error
	GoFireLogEntry `yaml:"inline"`
}

//PeerErrorMessage -
type PeerErrorMessage struct {
	Source microcontroller.Microcontroller
	Errors []GoFireError
}

//NewError -
func NewError(code int, model string, method string, msg string, err error) *GoFireError {
	notGoodThing := GoFireError{
		Code:    code,
		GoError: err,
	}

	notGoodThing.CreateDate = time.Now()
	notGoodThing.Model = model
	notGoodThing.Method = method
	notGoodThing.Message = msg

	return &notGoodThing
}

//UhOh -
func (e *ErrorHandler) UhOh(code int, model string, method string, msg string, err error) {
	e.insertError(NewError(code, model, method, msg, err))
}

//UhOhs - several new GoFireError to the log and Print it to the terminal
func (e *ErrorHandler) UhOhs(notGoodThings ...GoFireError) {
	//Local errors
	for _, notGoodThing := range notGoodThings {
		e.insertError(&notGoodThing)
		json, err := utilities.StringJSON(notGoodThing)
		if err != nil {
			e.insertError(NewError(-1, "GoFireErrors", "UhOh", "unable to json", err))
		}
		log.Println(json)
	}
}

//ClusterError - Log the errors, warn the others and then panic.
func (e *ErrorHandler) ClusterError(panicAfterWarning bool, panicCluster bool, notGoodThings ...GoFireError) {
	//Errors that render this microcontroller unusable, but do not effect the rest of the cluster
	e.UhOhs(notGoodThings...)
	if panicCluster {
		e.EverybodyPanic(*e...)
	} else {
		e.WarnTheOthers(*e...)
	}
	if panicAfterWarning {
		panic(e)
	}
}

//WarnTheOthers - POST Error(s) to cluster.
func (e *ErrorHandler) WarnTheOthers(notGoodThings ...GoFireError) {
	//This path should be used for errors that make this instance of GoFire unavailable
	e.tellTheOthers("/errors/warn", notGoodThings...)
}

//EverybodyPanic - Meant for Errors that should stop the entire cluster
func (e *ErrorHandler) EverybodyPanic(notGoodThings ...GoFireError) {
	e.tellTheOthers("errors/panic", notGoodThings...)
}

func (e *ErrorHandler) tellTheOthers(path string, notGoodThings ...GoFireError) {
	c := *ClusterRef
	c.UpdatePeers(
		path,
		PeerErrorMessage{
			Source: c.Me,
			Errors: notGoodThings,
		},
		[]microcontroller.Microcontroller{c.Me})
}

//insertError -
func (e *ErrorHandler) insertError(notGoodThing *GoFireError) {
	l := len(*e)
	target := *e
	if cap(*e) == l {
		target = make([]GoFireError, l+1, l+10)
		copy(target, *e)
		target[l] = *notGoodThing
	} else {
		target = append(target, *notGoodThing)
	}
	e = &target
}
