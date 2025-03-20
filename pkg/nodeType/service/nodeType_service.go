package nodeType_service

import (
	"errors"
	"fmt"
	"go-cms-service/pkg/nodeType/dynamic_struct"
	"go-cms-service/pkg/nodeType/model"
	"go-cms-service/pkg/nodeType/utils"
	"go-cms-service/pkg/shared/dto"
	"go-cms-service/pkg/shared/utils"
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
	if err := s.createNewTable(nodeType); err != nil {
		log.Printf("❌ Failed at create Table: %v", err)
		return nodeType.TID, err
	}
	if err := s.db.Create(&nodeType).Error; err != nil {
		log.Printf("❌ Failed at save NodeType: %v", err)
		return nodeType.TID, err
	}
	log.Printf("🎉 Helper - Load new %s schema successfully!", nodeType.TID)
	return nodeType.TID, nil
}

func (s *NodeTypeService) createNewTable(nodeType *nodeType_model.NodeType) error {
	return s.db.Exec(dynamic_struct.CreateDynamicTable(nodeType)).Error
}

func (s *NodeTypeService) deleteColumn(tid, pid string) error {
	log.Printf("Delete column: %s in table %s", pid, tid)
	sql := fmt.Sprintf("ALTER TABLE %s DROP COLUMN %s", tid, pid)
	return s.db.Exec(sql).Error
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

	for pid, pt := range currentMap {
		if newPT, ok := newMap[pid]; ok {
			if pt.ValueType != newPT.ValueType {
				if err := s.deleteColumn(newNodeType.TID, pid); err != nil {
					log.Printf("❌ Error update column %s: %v", pt.PID, err)
					return newNodeType.TID, nil
				}
			}
			pt.ValueType = newPT.ValueType
			if err := s.db.Save(pt).Error; err != nil {
				log.Printf("❌ Failed to update PropertyType (pid=%s): %v", pid, err)
				return newNodeType.TID, nil
			}
		} else {
			if err := s.deleteColumn(newNodeType.TID, pid); err != nil {
				log.Printf("❌ Error delete column %s: %v\n", pt.PID, err)
				return newNodeType.TID, nil
			}
			if err := s.db.Unscoped().Delete(pt).Error; err != nil {
				log.Printf("❌ Failed to delete PropertyType (pid=%s): %v", pid, err)
				return newNodeType.TID, nil
			}
		}
	}

	if err := s.db.Exec(dynamic_struct.CreateDynamicTable(newNodeType)).Error; err != nil {
		log.Printf("❌ Failed at AutoMigrate: %v", err)
		return newNodeType.TID, nil
	}
	for _, pt := range toCreate {
		fmt.Println("create new PropertyType", pt.PID)
		if err := s.db.Create(pt).Error; err != nil {
			log.Printf("❌ Failed to create new PropertyType: %v", err)
			return newNodeType.TID, nil
		}
	}

	newNodeType.ID = existing.ID
	if err := s.db.Omit("PropertyTypes").Save(&newNodeType).Error; err != nil {
		log.Printf("❌ Failed at update NodeType(%s): %v", newNodeType.TID, err)
		return newNodeType.TID, err
	}
	log.Printf("🎉 Helper - Load %s schema successfully!", newNodeType.TID)
	return newNodeType.TID, nil
}

func (s *NodeTypeService) CheckNodeTypeExist(tid string) bool {
	return s.db.Migrator().HasTable(tid)
}

func (s *NodeTypeService) FetchRecords(tid string) (*[]map[string]interface{}, error) {
	var result []map[string]interface{}
	if err := s.db.Table(tid).Find(&result).Error; err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *NodeTypeService) FetchRecord(tid string, id string) (*map[string]interface{}, error) {
	var result map[string]interface{}
	if err := s.db.Table(tid).Find(&result, "id = ?", id).Error; err != nil {
		return nil, err
	}
	if result == nil {
		return nil, nil
	}
	return &result, nil
}

func (s *NodeTypeService) CreateRecord(tid string, data map[string]interface{}) (*map[string]interface{}, error) {
	data["id"] = dynamic_struct.GenerateID()
	if result := s.db.Table(tid).Create(&data); result.Error != nil {
		return &data, result.Error
	}
	delete(data, "@id")
	return &data, nil
}

func (s *NodeTypeService) UpdateRecord(tid string, id string, data map[string]interface{}) (*map[string]interface{}, error) {
	delete(data, "id")
	result := s.db.Table(tid).Where("id = ?", id).Updates(&data)
	if result.Error != nil {
		return nil, result.Error
	}
	return &data, nil
}

func (s *NodeTypeService) DeleteRecord(tid string, id string) error {
	return s.db.Table(tid).Where("id = ?", id).Delete(nil).Error
}
