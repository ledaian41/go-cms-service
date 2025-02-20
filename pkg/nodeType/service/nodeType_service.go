package nodeType_service

import (
	"fmt"
	"go-product-service/pkg/nodeType/model"
	"go-product-service/pkg/nodeType/utils"
	"go-product-service/pkg/shared/dto"
	"gorm.io/gorm"
	"log"
)

type NodeTypeService struct {
	db *gorm.DB
}

func NewNodeTypeService(db *gorm.DB) *NodeTypeService {
	return &NodeTypeService{db: db}
}

func (s *NodeTypeService) InitDatabase() {
	err := s.db.AutoMigrate(&nodeType_model.NodeType{}, &nodeType_model.PropertyType{})
	if err != nil {
		log.Fatal("❌ Failed at AutoMigrate:", err)
	}
	log.Println("🎉 NodeType - Database migrate successfully")
}

func (s *NodeTypeService) FetchNodeTypes() *[]shared_dto.NodeTypeDTO {
	var nodeTypes []nodeType_model.NodeType
	if err := s.db.Preload("PropertyTypes").Find(&nodeTypes).Error; err != nil {
		log.Fatalf("❌ Failed at query NodeTypes: %v", err)
	}
	var dtos []shared_dto.NodeTypeDTO
	for _, n := range nodeTypes {
		dtos = append(dtos, n.NodeTypeDTO())
	}
	return &dtos
}

func (s *NodeTypeService) LoadSchema(path string) (string, error) {
	newNodeType, err := nodeType_utils.LoadSchema(path)
	if err != nil {
		log.Fatalf("❌ Failed at LoadSchema: %v", err)
		return path, err
	}

	var existing nodeType_model.NodeType
	if err := s.db.Preload("PropertyTypes").Where("t_id = ?", newNodeType.TID).First(&existing).Error; err != nil {
		if err := s.db.Create(&newNodeType).Error; err != nil {
			log.Fatalf("❌ Failed at save NodeType: %v", err)
			return newNodeType.TID, err
		}
		log.Println(fmt.Sprintf("🎉 Helper - Load new %s schema successfully!", newNodeType.TID))
		return newNodeType.TID, nil
	}
	existing.PropertyTypes = newNodeType.PropertyTypes

	if err := s.db.Model(&existing).Association("PropertyTypes").Replace(newNodeType.PropertyTypes); err != nil {
		log.Fatalf("❌ Failed at update PropertyTypes: %v", err)
		return existing.TID, err
	}

	if err := s.db.Save(&existing).Error; err != nil {
		log.Fatalf("❌ Failed at update NodeType: %v", err)
		return existing.TID, err
	}
	log.Println(fmt.Sprintf("🎉 Helper - Load %s schema successfully!"), existing.TID)
	return existing.TID, nil
}
