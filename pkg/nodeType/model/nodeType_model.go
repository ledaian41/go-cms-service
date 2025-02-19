package model

type NodeType struct {
	TID           string         `json:"tid"`
	PropertyTypes []PropertyType `json:"propertyTypes"`
}

type PropertyType struct {
	PID       string `json:"pid"`
	ValueType string `json:"valueType"`
}
