package lumberjack

import (
	"firecontroller/utilities"
	"log"
	"time"
)

//LoggingHandler -
type LoggingHandler []GoFireLogEntry

//GoFireLogEntry -
type GoFireLogEntry struct {
	CreateDate time.Time
	Model      string
	Method     string
	Message    string
}

func NewLog(method string, model string, msg string) *GoFireLogEntry {
	return &GoFireLogEntry{
		Model:   model,
		Method:  method,
		Message: msg,
	}
}

func (l *LoggingHandler) Log(model string, method string) {

}

//UhOh - add a new GoFireError to the log and Print it to the terminal
func (g *LoggingHandler) Logs(info ...GoFireLogEntry) {
	//Local Logging
	for _, entry := range info {
		g.insertLogEntry(&entry)
		json, err := utilities.StringJSON(entry)
		if err != nil {
			Geoffrey.Handle.UhOh(-1, "GoFireErrors", "UhOh", "unable to json", err)
		}
		log.Println(json)
	}
}

//insertLogEntry -
func (g *LoggingHandler) insertLogEntry(notGoodThing *GoFireLogEntry) {
	l := len(*g)
	target := *g
	if cap(*g) == l {
		target = make([]GoFireLogEntry, l+1, l+10)
		copy(target, *g)
		target[l] = *notGoodThing
	} else {
		target = append(target, *notGoodThing)
	}
	g = &target
}
