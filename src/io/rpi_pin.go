package io

import (
	"firecontroller/utilities"
	"strconv"
)

/*
  +-----+---------+----------+---------+-----+
  | BCM |   Name  | Physical | Name    | BCM |
  +-----+---------+----++----+---------+-----+
  |     |    3.3v |  1 || 2  | 5v      |     |
  |   2 |   SDA 1 |  3 || 4  | 5v      |     |
  |   3 |   SCL 1 |  5 || 6  | 0v      |     |
  |   4 | GPIO  7 |  7 || 8  | TxD     | 14  |
  |     |      0v |  9 || 10 | RxD     | 15  |
  |  17 | GPIO  0 | 11 || 12 | GPIO  1 | 18  |
  |  27 | GPIO  2 | 13 || 14 | 0v      |     |
  |  22 | GPIO  3 | 15 || 16 | GPIO  4 | 23  |
  |     |    3.3v | 17 || 18 | GPIO  5 | 24  |
  |  10 |    MOSI | 19 || 20 | 0v      |     |
  |   9 |    MISO | 21 || 22 | GPIO  6 | 25  |
  |  11 |    SCLK | 23 || 24 | CE0     | 8   |
  |     |      0v | 25 || 26 | CE1     | 7   |
  |   0 |   SDA 0 | 27 || 28 | SCL 0   | 1   |
  |   5 | GPIO 21 | 29 || 30 | 0v      |     |
  |   6 | GPIO 22 | 31 || 32 | GPIO 26 | 12  |
  |  13 | GPIO 23 | 33 || 34 | 0v      |     |
  |  19 | GPIO 24 | 35 || 36 | GPIO 27 | 16  |
  |  26 | GPIO 25 | 37 || 38 | GPIO 28 | 20  |
  |     |      0v | 39 || 40 | GPIO 29 | 21  |
  +-----+---------+----++----+---------+-----+
*/

//RpiPinMap -
type RpiPinMap struct {
	//BcmPin - The pin id on the processor, used by rpio for commands
	BcmPin uint8
	//Human readable name for the Raspi Pin
	Name string
	//Connection position on the header according to the block above^^
	HeaderPin int
}

//NoPin is a placeholder value for a pin that should not be assigned to a component
var NoPin uint8 = 255

func (r RpiPinMap) String() string {
	return "\t" + utilities.LabelString("Name", r.Name) + "\t" +
		utilities.LabelString("Header Pin", strconv.Itoa(r.HeaderPin)) + "\t" +
		utilities.LabelString("BCM Pin", strconv.Itoa(int(r.BcmPin)))
}

//GetPins - Returns Pins for Raspi 4
func GetPins() []RpiPinMap {
	return []RpiPinMap{
		{
			HeaderPin: 1,
			BcmPin:    NoPin,
			Name:      "3.3v",
		},
		{
			HeaderPin: 3,
			BcmPin:    2,
			Name:      "BCM 2/SDA 1",
		},
		{
			HeaderPin: 5,
			BcmPin:    3,
			Name:      "BCM 3/SCL 1",
		},
		{
			HeaderPin: 7,
			BcmPin:    4,
			Name:      "BCM 4/GPIOCLK0",
		},
		{
			HeaderPin: 9,
			BcmPin:    NoPin,
			Name:      "0v",
		},
		{
			HeaderPin: 11,
			BcmPin:    17,
			Name:      "BCM 17",
		},
		{
			HeaderPin: 13,
			BcmPin:    27,
			Name:      "BCM 27",
		},
		{
			HeaderPin: 15,
			BcmPin:    22,
			Name:      "BCM 22",
		},
		{
			HeaderPin: 17,
			BcmPin:    NoPin,
			Name:      "3.3v",
		},
		{
			HeaderPin: 19,
			BcmPin:    10,
			Name:      "BCM 10/MOSI",
		},
		{
			HeaderPin: 21,
			BcmPin:    9,
			Name:      "BCM 9/MISO",
		},
		{
			HeaderPin: 23,
			BcmPin:    11,
			Name:      "BCM 11/SCLK",
		},
		{
			HeaderPin: 25,
			BcmPin:    NoPin,
			Name:      "0v",
		},
		{
			HeaderPin: 27,
			BcmPin:    0,
			Name:      "BCM 0/SDA 0",
		},
		{
			HeaderPin: 29,
			BcmPin:    5,
			Name:      "BCM 5",
		},
		{
			HeaderPin: 31,
			BcmPin:    6,
			Name:      "BCM 6",
		},
		{
			HeaderPin: 33,
			BcmPin:    13,
			Name:      "BCM 13/PWM1",
		},
		{
			HeaderPin: 35,
			BcmPin:    19,
			Name:      "BCM 19/MISO",
		},
		{
			HeaderPin: 37,
			BcmPin:    26,
			Name:      "BCM 26",
		},
		{
			HeaderPin: 39,
			BcmPin:    NoPin,
			Name:      "0v",
		},
		{
			HeaderPin: 2,
			BcmPin:    NoPin,
			Name:      "5v",
		},
		{
			HeaderPin: 4,
			BcmPin:    NoPin,
			Name:      "5v",
		},
		{
			HeaderPin: 6,
			BcmPin:    NoPin,
			Name:      "0v",
		},
		{
			HeaderPin: 8,
			BcmPin:    14,
			Name:      "BCM 14/TxD",
		},
		{
			HeaderPin: 10,
			BcmPin:    15,
			Name:      "BCM 15/RxD",
		},
		{
			HeaderPin: 12,
			BcmPin:    18,
			Name:      "BCM 18/PWM0",
		},
		{
			HeaderPin: 14,
			BcmPin:    NoPin,
			Name:      "0v",
		},
		{
			HeaderPin: 16,
			BcmPin:    23,
			Name:      "BCM 23",
		},
		{
			HeaderPin: 18,
			BcmPin:    24,
			Name:      "BCM 24",
		},
		{
			HeaderPin: 20,
			BcmPin:    NoPin,
			Name:      "0v",
		},
		{
			HeaderPin: 22,
			BcmPin:    25,
			Name:      "BCM 25",
		},
		{
			HeaderPin: 24,
			BcmPin:    8,
			Name:      "BCM 8/CE0",
		},
		{
			HeaderPin: 26,
			BcmPin:    7,
			Name:      "BCM 7/CE1",
		},
		{
			HeaderPin: 28,
			BcmPin:    1,
			Name:      "BCM 1/ID_SC",
		},
		{
			HeaderPin: 30,
			BcmPin:    NoPin,
			Name:      "0v",
		},
		{
			HeaderPin: 32,
			BcmPin:    12,
			Name:      "BCM 12/PWM0",
		},
		{
			HeaderPin: 34,
			BcmPin:    NoPin,
			Name:      "0v",
		},
		{
			HeaderPin: 36,
			BcmPin:    16,
			Name:      "BCM 16",
		},
		{
			HeaderPin: 38,
			BcmPin:    20,
			Name:      "BCM 20/MOSI",
		},
		{
			HeaderPin: 40,
			BcmPin:    21,
			Name:      "BCM 21/SCLK",
		},
	}
}
