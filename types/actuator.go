package types

type Actuator struct {
	Name     string        `json:"name" validate:"required"`
	IP       string        `json:"ip" validate:"required,ip"`
	SingType string        `json:"singType" validate:"required,singType"`
	Tags     []ActuatorTag `json:"tags"`
}

type ActuatorTag struct {
	Tag       string  `json:"tag" validate:"required"`
	HighLimit float32 `json:"highLimit" validate:"required"`
	LowLimit  float32 `json:"low" validate:"required"`
	HighAjust float32 `json:"highAjust" validate:"required"`
	LowAjust  float32 `json:"lowAjust" validate:"required"`
}
