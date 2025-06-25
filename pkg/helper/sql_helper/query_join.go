package sql_helper

import (
	"encoding/json"
	"fmt"
	"github.com/iancoleman/strcase"
	"github.com/ledaian41/go-cms-service/pkg/shared/dto"
	shared_utils "github.com/ledaian41/go-cms-service/pkg/shared/utils"
	"strings"
)

type JoinCondition struct {
	Left  string
	Op    string
	Right string
}

type JoinTable struct {
	Name       string
	JoinType   string
	Conditions []JoinCondition
	Alias      string
	Fields     []string
}

type JoinSpec struct {
	Tables []JoinTable
}

func NewJoinSpec(tid string, propertyTypes []shared_dto.PropertyTypeDTO) JoinSpec {
	spec := JoinSpec{
		Tables: make([]JoinTable, 0, len(propertyTypes)),
	}
	for _, pt := range propertyTypes {
		if len(pt.ReferenceType) == 0 {
			continue
		}

		joinTable := JoinTable{
			Name:     pt.ReferenceType,
			JoinType: "LEFT",
			Conditions: []JoinCondition{
				{
					Left:  fmt.Sprintf("%s.%s", tid, pt.PID),
					Op:    "=",
					Right: fmt.Sprintf("%s.id", pt.PID),
				},
			},
			Alias:  pt.PID,
			Fields: []string{},
		}
		spec.Tables = append(spec.Tables, joinTable)
	}
	return spec
}

func QueryJoin(spec JoinSpec) string {
	query := ""
	for _, table := range spec.Tables {
		query += fmt.Sprintf("%s JOIN %s AS %s ON ", table.JoinType, strcase.ToSnake(table.Name), table.Alias)
		conditions := make([]string, 0)

		for _, cond := range table.Conditions {
			conditions = append(conditions,
				fmt.Sprintf("%s %s %s", cond.Left, cond.Op, cond.Right))
		}

		query += strings.Join(conditions, " AND ")
	}
	return query
}

func BuildSelectFields(typeId string, spec JoinSpec) string {
	var fields []string

	fields = append(fields, fmt.Sprintf("%s.*", typeId))

	for _, table := range spec.Tables {
		if len(table.Fields) == 0 {
			fields = append(fields, fmt.Sprintf("row_to_json(%s.*) as %s", table.Alias, table.Alias))
		} else {
			for _, field := range table.Fields {
				fields = append(fields, fmt.Sprintf("%s.%s as %s_%s",
					table.Alias,
					field,
					table.Alias,
					field))
			}
		}
	}

	return strings.Join(fields, ", ")

}

func FormatJoinResponse(records []map[string]interface{}, spec JoinSpec) []map[string]interface{} {
	formattedRecords := make([]map[string]interface{}, len(records))
	for i, record := range records {
		formattedRecord := make(map[string]interface{})

		for key, value := range record {
			handled := false

			for _, join := range spec.Tables {
				prefix := join.Alias + "_"
				if strings.HasPrefix(key, prefix) {
					tableName := join.Name
					if _, exists := formattedRecord[tableName]; !exists {
						formattedRecord[tableName] = make(map[string]interface{})
					}

					fieldName := strings.TrimPrefix(key, prefix)
					formattedRecord[tableName].(map[string]interface{})[fieldName] = value
					handled = true
					break
				}
				strValue, isString := value.(string)
				if key == join.Alias && isString && shared_utils.IsJSON(strValue) {
					if err := json.Unmarshal([]byte(strValue), &value); err == nil {
						formattedRecord[key] = value
					}
					handled = true
					break
				}
			}

			if !handled {
				formattedRecord[key] = value
			}
		}
		formattedRecords[i] = formattedRecord
	}

	return formattedRecords
}
