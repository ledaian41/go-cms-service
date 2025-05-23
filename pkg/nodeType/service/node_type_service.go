package nodeType_service

import (
	"go-cms-service/pkg/nodetype/model"
	"go-cms-service/pkg/nodetype/sql_helper"
	"go-cms-service/pkg/shared/dto"
	"go-cms-service/pkg/shared/interface"
	"go-cms-service/pkg/shared/utils"
	"gorm.io/gorm"
	"log"
)

type NodeTypeService struct {
	db          *gorm.DB
	fileService shared_interface.FileService
}

func NewNodeTypeService(db *gorm.DB, fileService shared_interface.FileService) *NodeTypeService {
	return &NodeTypeService{db: db, fileService: fileService}
}

func (s *NodeTypeService) InitDatabase() {
	err := s.db.AutoMigrate(&nodeType_model.NodeType{}, &nodeType_model.PropertyType{})
	if err != nil {
		log.Printf("❌ Failed at AutoMigrate: %v", err)
	}
	log.Println("🎉 NodeType - Database migrate successfully")
}

func (s *NodeTypeService) FetchNodeTypes() *[]shared_dto.NodeTypeDTO {
	var nodeTypes []nodeType_model.NodeType
	if err := s.db.Preload("PropertyTypes").Find(&nodeTypes).Error; err != nil {
		log.Printf("❌ Failed at query NodeTypes: %v", err)
	}
	var dtos []shared_dto.NodeTypeDTO
	for _, n := range nodeTypes {
		dtos = append(dtos, n.NodeTypeDTO())
	}
	return &dtos
}

func (s *NodeTypeService) FetchNodeType(tid string) shared_dto.NodeTypeDTO {
	var node nodeType_model.NodeType
	if err := s.db.Preload("PropertyTypes").Where("tid = ?", tid).First(&node).Error; err != nil {
		log.Printf("❌ Failed at query NodeTypes: %v", err)
	}
	return node.NodeTypeDTO()
}

func (s *NodeTypeService) DeleteNodeType(tid string) (bool, error) {
	var node nodeType_model.NodeType
	if err := s.db.Where("tid = ?", tid).First(&node).Error; err != nil {
		log.Printf("❌ NodeType not found: %v", err)
		return false, err
	}

	if err := s.db.Unscoped().Where("node_type_refer = ?", node.ID).Delete(&nodeType_model.PropertyType{}).Error; err != nil {
		log.Printf("❌ Failed to delete PropertyType: %v", err)
	}

	if err := s.db.Unscoped().Delete(&node).Error; err != nil {
		log.Printf("❌ Failed at delete NodeType: %v", err)
		return false, err
	}
	return true, nil
}

func (s *NodeTypeService) CheckNodeTypeExist(tid string) bool {
	return s.db.Migrator().HasTable(tid)
}

func (s *NodeTypeService) FetchRecords(tid string, option shared_utils.QueryOption) ([]map[string]interface{}, *shared_dto.PaginationDTO, error) {
	var records []map[string]interface{}
	if option.Page < 1 {
		option.Page = 1 // default page number
	}

	if option.PageSize < 1 || option.PageSize > 100 {
		option.PageSize = 10 // default page size
	}
	offset := (int(option.Page) - 1) * int(option.PageSize)
	if err := s.db.Table(tid).Limit(int(option.PageSize)).Offset(offset).Find(&records).Error; err != nil {
		return nil, nil, err
	}

	var total int64
	if err := s.db.Table(tid).Count(&total).Error; err != nil {
		return nil, nil, err
	}
	pagination := &shared_dto.PaginationDTO{
		Page:     option.Page,
		PageSize: option.PageSize,
		Total:    total,
	}
	pagination.CalculateTotalPage()

	return records, pagination, nil
}

func (s *NodeTypeService) FetchRecord(tid string, id string) (map[string]interface{}, error) {
	var result map[string]interface{}
	if err := s.db.Table(tid).Find(&result, "id = ?", id).Error; err != nil {
		return nil, err
	}
	if result == nil {
		return nil, nil
	}
	return result, nil
}

func (s *NodeTypeService) CreateRecord(tid string, data map[string]interface{}) (map[string]interface{}, error) {
	data["id"] = sql_helper.GenerateID()
	if result := s.db.Table(tid).Create(&data); result.Error != nil {
		return data, result.Error
	}
	delete(data, "@id")
	return data, nil
}

func (s *NodeTypeService) UpdateRecord(tid string, id string, data map[string]interface{}) (map[string]interface{}, error) {
	delete(data, "id")
	result := s.db.Table(tid).Where("id = ?", id).Updates(&data)
	if result.Error != nil {
		return nil, result.Error
	}
	return data, nil
}

func (s *NodeTypeService) DeleteRecord(tid string, id string) error {
	return s.db.Table(tid).Where("id = ?", id).Delete(nil).Error
}
