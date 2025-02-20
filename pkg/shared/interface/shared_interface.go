package shared_interface

import (
	"go-product-service/pkg/shared/dto"
)

type NodeTypeServiceInterface interface {
	FetchNodeTypes() *[]shared_dto.NodeTypeDTO
	LoadSchema(filePath string) (string, error)
}
