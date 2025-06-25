package node_type_service

import (
	"github.com/iancoleman/strcase"
	"github.com/ledaian41/go-cms-service/pkg/node_type/model"
	"github.com/ledaian41/go-cms-service/pkg/shared/dto"
	"github.com/ledaian41/go-cms-service/pkg/shared/interface"
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
	err := s.db.AutoMigrate(&node_type_model.NodeType{}, &node_type_model.PropertyType{})
	if err != nil {
		log.Printf("❌ Failed at AutoMigrate: %v", err)
	}
	log.Println("🎉 NodeType - Database migrate successfully")
}

func (s *NodeTypeService) FetchNodeTypes() *[]shared_dto.NodeTypeDTO {
	var nodeTypes []node_type_model.NodeType
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
	var node node_type_model.NodeType
	if err := s.db.Preload("PropertyTypes").Where("tid = ?", tid).First(&node).Error; err != nil {
		log.Printf("❌ Failed at query NodeTypes: %v", err)
	}
	return node.NodeTypeDTO()
}

func (s *NodeTypeService) DeleteNodeType(tid string) (bool, error) {
	var node node_type_model.NodeType
	if err := s.db.Where("tid = ?", tid).First(&node).Error; err != nil {
		log.Printf("❌ NodeType not found: %v", err)
		return false, err
	}

	if err := s.db.Unscoped().Where("node_type_refer = ?", node.ID).Delete(&node_type_model.PropertyType{}).Error; err != nil {
		log.Printf("❌ Failed to delete PropertyType: %v", err)
	}

	if err := s.db.Unscoped().Delete(&node).Error; err != nil {
		log.Printf("❌ Failed at delete NodeType: %v", err)
		return false, err
	}
	return true, nil
}

func (s *NodeTypeService) CheckNodeTypeExist(tid string) bool {
	var count int64
	if err := s.db.Model(&node_type_model.NodeType{}).Where("tid = ?", strcase.ToLowerCamel(tid)).Count(&count).Error; err != nil {
		return false
	}
	return count > 0
}

func (s *NodeTypeService) FetchPropertyTypesByTid(tid string) []shared_dto.PropertyTypeDTO {
	var nodeTypeId string
	s.db.Table("node_types").Select("id").Where("tid = ?", tid).Scan(&nodeTypeId)
	var propertyTypes []node_type_model.PropertyType
	s.db.Table("property_types").Where("node_type_refer = ?", nodeTypeId).Find(&propertyTypes)
	result := make([]shared_dto.PropertyTypeDTO, 0)
	for _, pt := range propertyTypes {
		result = append(result, pt.PropertyTypeDTO())
	}
	return result
}
