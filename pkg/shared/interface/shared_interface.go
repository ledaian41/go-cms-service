package shared_interface

import (
	"go-cms-service/pkg/shared/dto"
)

type NodeTypeServiceInterface interface {
	FetchNodeTypes() *[]shared_dto.NodeTypeDTO
	LoadSchema(filePath string, ch chan<- string)
	DeleteNodeType(tid string) (bool, error)
}
