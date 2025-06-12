package valuetype

import "go-cms-service/pkg/nodetype/model"

func MapValueTypeToSQL(pt *nodeType_model.PropertyType) string {
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
