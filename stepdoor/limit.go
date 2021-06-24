package stepdoor

import "github.com/stianeikeland/go-rpio/v4"


type Limit interface {
	Start()
	AtLimit() bool
	Sleep()
}

type LimitSwitch struct {
  gpioPin rpio.Pin
}

func NewLimit(pinNumber int) *LimitSwitch {
	pin := rpio.Pin(pinNumber)
	return &LimitSwitch{gpioPin: pin}
}

func (ls *LimitSwitch) Start() {
	ls.gpioPin.Input()
}

func (ls *LimitSwitch) AtLimit() bool {
	ls.gpioPin.PullUp()
	return ls.gpioPin.Read() == rpio.Low
}

func (ls *LimitSwitch) Sleep() {
	ls.gpioPin.PullOff()
}