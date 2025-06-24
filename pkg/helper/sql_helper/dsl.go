package sql_helper

import (
	"fmt"
	"github.com/ledaian41/go-cms-service/pkg/shared/dto"
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
}

type JoinSpec struct {
	Tables []JoinTable
}

func NewJoinSpec(tid string, propertyTypes []shared_dto.PropertyTypeDTO) *JoinSpec {
	spec := &JoinSpec{
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
					Right: fmt.Sprintf("%s.id", pt.ReferenceType),
				},
			},
		}
		spec.Tables = append(spec.Tables, joinTable)
	}
	return spec
}

func QueryJoin(spec JoinSpec) string {
	query := ""
	for _, table := range spec.Tables {
		query += fmt.Sprintf("%s JOIN %s ON ", table.JoinType, table.Name)
		conditions := make([]string, 0)

		for _, cond := range table.Conditions {
			conditions = append(conditions,
				fmt.Sprintf("%s %s %s", cond.Left, cond.Op, cond.Right))
		}

		query += strings.Join(conditions, " AND ")
	}
	return query
}
