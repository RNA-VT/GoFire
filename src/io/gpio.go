package io

import (
	"errors"
	"firecontroller/io/mock"
	"firecontroller/utilities"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/spf13/viper"
	"github.com/stianeikeland/go-rpio/v4"
)

//Gpio - a raspberrypi digital gpio pin
type Gpio struct {
	Pin     GpioPin
	PinInfo RpiPinMap
	Failed  bool
}

//GpioPin - an interface to support mocking rpio.Pin
type GpioPin interface {
	Output()
	High()
	Low()
	Toggle()
	Write(rpio.State)
	Read() rpio.State
}

func (g Gpio) String() string {
	var pinString string
	pinString, err := utilities.StringJSON(g.Pin)
	if err != nil {
		pinString = err.Error()
	}
	return utilities.LabelString("\tFAILED", strconv.FormatBool(g.Failed)) +
		utilities.LabelString("\tPin", pinString) +
		utilities.LabelString("\tPin Info", g.PinInfo.String())
}

//Init - create gpio pin object and set modes
func (g *Gpio) Init(headerPin int, initHigh bool) error {
	if err := g.loadPinInfoByHeader(headerPin); err != nil {
		return err
	}
	log.Println("BCM Pin:", g.PinInfo.BcmPin)
	//This pin theoretically checks out, but is it real?
	thisIsNotReal := viper.GetBool("GOFIRE_MOCK_GPIO") || viper.GetString("ENV") == "local"
	if thisIsNotReal {
		//Nothing is real and this pin, especially, is laughable. Mock it.
		g.Pin = mock.Pin(g.PinInfo.BcmPin)
	} else {
		g.Pin = rpio.Pin(g.PinInfo.BcmPin)
		if err := rpio.Open(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	g.Pin.Output()
	if initHigh {
		g.Pin.High()
	} else {
		g.Pin.Low()
	}
	if !thisIsNotReal {
		rpio.Close()
	}

	return nil
}

//HandleEnable -
func (g *Gpio) HandleEnable() bool {
	if err := rpio.Open(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer rpio.Close()
	g.Pin.Low()
	return true
}

//HandleDisable -
func (g *Gpio) HandleDisable() bool {
	if err := rpio.Open(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer rpio.Close()
	g.Pin.Low()
	return true
}

/*CurrentStateString just for pretty printing the device info */
func (g *Gpio) CurrentStateString() string {
	state := "OFF"

	if g.Pin.Read() == rpio.High {
		state = "ON"
	}
	return state
}

func (g *Gpio) loadPinInfoByHeader(headerPin int) error {
	pins := GetPins()
	for i := 0; i < len(pins); i++ {
		if pins[i].HeaderPin == headerPin {
			if pins[i].BcmPin != NoPin {
				g.Failed = false
				g.PinInfo = pins[i]
				return nil
			}
			g.Failed = true
			return errors.New("[GPIO PIN INIT]: Header Pin:" + strconv.Itoa(headerPin) + " is not a GPIO pin. Type: " + pins[i].Name)
		}
	}
	g.Failed = true
	return errors.New("[GPIO PIN INIT]: Header Pin:" + strconv.Itoa(headerPin) + "not found.")
}
