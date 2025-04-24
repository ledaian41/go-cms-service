package shared_interface

import (
	file_model "go-cms-service/pkg/file/model"
	"go-cms-service/pkg/shared/dto"
	"mime/multipart"
)

type NodeTypeServiceInterface interface {
	FetchNodeTypes() *[]shared_dto.NodeTypeDTO
	FetchNodeType(tid string) shared_dto.NodeTypeDTO
	LoadSchema(filePath string, ch chan<- string)
	DeleteNodeType(tid string) (bool, error)
	CheckNodeTypeExist(tid string) bool
	FetchRecords(tid string) ([]map[string]interface{}, error)
	FetchRecord(tid string, id string) (map[string]interface{}, error)
	CreateRecord(tid string, data map[string]interface{}) (map[string]interface{}, error)
	UpdateRecord(tid string, id string, data map[string]interface{}) (map[string]interface{}, error)
	DeleteRecord(tid string, id string) error
	PreprocessData(nodeTypeDTO shared_dto.NodeTypeDTO, rawData map[string]interface{}) map[string]interface{}
}

type HelperServiceInterface interface {
	LoadJsonData(filePath string, ch chan<- string)
}

type FileService interface {
	SaveFile(file *multipart.FileHeader, uploadDir string) (*file_model.FileInfo, error)
}
