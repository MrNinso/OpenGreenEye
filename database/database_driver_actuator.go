//+build actuator

package database

import "OpenGreenEye/types"

type DriverActuator interface {
	DriverGlobal
	RegisterActuator(actuator types.Actuator) (token string)
	UpdateActuator(actuator types.Actuator) (error uint8)
	LogActuator(id int, tag, limitType string)
	ListActuators() []types.Actuator
	ListActuatorsByTags(tags []string) []types.ActuatorTag
}
