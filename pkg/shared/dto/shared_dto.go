package shared_dto

type NodeTypeDTO struct {
	TID           string            `json:"tid"`
	PropertyTypes []PropertyTypeDTO `json:"propertyTypes"`
}

type PropertyTypeDTO struct {
	PID       string `json:"pid"`
	ValueType string `json:"valueType"`
}
