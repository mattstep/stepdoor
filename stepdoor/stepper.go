package stepdoor

import (
	"fmt"
	"github.com/stianeikeland/go-rpio/v4"
	"time"
)

type Stepper interface {
	Step(count int) LimitError
	Clockwise()
	CounterClockwise()
	Sleep()
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

const SteppingPeriod = 2 * time.Microsecond;

// 4x microstepping * 580mm / (1.8deg/step * 8mm/360deg) = 240,000 microsteps
const maxSteps = 58000

func NewStepperMotor(stepPinNumber int, directionPinNumber int, sleepPinNumber int) *StepperMotor {
	sm := &StepperMotor{
		stepPin:      rpio.Pin(stepPinNumber),
		directionPin: rpio.Pin(directionPinNumber),
		sleepPin:	  rpio.Pin(sleepPinNumber),
		direction: Clockwise,
		steps: 0,
	}

	sm.stepPin.Output()
	sm.stepPin.Low()

	sm.directionPin.Output()
	sm.directionPin.Low()

	sm.sleepPin.Output()
	sm.sleepPin.High()

	return sm
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
			s.steps++
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
