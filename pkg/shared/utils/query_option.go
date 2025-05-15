package shared_utils

import "strings"

type QueryOption struct {
	ReferenceView string
	Page          int32
	PageSize      int8
}

func (qo QueryOption) GetReferenceViewKeys() []string {
	if len(qo.ReferenceView) == 0 {
		return nil
	}
	return strings.Split(qo.ReferenceView, ",")
}
