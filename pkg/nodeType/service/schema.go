package nodeType_service

import (
	"errors"
	"fmt"
	"go-cms-service/pkg/nodetype/model"
	"go-cms-service/pkg/nodetype/sql_helper"
	"go-cms-service/pkg/nodetype/utils"
	"go-cms-service/pkg/shared/utils"
	"go-cms-service/pkg/valuetype"
	"gorm.io/gorm"
	"log"
	"strings"
	"sync"
)

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
	if err := s.db.Exec(sql_helper.QueryCreateNewTable(nodeType)).Error; err != nil {
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

func (s *NodeTypeService) deleteColumn(tid, pid string) error {
	log.Printf("Delete column: %s in table %s", pid, tid)
	return s.db.Exec(sql_helper.QueryDeleteColumnFromTable(tid, pid)).Error
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
			if valuetype.MapValueTypeToSQL(pt.ValueType) != valuetype.MapValueTypeToSQL(newPT.ValueType) {
				if err := s.deleteColumn(newNodeType.TID, pid); err != nil {
					if !strings.Contains(err.Error(), "no such column") {
						log.Printf("❌ Error delete column %s: %v\n", pt.PID, err)
						return newNodeType.TID, nil
					}
				}
				if err := s.db.Unscoped().Delete(pt).Error; err != nil {
					log.Printf("❌ Failed to delete PropertyType (pid=%s): %v", pid, err)
					return newNodeType.TID, nil
				}
				toCreate = append(toCreate, newPT)
			} else {
				pt.ValueType = newPT.ValueType
				if err := s.db.Save(pt).Error; err != nil {
					log.Printf("❌ Failed to update PropertyType (pid=%s): %v", pid, err)
					return newNodeType.TID, nil
				}
			}
		} else {
			if err := s.deleteColumn(newNodeType.TID, pid); err != nil {
				if !strings.Contains(err.Error(), "no such column") {
					log.Printf("❌ Error delete column %s: %v\n", pt.PID, err)
					return newNodeType.TID, nil
				}
			}
			if err := s.db.Unscoped().Delete(pt).Error; err != nil {
				log.Printf("❌ Failed to delete PropertyType (pid=%s): %v", pid, err)
				return newNodeType.TID, nil
			}
		}
	}

	for _, pt := range toCreate {
		sql := sql_helper.QueryAddColumnToTable(newNodeType.TID, pt)
		if len(sql) == 0 {
			continue
		}
		if err := s.db.Exec(sql).Error; err != nil {
			log.Printf("❌ Failed at AutoMigrate: %v", err)
			return newNodeType.TID, nil
		}
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
