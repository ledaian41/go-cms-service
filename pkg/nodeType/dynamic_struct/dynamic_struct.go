package dynamic_struct

import (
	"fmt"
	"go-cms-service/pkg/nodeType/model"
	"reflect"
	"strings"
)

func mapValueType(valueType string) reflect.Type {
	switch valueType {
	case "STRING":
		return reflect.TypeOf("")
	case "INT":
		return reflect.TypeOf(0)
	case "DOUBLE":
		return reflect.TypeOf(0.0)
	case "BOOLEAN":
		return reflect.TypeOf(true)
	default:
		return reflect.TypeOf(nil) // Giá trị không hợp lệ
	}
}

func mapValueTypeToSQL(valueType string) string {
	switch valueType {
	case "STRING":
		return "text"
	case "INT":
		return "integer"
	case "DOUBLE":
		return "real"
	case "BOOLEAN":
		return "boolean"
	default:
		return "text"
	}
}

func toExportedName(name string) string {
	if name == "" {
		return ""
	}
	return strings.ToUpper(string(name[0])) + name[1:]
}

func CreateDynamicStruct(nodeType *nodeType_model.NodeType) reflect.Type {
	fields := []reflect.StructField{
		{
			Name: "ID",
			Type: reflect.TypeOf(""),
			Tag:  reflect.StructTag(`gorm:"primaryKey;type:text"`),
		},
	}

	for _, pt := range nodeType.PropertyTypes {
		exportedName := toExportedName(pt.PID)
		field := reflect.StructField{
			Name: exportedName,
			Type: mapValueType(pt.ValueType),
			Tag:  reflect.StructTag(fmt.Sprintf(`gorm:"type:%s"`, mapValueTypeToSQL(pt.ValueType))),
		}
		fields = append(fields, field)
	}

	return reflect.StructOf(fields)
}
