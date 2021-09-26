//+build sensors

package database

import "OpenGreenEye/types"

type DriverSensors interface {
	DriverGlobal
	RegisterSensor(sensor types.Sensor) (token string)
	RegisterSensorValue(id int, value float32) (error uint8)
	UpdateSensor(id int, sensor types.Sensor) (error uint8)
	ListSensors() []types.Sensor
}
