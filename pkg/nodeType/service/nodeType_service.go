package nodeType_service

import (
	"errors"
	"fmt"
	"go-product-service/pkg/nodeType/model"
	"go-product-service/pkg/nodeType/utils"
	"go-product-service/pkg/shared/dto"
	"go-product-service/pkg/shared/utils"
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
	nodeType, err := nodeType_utils.ReadSchemaJson(path)
	if err != nil {
		log.Printf("❌ Failed at LoadSchema: %v", err)
		close(ch)
		return
	}
	s.loadNodeTypeToDB(nodeType, ch)
	close(ch)
}

func (s *NodeTypeService) loadSchemaDirectory(path string, ch chan<- string) {
	nodeTypes, err := nodeType_utils.ReadSchemasFromDir(path)
	if err != nil {
		log.Printf("❌ Failed at LoadSchema: %v", err)
		close(ch)
		return
	}

	var wg sync.WaitGroup
	for _, nodeType := range nodeTypes {
		wg.Add(1)
		go func(nodeType *nodeType_model.NodeType) {
			defer wg.Done()
			s.loadNodeTypeToDB(nodeType, ch)
		}(nodeType)
	}
	go func() {
		wg.Wait()
		close(ch)
	}()
}

func (s *NodeTypeService) loadNodeTypeToDB(nodeType *nodeType_model.NodeType, ch chan<- string) {
	var existing nodeType_model.NodeType
	if err := s.db.Preload("PropertyTypes").Where("tid = ?", nodeType.TID).First(&existing).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			tid, _ := s.createNewNodeType(nodeType)
			ch <- tid
		} else {
			log.Printf("❌ Error loading record: %v", err)
		}
		return
	}
	tid, _ := s.updateNodeType(&existing, nodeType)
	ch <- tid
}

func (s *NodeTypeService) createNewNodeType(nodeType *nodeType_model.NodeType) (string, error) {
	if err := s.db.Create(&nodeType).Error; err != nil {
		log.Printf("❌ Failed at save NodeType: %v", err)
		return nodeType.TID, err
	}
	log.Printf("🎉 Helper - Load new %s schema successfully!", nodeType.TID)
	return nodeType.TID, nil
}

func (s *NodeTypeService) updateNodeType(existing *nodeType_model.NodeType, newNodeType *nodeType_model.NodeType) (string, error) {
	var currentPTs []*nodeType_model.PropertyType
	if err := s.db.Model(&existing).Association("PropertyTypes").Find(&currentPTs); err != nil {
		log.Printf("❌ Failed at query PropertyTypes: %v", err)
		return newNodeType.TID, nil
	}

	currentMap := make(map[string]*nodeType_model.PropertyType)
	for _, pt := range currentPTs {
		currentMap[pt.PID] = pt
	}

	newMap := make(map[string]*nodeType_model.PropertyType)
	var toCreate []*nodeType_model.PropertyType
	for _, pt := range newNodeType.PropertyTypes {
		pt.NodeTypeRefer = existing.ID
		if currentMap[pt.PID] != nil {
			newMap[pt.PID] = pt
		} else {
			toCreate = append(toCreate, pt)
		}
	}

	tx := s.db.Begin()
	for pid, pt := range currentMap {
		if newPT, ok := newMap[pid]; ok {
			pt.ValueType = newPT.ValueType
			if err := tx.Save(pt).Error; err != nil {
				tx.Rollback()
				log.Printf("failed to update PropertyType (pid=%s): %v", pid, err)
				return newNodeType.TID, nil
			}
		} else {
			if err := tx.Unscoped().Delete(pt).Error; err != nil {
				tx.Rollback()
				log.Printf("failed to delete PropertyType (pid=%s): %v", pid, err)
				return newNodeType.TID, nil
			}
		}
	}
	for _, pt := range toCreate {
		fmt.Println("create new PropertyType", pt.PID)
		if err := tx.Create(pt).Error; err != nil {
			tx.Rollback()
			log.Printf("❌ Failed to create new PropertyType: %v", err)
			return newNodeType.TID, nil
		}
	}
	if err := tx.Commit().Error; err != nil {
		log.Printf("❌ Commit transaction failed: %v", err)
	}

	newNodeType.ID = existing.ID
	if err := s.db.Omit("PropertyTypes").Save(&newNodeType).Error; err != nil {
		log.Printf("❌ Failed at update NodeType(%s): %v", newNodeType.TID, err)
		return newNodeType.TID, err
	}
	log.Printf("🎉 Helper - Load %s schema successfully!", newNodeType.TID)
	return newNodeType.TID, nil
}
