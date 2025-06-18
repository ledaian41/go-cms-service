package shared_utils

import (
	"net/url"
	"strconv"
	"strings"
)

type SearchQuery struct {
	Field    string
	Operator string
	Value    string
}

type QueryOption struct {
	ReferenceView string
	Page          int32
	PageSize      int8
	SortBy        string
	Query         url.Values
}

var validOperators = map[string]bool{
	"equal":   true,
	"include": true,
	"in":      true,
}

func (qo QueryOption) GetReferenceViewKeys() []string {
	if len(qo.ReferenceView) == 0 {
		return nil
	}
	return strings.Split(qo.ReferenceView, ",")
}

func (qo QueryOption) GetSearchQuery() []SearchQuery {
	queries := make([]SearchQuery, 0)
	for key, values := range qo.Query {
		if len(values) == 0 {
			continue
		}

		parts := strings.Split(key, "_")
		if len(parts) < 2 {
			continue
		}

		field := parts[0]
		operator := parts[len(parts)-1]

		if !validOperators[operator] {
			continue
		}

		queries = append(queries, SearchQuery{
			Field:    field,
			Operator: operator,
			Value:    values[0],
		})
	}
	return queries
}

func ParseInt(value string) int {
	if result, err := strconv.Atoi(value); err == nil {
		return result
	}
	return 0
}
