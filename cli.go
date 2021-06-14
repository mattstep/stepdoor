package main

import (
	"fmt"
	"github.com/alecthomas/kong"
	"stepdoor/stepdoor"
)

var CLI struct {
	Open struct {} `cmd: help:"Open the door"`
	Close struct {} `cmd: help:"Close the door"`
	Interrupt struct{} `cmd: help:"Interrupt an existing operation on the door"`
	State struct{} `cmd: help:"Current state of the door"`
}

var pinMapping = stepdoor.DoorPinMapping{
	TopLimitSwitchPin:        17,
	BottomLimitSwitchPin:     18,
	StepperMotorDirectionPin: 4,
	StepperMotorStepPin:      14,
	StepperMotorSleepPin:     15,
}

func main() {
	ctx := kong.Parse(&CLI)

	switch ctx.Command() {
	case "open" :
		err := stepdoor.NewStepDoor(pinMapping).Open()
		if err != nil {
			panic(err)
		}
		fmt.Printf("Door has been opened.")
	case "close" :
		err := stepdoor.NewStepDoor(pinMapping).Close()
		if err != nil {
			panic(err)
		}
		fmt.Printf("Door has been closed.")
	case "interrupt" :
		door := stepdoor.NewStepDoor(pinMapping)
		door.Interrupt()
		fmt.Printf("Door has been interrupted. Current state is %v", door.Current())
	case "state" :
		fmt.Printf("Door's current state is %v", stepdoor.NewStepDoor(pinMapping).Current())
	default:
		panic(ctx.Command())
	}
}
