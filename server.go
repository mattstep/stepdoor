package main

import (
  "github.com/gin-gonic/gin"
  . "net/http"
  "stepdoor/stepdoor"
)

func main() {
  door := stepdoor.NewStepDoor(stepdoor.DoorPinMapping{
    TopLimitSwitchPin:        17,
    BottomLimitSwitchPin:     18,
    StepperMotorDirectionPin: 4,
    StepperMotorStepPin:      14,
    StepperMotorSleepPin:     15,
  })
  
  r:= gin.Default()

  r.GET("/state", func(c *gin.Context) {
    doorStatusResponse(c, door)
  })

  r.POST("/open", func(c *gin.Context) {
    if err := door.Close(); err != nil {
      doorErrorResponse(c, err)
    } else {
      doorStatusResponse(c, door)
    }
  })

  r.POST("/close", func(c *gin.Context) {
    if err := door.Close(); err != nil {
      doorErrorResponse(c, err)
    } else {
      doorStatusResponse(c, door)
    }
  })

  r.POST("/interrupt", func(c *gin.Context) {
    door.Interrupt()
    doorStatusResponse(c, door)
  })

  r.Run()
}

func doorErrorResponse(c *gin.Context, err error) {
  c.Error(err)
  c.String(StatusInternalServerError, err.Error())
}

func doorStatusResponse(c *gin.Context, door *stepdoor.StepDoor) {
  c.String(StatusOK, door.Current().String())
}