package sql_helper

import (
	"fmt"
	"github.com/bwmarrin/snowflake"
	"go-cms-service/pkg/nodeType/model"
	"go-cms-service/pkg/nodeType/valueType"
	"hash/fnv"
	"math/big"
	"strings"
)

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

func QueryCreateNewTable(nodeType *nodeType_model.NodeType) string {
	query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (id text PRIMARY KEY, ", nodeType.TID)
	var columnDefs []string

	for _, pt := range nodeType.PropertyTypes {
		columnDefs = append(columnDefs, fmt.Sprintf("%s %s", pt.PID, valueType.MapValueTypeToSQL(pt.ValueType)))
	}
	query += strings.Join(columnDefs, ", ") + ");"
	return query
}

func QueryAddColumnToTable(tid string, pt *nodeType_model.PropertyType) string {
	query := fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s;",
		tid, pt.PID, valueType.MapValueTypeToSQL(pt.ValueType))
	return query
}

func QueryDeleteColumnFromTable(tid, pid string) string {
	return fmt.Sprintf("ALTER TABLE %s DROP COLUMN %s", tid, pid)
}
