package stepdoor

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/stianeikeland/go-rpio/v4"
	"sync"
	"time"
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

const StepsPerLimitCheck = 4

const LoggingFrequency = 1 * time.Second

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
		stepper:     NewStepperMotor(mapping.StepperMotorStepPin, mapping.StepperMotorDirectionPin, mapping.StepperMotorSleepPin),
		topLimit:    NewLimit(mapping.TopLimitSwitchPin),
		bottomLimit: NewLimit(mapping.BottomLimitSwitchPin),
	}
	if err := rpio.Open(); err != nil {
		log.Errorf("unable to initialize door : %v", err)
		return nil
	}
	stepDoor.stepper.Start()
	stepDoor.topLimit.Start()
	stepDoor.bottomLimit.Start()
	return stepDoor
}

func (s *StepDoor) Open() error {
	ticker := time.NewTicker(LoggingFrequency)
	quit := make(chan bool)
	go s.periodicLogging(ticker, quit)

	s.stepper.CounterClockwise()

	err := s.moveToLimit(s.bottomLimit)

	quit <- true

	if err != nil {
		log.Errorf("Stepper motor halted due to %v while opening", err)
		return err
	}

	s.stepper.Sleep()

	return nil
}

func (s *StepDoor) Close() error {
	ticker := time.NewTicker(LoggingFrequency)
	quit := make(chan bool)
	go s.periodicLogging(ticker, quit)

	s.stepper.Clockwise()

	err := s.moveToLimit(s.topLimit)

	quit <- true

	if err != nil {
		log.Errorf("Stepper motor halted due to %v while closing", err)
		return err
	}

	return nil
}

func (s *StepDoor) Current() DoorState {
	if s.topLimit.AtLimit() {
		return Closed
	}
	if s.bottomLimit.AtLimit() {
		return Open
	}
	return Semi
}

func (s *StepDoor) Interrupt() {
	s.interrupt <- true
}

func (s *StepDoor) moveToLimit(limit Limit) error {
	s.movingLock.Lock()
	defer s.movingLock.Unlock()
	defer limit.Sleep()

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
		default:
			//noop
		}
	}
	return nil
}

func (s *StepDoor) periodicLogging(ticker *time.Ticker, quit chan bool) {
	for {
		select {
		case <-ticker.C:
			s.stepper.LogSteps()
		case <-quit:
			ticker.Stop()
			return
		}
	}
}
