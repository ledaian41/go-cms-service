package sql_helper

import (
	"fmt"
	"github.com/bwmarrin/snowflake"
	"github.com/iancoleman/strcase"
	"github.com/ledaian41/go-cms-service/pkg/node_type/model"
	"github.com/ledaian41/go-cms-service/pkg/shared/utils"
	"github.com/ledaian41/go-cms-service/pkg/value_type"
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

func QueryCreateNewTable(nodeType *node_type_model.NodeType) string {
	query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (id text PRIMARY KEY, created_at timestamptz, created_by, modified_at timestamptz, modified_by text, deleted_at timestamptz, deleted_by text, ", strcase.ToSnake(nodeType.TID))
	var columnDefs []string

	for _, pt := range nodeType.PropertyTypes {
		sqlType := value_type.MapValueTypeToSQL(pt)
		if len(sqlType) == 0 {
			continue
		}
		columnDefs = append(columnDefs, fmt.Sprintf("%s %s", strcase.ToSnake(pt.PID), sqlType))
	}
	query += strings.Join(columnDefs, ", ") + ");"
	return query
}

func QueryAddColumnToTable(tid string, pt *node_type_model.PropertyType) string {
	sqlType := value_type.MapValueTypeToSQL(pt)
	if len(sqlType) == 0 {
		return ""
	}
	query := fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s;",
		tid, pt.PID, sqlType)
	return query
}

func QueryDeleteColumnFromTable(tid, pid string) string {
	return fmt.Sprintf("ALTER TABLE %s DROP COLUMN %s", strcase.ToSnake(tid), strcase.ToSnake(pid))
}

func QueryTableColumns(tid string) string {
	return fmt.Sprintf(`
		SELECT column_name
		FROM information_schema.columns
		WHERE table_schema = 'public' AND table_name = '%s'
		ORDER BY ordinal_position
	`, tid)
}

func QueryTableExist(tid string) string {
	return fmt.Sprintf(`
		SELECT EXISTS (
			SELECT FROM information_schema.tables
			WHERE table_schema = 'public' AND table_name = '%s'
		)
	`, tid)
}

func BuildSearchConditions(queries []shared_utils.SearchQuery) (string, []interface{}) {
	var conditions []string
	var values []interface{}

	for _, query := range queries {
		switch query.Operator {
		case "equal":
			conditions = append(conditions, fmt.Sprintf("%s = ?", query.Field))
			values = append(values, query.Value)
		case "include":
			conditions = append(conditions, fmt.Sprintf("%s ILIKE ?", query.Field))
			values = append(values, fmt.Sprintf("%%%s%%", strings.TrimSpace(query.Value)))
		case "in":
			conditions = append(conditions, fmt.Sprintf("%s IN ?", query.Field))
			queryValue := strings.Split(query.Value, ",")
			interfaceValues := make([]interface{}, len(queryValue))
			for i, v := range queryValue {
				interfaceValues[i] = v
			}
			values = append(values, interfaceValues)
		case "from":
			conditions = append(conditions, fmt.Sprintf("%s >= ?", query.Field))
			values = append(values, query.Value)
		case "to":
			conditions = append(conditions, fmt.Sprintf("%s <= ?", query.Field))
			values = append(values, query.Value)
		case "fromto":
			fromTo := strings.Split(query.Value, ",")
			conditions = append(conditions, fmt.Sprintf("%s BETWEEN ? AND ?", query.Field))
			values = append(values, strings.TrimSpace(fromTo[0]), strings.TrimSpace(fromTo[1]))
		}
	}

	if len(conditions) == 0 {
		return "", nil
	}

	return strings.Join(conditions, " AND "), values
}
