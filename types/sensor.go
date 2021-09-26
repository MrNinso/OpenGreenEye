package types

type Sensor struct {
	Name     string   `json:"name" validate:"required"`
	IP       string   `json:"ip" validate:"required,ip"`
	SingType string   `json:"singType" validate:"required,singType"`
	Tags     []string `json:"tags"`
}
