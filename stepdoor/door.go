package stepdoor

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/stianeikeland/go-rpio/v4"
	"sync"
)

type Door interface {
	Open() error
	Close() error
	Current() DoorState
	Interrupt()
}

//go:generate stringer -type=DoorState
type DoorState int

const (
	Closed DoorState = iota
	Semi
	Open
)

const StepsPerLimitCheck = 8

type StepDoor struct {
	topLimit Limit
	bottomLimit Limit
	stepper Stepper
	interrupt chan bool
	movingLock sync.Mutex
}

type DoorPinMapping struct {
	TopLimitSwitchPin int
	BottomLimitSwitchPin int
	StepperMotorDirectionPin int
	StepperMotorStepPin int
	StepperMotorSleepPin int
}

func NewStepDoor(mapping DoorPinMapping) *StepDoor {
	stepDoor := &StepDoor{
		topLimit:    NewLimit(mapping.TopLimitSwitchPin),
		bottomLimit: NewLimit(mapping.BottomLimitSwitchPin),
		stepper:     NewStepperMotor(mapping.StepperMotorStepPin, mapping.StepperMotorDirectionPin, mapping.StepperMotorSleepPin),
	}
	if err := rpio.Open(); err != nil {
		log.Errorf("unable to initialize door : %v", err)
		return nil
	}
	return stepDoor
}

func (s StepDoor) Open() error {
	s.stepper.CounterClockwise()

	err := s.moveToLimit(s.bottomLimit)

	if err != nil {
		log.Errorf("Stepper motor halted due to %v while opening", err)
		return err
	}

	s.stepper.Sleep()

	return nil
}

func (s StepDoor) Close() error {
	s.stepper.Clockwise()

	err := s.moveToLimit(s.topLimit)

	if err != nil {
		log.Errorf("Stepper motor halted due to %v while closing", err)
		return err
	}

	return nil
}

func (s StepDoor) Current() DoorState {
	if s.topLimit.AtLimit() {
		return Closed
	}
	if s.bottomLimit.AtLimit() {
		return Open
	}
	return Semi
}

func (s StepDoor) Interrupt() {
	s.interrupt <- true
}

func (s StepDoor) moveToLimit(limit Limit) error {
	s.movingLock.Lock()
	defer s.movingLock.Unlock()

	for !limit.AtLimit() {

		err := s.stepper.Step(StepsPerLimitCheck)

		if err != nil {
			return err
		}

		select {
		case interrupted := <- s.interrupt:
			if interrupted {
				return fmt.Errorf("door moving interrupted")
			}
		}
	}
	return nil
}
