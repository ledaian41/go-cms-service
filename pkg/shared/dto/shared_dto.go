package shared_dto

import "math"

type NodeTypeDTO struct {
	TID           string            `json:"tid"`
	PropertyTypes []PropertyTypeDTO `json:"propertyTypes"`
}

type PropertyTypeDTO struct {
	PID            string `json:"pid"`
	ValueType      string `json:"valueType"`
	ReferenceType  string `json:"referenceType"`
	ReferenceValue string `json:"referenceValue"`
}

type PaginationDTO struct {
	Page      int32 `json:"page"`
	PageSize  int8  `json:"pageSize"`
	Total     int64 `json:"total"`
	TotalPage int   `json:"totalPage"`
}

func (pagination *PaginationDTO) CalculateTotalPage() {
	pagination.TotalPage = int(math.Ceil(float64(pagination.Total) / float64(pagination.PageSize)))
}
