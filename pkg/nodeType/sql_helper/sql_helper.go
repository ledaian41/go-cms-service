package sql_helper

import (
	"fmt"
	"github.com/bwmarrin/snowflake"
	"go-cms-service/pkg/nodetype/model"
	"go-cms-service/pkg/valuetype"
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
		sqlType := valuetype.MapValueTypeToSQL(pt)
		if len(sqlType) == 0 {
			continue
		}
		columnDefs = append(columnDefs, fmt.Sprintf("%s %s", pt.PID, sqlType))
	}
	query += strings.Join(columnDefs, ", ") + ");"
	return query
}

func QueryAddColumnToTable(tid string, pt *nodeType_model.PropertyType) string {
	sqlType := valuetype.MapValueTypeToSQL(pt)
	if len(sqlType) == 0 {
		return ""
	}
	query := fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s;",
		tid, pt.PID, sqlType)
	return query
}

func QueryDeleteColumnFromTable(tid, pid string) string {
	return fmt.Sprintf("ALTER TABLE %s DROP COLUMN %s", tid, pid)
}

func QueryTableColumns(tid string) string {
	return fmt.Sprintf("PRAGMA table_info(%s)", tid)
}
