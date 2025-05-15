package valuetype

func MapValueTypeToSQL(valueType string) string {
	vt, err := ParseValueType(valueType)
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
