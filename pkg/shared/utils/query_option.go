package shared_utils

import (
	"strconv"
	"strings"
)

type QueryOption struct {
	ReferenceView string
	Page          int32
	PageSize      int8
	SortBy        string
}

func (qo QueryOption) GetReferenceViewKeys() []string {
	if len(qo.ReferenceView) == 0 {
		return nil
	}
	return strings.Split(qo.ReferenceView, ",")
}

func ParseInt(value string) int {
	if result, err := strconv.Atoi(value); err == nil {
		return result
	}
	return 0
}
