package shared_interface

import (
	"github.com/ledaian41/go-cms-service/pkg/shared/dto"
	"github.com/ledaian41/go-cms-service/pkg/shared/utils"
)

type NodeTypeService interface {
	LoadSchema(filePath string, ch chan<- string)
	FetchNodeTypes() *[]shared_dto.NodeTypeDTO
	FetchNodeType(tid string) shared_dto.NodeTypeDTO
	DeleteNodeType(tid string) (bool, error)
	CheckNodeTypeExist(tid string) bool
	FetchRecords(tid string, option shared_utils.QueryOption) ([]map[string]interface{}, *shared_dto.PaginationDTO, error)
	FetchRecord(tid string, id string) (map[string]interface{}, error)
	CreateRecord(tid string, data map[string]interface{}) (map[string]interface{}, error)
	UpdateRecord(tid string, id string, data map[string]interface{}) (map[string]interface{}, error)
	DeleteRecord(tid string, id string) error
	PreprocessFile(nodeTypeDTO shared_dto.NodeTypeDTO, rawData map[string]interface{}) (map[string]interface{}, error)
}
