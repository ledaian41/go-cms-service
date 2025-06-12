package value_type

import "fmt"

type ValueType string

const (
	String     ValueType = "STRING"
	Integer    ValueType = "INT"
	Boolean    ValueType = "BOOLEAN"
	Double     ValueType = "DOUBLE"
	Float      ValueType = "FLOAT"
	File       ValueType = "FILE"
	Reference  ValueType = "REFERENCE"
	References ValueType = "REFERENCES"
)

var validValueTypes = map[ValueType]bool{
	String:     true,
	Integer:    true,
	Boolean:    true,
	Double:     true,
	Float:      true,
	File:       true,
	Reference:  true,
	References: true,
}

func ParseValueType(value string) (ValueType, error) {
	vt := ValueType(value)
	if !validValueTypes[vt] {
		return "", fmt.Errorf("invalid value type: %s", value)
	}
	return vt, nil
}
