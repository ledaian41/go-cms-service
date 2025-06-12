package value_type

import "github.com/ledaian41/go-cms-service/pkg/node_type/model"

func MapValueTypeToSQL(pt *node_type_model.PropertyType) string {
	vt, err := ParseValueType(pt.ValueType)
	if err != nil {
		return ""
	}

	switch vt {
	case Integer, Boolean:
		return "integer"
	case Double, Float:
		return "real"
	default:
		return "text"
	}
}
