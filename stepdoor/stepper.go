package stepdoor

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/stianeikeland/go-rpio/v4"
	"time"
)

type Stepper interface {
	Start()
	Step(count int) LimitError
	Clockwise()
	CounterClockwise()
	Sleep()
	LogSteps()
}

type StepperMotor struct {
	stepPin rpio.Pin
	directionPin rpio.Pin
	sleepPin rpio.Pin
	direction Direction
	steps int
	asleep bool
}

type Direction int8

type LimitError error

const (
	Clockwise Direction = iota
	CounterClockwise
)

const SteppingPeriod = 1 * time.Millisecond

// 4x microstepping * 580mm / (1.8deg/step * 8mm/360deg) = 240,000 microsteps
const maxSteps = 58000

func (s StepperMotor) LogSteps() {
	log.WithFields(log.Fields{
		"steps": s.steps,
	}).Info("current step count")
}

func NewStepperMotor(stepPinNumber int, directionPinNumber int, sleepPinNumber int) *StepperMotor {
	sm := &StepperMotor{
		stepPin:      rpio.Pin(stepPinNumber),
		directionPin: rpio.Pin(directionPinNumber),
		sleepPin:	  rpio.Pin(sleepPinNumber),
		direction: Clockwise,
		steps: 0,
	}

	return sm
}

func (s StepperMotor) Start() {
	s.stepPin.Output()
	s.stepPin.High()

	s.directionPin.Output()
	s.directionPin.Low()

	s.sleepPin.Output()
	s.sleepPin.High()
}

func (s StepperMotor) Step(count int) LimitError {
	if s.asleep {
		s.wake()
	}
	if s.steps + count < maxSteps {
		for i:=0; i < count; i++ {
			s.stepPin.Low()
			time.Sleep(SteppingPeriod / 2)
			s.stepPin.High()
			time.Sleep(SteppingPeriod / 2)
			s.steps = s.steps + 1
		}
	} else {
		return LimitError(fmt.Errorf("unable to continue stepping, limit reached"))
	}

	return nil
}

func (s StepperMotor) Clockwise() {
	if s.direction == CounterClockwise {
		s.direction = Clockwise
		s.steps = 0
	}
	s.syncDirection()
}

func (s StepperMotor) CounterClockwise() {
	if s.direction == Clockwise {
		s.direction = CounterClockwise
		s.steps = 0
	}
	s.syncDirection()
}

func (s StepperMotor) syncDirection() {
	if s.direction == Clockwise {
		s.directionPin.Write(rpio.Low)
	} else {
		s.directionPin.Write(rpio.High)
	}
}

func (s StepperMotor) Sleep() {
	s.sleepPin.Low()
	s.asleep = true
}

func (s StepperMotor) wake() {
	s.sleepPin.High()
	time.Sleep(1 * time.Millisecond)
	s.asleep = false
}
