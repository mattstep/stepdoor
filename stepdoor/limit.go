package stepdoor

import "github.com/stianeikeland/go-rpio/v4"


type Limit interface {
	AtLimit() bool
}

type LimitSwitch struct {
  gpioPin rpio.Pin
}

func NewLimit(pinNumber int) *LimitSwitch {
	pin := rpio.Pin(pinNumber)
	pin.Input()
	return &LimitSwitch{gpioPin: pin}
}

func (ls LimitSwitch) AtLimit() bool {
	ls.gpioPin.PullUp()
	return ls.gpioPin.Read() == rpio.Low
}
