package shared_interface

import (
	"go-cms-service/pkg/shared/dto"
)

type NodeTypeServiceInterface interface {
	FetchNodeTypes() *[]shared_dto.NodeTypeDTO
	LoadSchema(filePath string, ch chan<- string)
	DeleteNodeType(tid string) (bool, error)
	CheckNodeTypeExist(tid string) bool
	FetchRecords(tid string) (*[]map[string]interface{}, error)
	FetchRecord(tid string, id string) (*map[string]interface{}, error)
	CreateRecord(tid string, data map[string]interface{}) (*map[string]interface{}, error)
	UpdateRecord(tid string) (*map[string]interface{}, error)
	DeleteRecord(tid string, id string) error
}
