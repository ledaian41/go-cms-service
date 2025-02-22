package nodeType_service

import (
	"fmt"
	"go-product-service/pkg/nodeType/model"
	"go-product-service/pkg/nodeType/utils"
	"go-product-service/pkg/shared/dto"
	shared_utils "go-product-service/pkg/shared/utils"
	"gorm.io/gorm"
	"log"
	"sync"
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
		log.Println("❌ Failed at query NodeTypes: %v", err)
	}
	var dtos []shared_dto.NodeTypeDTO
	for _, n := range nodeTypes {
		dtos = append(dtos, n.NodeTypeDTO())
	}
	return &dtos
}

func (s *NodeTypeService) LoadSchema(path string, ch chan<- string) {
	if shared_utils.IsJsonPath(path) {
		s.loadSchemaFile(path, ch)
		return
	}
	if shared_utils.IsDirectory(path) {
		s.loadSchemaDirectory(path, ch)
		return
	}
}

func (s *NodeTypeService) loadSchemaFile(path string, ch chan<- string) {
	newNodeType, err := nodeType_utils.ReadSchemaJson(path)
	if err != nil {
		log.Println("❌ Failed at LoadSchema: %v", err)
		close(ch)
		return
	}

	var existing nodeType_model.NodeType
	if err := s.db.Preload("PropertyTypes").Where("t_id = ?", newNodeType.TID).First(&existing).Error; err != nil {
		tid, _ := s.createNewNodeType(newNodeType)
		ch <- tid
		close(ch)
		return
	}

	tid, _ := s.updateNodeType(&existing, newNodeType)
	ch <- tid
	close(ch)
	return
}

func (s *NodeTypeService) loadSchemaDirectory(path string, ch chan<- string) {
	nodeTypes, err := nodeType_utils.ReadSchemasFromDir(path)
	if err != nil {
		log.Println("❌ Failed at LoadSchema: %v", err)
		close(ch)
		return
	}

	var wg sync.WaitGroup
	for _, nodeType := range nodeTypes {
		wg.Add(1)
		go func(nodeType *nodeType_model.NodeType) {
			defer wg.Done()
			var existing nodeType_model.NodeType
			if err := s.db.Preload("PropertyTypes").Where("t_id = ?", nodeType.TID).First(&existing).Error; err != nil {
				tid, _ := s.createNewNodeType(nodeType)
				ch <- tid
			} else {
				tid, _ := s.updateNodeType(&existing, nodeType)
				ch <- tid
			}
		}(nodeType)
	}
	go func() {
		wg.Wait()
		close(ch)
	}()
}

func (s *NodeTypeService) createNewNodeType(nodeType *nodeType_model.NodeType) (string, error) {
	if err := s.db.Create(&nodeType).Error; err != nil {
		log.Println("❌ Failed at save NodeType: %v", err)
		return nodeType.TID, err
	}
	log.Println(fmt.Sprintf("🎉 Helper - Load new %s schema successfully!", nodeType.TID))
	return nodeType.TID, nil
}

func (s *NodeTypeService) updateNodeType(existing *nodeType_model.NodeType, newNodeType *nodeType_model.NodeType) (string, error) {
	existing.PropertyTypes = newNodeType.PropertyTypes
	if err := s.db.Model(&existing).Association("PropertyTypes").Replace(newNodeType.PropertyTypes); err != nil {
		log.Println("❌ Failed at update PropertyTypes: %v", err)
		return existing.TID, err
	}

	if err := s.db.Save(&existing).Error; err != nil {
		log.Println("❌ Failed at update NodeType: %v", err)
		return existing.TID, err
	}
	log.Println(fmt.Sprintf("🎉 Helper - Load %s schema successfully!", existing.TID))
	return existing.TID, nil
}
