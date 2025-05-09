package valueType

func MapValueTypeToSQL(valueType string) string {
	switch valueType {
	case "INT", "BOOLEAN":
		return "integer"
	case "DOUBLE", "FLOAT":
		return "real"
	case "STRING", "JSON", "FILE":
		return "text"
	default:
		return ""
	}
}
