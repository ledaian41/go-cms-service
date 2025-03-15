package dynamic_struct

import (
	"fmt"
	"github.com/bwmarrin/snowflake"
	"go-cms-service/pkg/nodeType/model"
	"hash/fnv"
	"math/big"
	"strings"
)

func MapValueTypeToSQL(valueType string) string {
	switch valueType {
	case "STRING":
		return "text"
	case "INT":
		return "integer"
	case "DOUBLE", "FLOAT":
		return "real"
	case "BOOLEAN":
		return "integer"
	default:
		return "text"
	}
}

func generateSnowflakeID() int64 {
	node, _ := snowflake.NewNode(1)
	return node.Generate().Int64()
}

func GenerateID() string {
	id := generateSnowflakeID()
	h := fnv.New32a()
	h.Write([]byte(fmt.Sprintf("%d", id)))
	hashed := h.Sum32()
	return new(big.Int).SetUint64(uint64(hashed)).Text(36)
}

func CreateDynamicTable(nodeType *nodeType_model.NodeType) string {
	query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (id text PRIMARY KEY, ", nodeType.TID)
	var columnDefs []string

	for _, pt := range nodeType.PropertyTypes {
		columnDefs = append(columnDefs, fmt.Sprintf("%s %s", pt.PID, MapValueTypeToSQL(pt.ValueType)))
	}
	query += strings.Join(columnDefs, ", ") + ");"
	return query
}
