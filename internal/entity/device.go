package entity

import (
	"github.com/google/uuid"
)

type Device struct {
	ID               uuid.UUID
	Name             string
	Description      string
	PlantType        string
	IsOnline         bool
	AnalogSensors    []analogSensor
	DiscreteSensors  []discreteSensor
	DiscreteActuator []discreteActuator
}

type analogSensor struct {
	ID          uuid.UUID
	Name        string
	Description string
	Type        string
	Unit        string
	Value       int
	Scale       int
}

type discreteSensor struct {
	ID          uuid.UUID
	Name        string
	Description string
	Value       bool
}

type discreteActuator struct {
	ID          uuid.UUID
	Name        string
	Description string
	Value       bool
}
